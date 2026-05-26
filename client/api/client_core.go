package api

import (
	"context"
	"sync"

	"github.com/gabbla05/KittyProtocol/internal/protection"
)

// ClientState represents the high-level lifecycle state of KittyClient.
// It is intentionally coarse-grained and UI-agnostic.
type ClientState int

const (
	StateDisconnected ClientState = iota
	StateHandshaking
	StateAuthenticating
	StateRegistering
	StateSelectingTarget
	StateEstablished
)

// AppPayloadHandler is invoked for application-level payloads that are not
// recognized as chat control frames. It is typically used by higher layers
// (e.g. GUI, bots) to process custom messages.
type AppPayloadHandler func(sender string, payload []byte)

// ErrorHandler receives protocol-level errors that are not directly tied
// to HELLO / AUTH / REGISTER operations.
type ErrorHandler func(code, desc string)

// StatusHandler receives presence/status updates for a given user.
type StatusHandler func(target, status string)

// DisconnectHandler is invoked when the underlying transport is broken
// or the receiver loop terminates due to a read error.
type DisconnectHandler func(err error)

type peerKeys struct {
	kEnc []byte
	kMac []byte
}

// KittyClient is the main high-level client for KittyProtocol.
//
// It encapsulates:
//   - QUIC transport (ConnAdapter / StreamAdapter),
//   - TLS + TOFU certificate pinning,
//   - E2EE key management,
//   - replay protection,
//   - ACK tracking,
//   - chat events and application payload dispatch,
//   - lifecycle state machine.
//
// The type is safe for concurrent use by multiple goroutines.
type KittyClient struct {
	mu    sync.Mutex
	state ClientState

	conn   ConnAdapter
	stream StreamAdapter

	user   string
	target string

	ackMgr   *AckManager
	replay   *protection.ReplayDetector
	stopPing chan struct{}
	stopRecv chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc

	// lastFrame is a development-only helper used for replay testing.
	lastFrame []byte

	peerKeys map[string]peerKeys

	appHandler        AppPayloadHandler
	errHandler        ErrorHandler
	statusHandler     StatusHandler
	disconnectHandler DisconnectHandler

	helloCh    chan OpResult
	authCh     chan OpResult
	registerCh chan OpResult

	chatReqCh    chan ChatRequestEvent
	chatAcceptCh chan ChatAcceptEvent
	chatRefuseCh chan ChatRefuseEvent
	chatEndCh    chan ChatEndEvent
	chatMsgCh    chan ChatMessageEvent
}

// NewKittyClient constructs a new client instance with a fresh internal
// context, replay detector, ACK manager and all event channels.
//
// The returned client is in StateDisconnected and has no active transport.
func NewKittyClient() *KittyClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &KittyClient{
		state:    StateDisconnected,
		ackMgr:   NewAckManager(),
		replay:   protection.NewReplayDetector(),
		stopPing: make(chan struct{}),
		stopRecv: make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
		peerKeys: make(map[string]peerKeys),

		helloCh:    make(chan OpResult, 1),
		authCh:     make(chan OpResult, 1),
		registerCh: make(chan OpResult, 1),

		chatReqCh:    make(chan ChatRequestEvent, 4),
		chatAcceptCh: make(chan ChatAcceptEvent, 4),
		chatRefuseCh: make(chan ChatRefuseEvent, 4),
		chatEndCh:    make(chan ChatEndEvent, 4),
		chatMsgCh:    make(chan ChatMessageEvent, 16),
	}
}
