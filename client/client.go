package main

import (
	"context"
	"sync"

	"github.com/quic-go/quic-go"
)

// ClientState represents the high-level state of the client.
type ClientState int

const (
	StateDisconnected    ClientState = iota
	StateHandshaking                 // After sending HELLO
	StateAuthenticating              // After sending AUTH
	StateSelectingTarget             // Waiting for target selection (UI-level)
	StateEstablished                 // Ready for DATA exchange
)

// KittyClient is the core client structure.
// It is UI-agnostic: no stdin/stdout, no CLI, no Wails-specific code.
// Any UI (CLI, Wails, Fyne) should talk to this type.
type KittyClient struct {
	mu    sync.Mutex
	state ClientState

	conn   *quic.Conn
	stream *quic.Stream

	user   string
	target string

	ackMgr   *AckManager
	replay   *ReplayDetector
	stopPing chan struct{}
	stopRecv chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc

	lastFrame []byte
}

// NewKittyClient creates a new client instance in the Disconnected state.
func NewKittyClient() *KittyClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &KittyClient{
		state:    StateDisconnected,
		ackMgr:   NewAckManager(),
		replay:   NewReplayDetector(),
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

// setState changes the client state (internal use only).
func (c *KittyClient) setState(newState ClientState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = newState
}

// SetTarget sets the current chat target (UI-level decision).
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

// Close gracefully shuts down the client:
// - stops ping and receiver loops,
// - closes the QUIC stream and connection,
// - cancels the internal context.
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
	}
	if c.conn != nil {
		_ = c.conn.CloseWithError(0, "client closed")
	}
}
