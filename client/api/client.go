package api

import (
	"context"
	"sync"

	"github.com/gabbla05/KittyProtocol/internal/protection"
)

// ClientState represents the high-level lifecycle state of the client.
// It is intentionally coarse-grained: UI layers decide how to interpret it.
type ClientState int

// AppPayloadHandler is a callback used by the application layer (client/app)
// to receive decrypted DATA payloads from KittyClient.
type AppPayloadHandler func(sender string, payload []byte)

// ErrorHandler is invoked when the Hub sends an ERROR frame.
type ErrorHandler func(code, desc string)

// StatusHandler is invoked when the Hub sends a STATUS_RES frame.
type StatusHandler func(target, status string)

// DisconnectHandler is invoked when the underlying connection is closed
// or the receiver loop terminates with an error.
type DisconnectHandler func(err error)

const (
	StateDisconnected    ClientState = iota // No QUIC connection
	StateHandshaking                        // HELLO sent, waiting for MEOW_OK
	StateAuthenticating                     // AUTH sent, waiting for MEOW_OK
	StateSelectingTarget                    // Logged in, waiting for UI to choose target
	StateEstablished                        // Ready for encrypted DATA exchange
)

// peerKeys holds derived encryption and MAC keys for a single peer.
type peerKeys struct {
	kEnc []byte
	kMac []byte
}

// KittyClient is the core client structure.
// It contains no UI logic and is safe to use from any frontend (CLI, GUI, etc.).
type KittyClient struct {
	mu    sync.Mutex
	state ClientState

	// QUIC transport (wrapped in adapters for testability and GUI-friendliness)
	conn   ConnAdapter
	stream StreamAdapter

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

	// E2EE keys per peer (logical username → keys)
	peerKeys map[string]peerKeys

	// Application-level payload handler (chat, etc.)
	appHandler AppPayloadHandler

	// Event handlers for UI/frontends
	errHandler        ErrorHandler
	statusHandler     StatusHandler
	disconnectHandler DisconnectHandler
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
		peerKeys: make(map[string]peerKeys),
	}
}

func (c *KittyClient) RegisterAppPayloadHandler(h AppPayloadHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.appHandler = h
}

func (c *KittyClient) OnError(h ErrorHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errHandler = h
}

func (c *KittyClient) OnStatus(h StatusHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statusHandler = h
}

func (c *KittyClient) OnDisconnected(h DisconnectHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.disconnectHandler = h
}

func (c *KittyClient) User() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.user
}

// getKeysForPeer returns derived keys for a given peer, if present.
// Caller MUST hold c.mu.
func (c *KittyClient) getKeysForPeer(peer string) (kEnc, kMac []byte, ok bool) {
	if c.peerKeys == nil {
		return nil, nil, false
	}
	pk, exists := c.peerKeys[peer]
	if !exists {
		return nil, nil, false
	}
	return pk.kEnc, pk.kMac, true
}
