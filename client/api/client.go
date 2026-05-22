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

// AppPayloadHandler is a callback used by the application layer (client/app)
// to receive decrypted DATA payloads from KittyClient.
type AppPayloadHandler func(sender string, payload []byte)

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

	// E2EE keys per peer (logical username → keys)
	peerKeys map[string]peerKeys

	// Application-level payload handler (chat, etc.)
	appHandler AppPayloadHandler
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

func (c *KittyClient) User() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.user
}

// getKeysForPeer returns derived keys for a given peer, if present.
func (c *KittyClient) getKeysForPeer(peer string) (kEnc, kMac []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.peerKeys == nil {
		return nil, nil, false
	}
	pk, exists := c.peerKeys[peer]
	if !exists {
		return nil, nil, false
	}
	return pk.kEnc, pk.kMac, true
}
