package api

import (
	"sync"
	"time"
)

// AckEventHandler defines callbacks for delivery acknowledgment events.
// UI layers (CLI, GUI) or higher-level application components may subscribe
// to receive notifications about message delivery or timeout.
//
// CONTRACT:
//   - OnDelivered(msgID) is called exactly once when MEOW_OK arrives.
//   - OnTimeout(msgID) is called exactly once when the timeout expires.
//   - A message will NEVER trigger both events.
//   - Handlers are invoked synchronously in the goroutine that processes
//     the event (no extra goroutines are spawned per handler).
type AckEventHandler interface {
	OnDelivered(msgID int64)
	OnTimeout(msgID int64)
}

// AckManager tracks pending messages and their timeout goroutines.
// Each pending message has a dedicated channel that is closed when MEOW_OK
// arrives. If the timeout fires first, OnTimeout is invoked.
//
// THREAD SAFETY:
//   - All operations are protected by a mutex.
//   - Safe for concurrent use by sender and receiver goroutines.
type AckManager struct {
	mu       sync.Mutex
	pending  map[int64]chan struct{}
	handlers []AckEventHandler
	timeout  time.Duration
}

// NewAckManager creates a new manager with a default timeout.
func NewAckManager() *AckManager {
	return &AckManager{
		pending: make(map[int64]chan struct{}),
		timeout: defaultAckTimeout,
	}
}

// RegisterHandler registers a component to receive ACK events.
// Handlers are invoked synchronously in the event-processing goroutine.
//
// It is safe to call RegisterHandler before or after messages are added,
// but handlers registered after AddPending() may miss earlier events.
func (a *AckManager) RegisterHandler(h AckEventHandler) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.handlers = append(a.handlers, h)
}

// AddPending registers a message ID and starts a timeout goroutine.
//
// BEHAVIOR:
//   - If MEOW_OK arrives → NotifyDelivered() closes the channel.
//   - If timeout expires → OnTimeout() is invoked.
//   - Only one of these events will fire for a given msgID.
func (a *AckManager) AddPending(msgID int64) {
	ch := make(chan struct{})

	a.mu.Lock()
	a.pending[msgID] = ch
	timeout := a.timeout
	a.mu.Unlock()

	go func() {
		select {
		case <-ch:
			// Delivered — nothing more to do.
			return

		case <-time.After(timeout):
			a.mu.Lock()
			// Check if still pending (may have been delivered just before timeout).
			if _, ok := a.pending[msgID]; ok {
				delete(a.pending, msgID)
				handlers := append([]AckEventHandler(nil), a.handlers...)
				a.mu.Unlock()

				for _, h := range handlers {
					h.OnTimeout(msgID)
				}
				return
			}
			a.mu.Unlock()
		}
	}()
}

// NotifyDelivered is called when MEOW_OK arrives.
// It removes the pending entry and notifies all handlers.
//
// If the message is no longer pending (e.g. timeout already fired),
// this call is a no-op.
func (a *AckManager) NotifyDelivered(msgID int64) {
	a.mu.Lock()
	ch, ok := a.pending[msgID]
	if ok {
		delete(a.pending, msgID)
	}
	handlers := append([]AckEventHandler(nil), a.handlers...)
	a.mu.Unlock()

	if !ok {
		return
	}

	// Closing the channel unblocks the timeout goroutine (if still waiting).
	close(ch)

	for _, h := range handlers {
		h.OnDelivered(msgID)
	}
}

// RegisterAckHandler is the public API exposed by KittyClient.
// It simply forwards the handler to the underlying AckManager.
//
// This keeps the AckManager internal while exposing a stable interface
// for UI layers (CLI, Wails GUI, tests, etc.).
func (c *KittyClient) RegisterAckHandler(h AckEventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ackMgr != nil {
		c.ackMgr.RegisterHandler(h)
	}
}
