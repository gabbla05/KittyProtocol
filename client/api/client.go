package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
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

	conn   *quic.Conn
	stream *quic.Stream

	user   string
	target string

	ackMgr   *AckManager
	replay   *protection.ReplayDetector
	stopPing chan struct{}
	stopRecv chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc

	lastFrame []byte // last raw frame (used only for replay testing)

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

// State returns the current client state (thread-safe).
func (c *KittyClient) State() ClientState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// setState updates the internal state (not exported).
func (c *KittyClient) setState(newState ClientState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = newState
}

// SetTarget sets the current chat target.
// This is a UI-level decision and does not involve protocol logic.
func (c *KittyClient) SetTarget(target string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.target = target
}

// Target returns the current chat target.
func (c *KittyClient) Target() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.target
}

// SetSharedSecret derives encryption and MAC keys from the shared secret.
// Keys are stored securely and zeroized on Close().
func (c *KittyClient) SetSharedSecret(secret []byte) error {
	kEnc, kMac, err := cryptoee.DeriveKeysFromSecret(secret)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.kEnc = kEnc
	c.kMac = kMac
	return nil
}

// RegisterAckHandler registers a UI or application component to receive ACK events.
func (c *KittyClient) RegisterAckHandler(h AckEventHandler) {
	if c.ackMgr != nil {
		c.ackMgr.RegisterHandler(h)
	}
}

// ReplayLastFrame resends the last raw frame written to the stream.
// This is used exclusively for replay protection testing.
func (c *KittyClient) ReplayLastFrame() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lastFrame == nil {
		return fmt.Errorf("no frame to replay")
	}
	if c.stream == nil {
		return fmt.Errorf("stream is nil")
	}

	_, err := c.stream.Write(c.lastFrame)
	return err
}

// Close gracefully shuts down the client:
//   - stops ping and receiver loops,
//   - closes the QUIC stream and connection,
//   - cancels the internal context,
//   - zeroizes encryption keys,
//   - resets replay detector and ACK manager,
//   - clears session state (target, lastFrame).
func (c *KittyClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Stop background loops
	select {
	case <-c.stopPing:
	default:
		close(c.stopPing)
	}
	select {
	case <-c.stopRecv:
	default:
		close(c.stopRecv)
	}

	// Cancel context
	if c.cancel != nil {
		c.cancel()
	}

	// Close stream and connection
	if c.stream != nil {
		_ = c.stream.Close()
		c.stream = nil
	}
	if c.conn != nil {
		_ = c.conn.CloseWithError(0, "client closed")
		c.conn = nil
	}

	// Zeroize keys
	if c.kEnc != nil {
		cryptoee.Zeroize(c.kEnc)
		c.kEnc = nil
	}
	if c.kMac != nil {
		cryptoee.Zeroize(c.kMac)
		c.kMac = nil
	}

	// Reset session state
	c.target = ""
	c.lastFrame = nil
	c.replay = protection.NewReplayDetector()
	c.ackMgr = NewAckManager()

	c.state = StateDisconnected
}
