package api

import (
	"context"
	"sync"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

// ClientState represents the high-level lifecycle state of the client.
// It is intentionally coarse-grained: UI layers decide how to interpret it.
type ClientState int

const (
	StateDisconnected    ClientState = iota // No QUIC connection
	StateHandshaking                        // HELLO sent, waiting for MEOW_OK
	StateAuthenticating                     // AUTH sent, waiting for MEOW_OK
	StateSelectingTarget                    // Logged in, waiting for UI to choose target
	StateEstablished                        // Ready for encrypted DATA exchange
)

// KittyClient is the core client structure.
// It contains no UI logic and is safe to use from any frontend (CLI, GUI, etc.).
type KittyClient struct {
	mu    sync.Mutex
	state ClientState

	// QUIC transport
	conn   *quic.Conn
	stream *quic.Stream

	// Session metadata
	user   string
	target string

	// Subsystems
	ackMgr   *AckManager
	replay   *protection.ReplayDetector
	stopPing chan struct{}
	stopRecv chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc

	// Debug/testing
	lastFrame []byte // last raw frame (used only for replay testing)

	// E2EE keys
	kEnc []byte // encryption key (AES-GCM)
	kMac []byte // MAC key (HMAC-SHA256)
}

// NewKittyClient creates a new client instance in the Disconnected state.
// It initializes ACK manager, replay detector, and background control channels.
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
	}
}
