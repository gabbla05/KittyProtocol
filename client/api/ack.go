package api

import (
	"sync"
	"time"
)

// AckEventHandler defines callbacks for delivery acknowledgment events.
// UI layers (CLI, GUI) or application components may subscribe to receive
// notifications about message delivery or timeout.
//
// CONTRACT:
//   - OnDelivered(msgID) is called exactly once when MEOW_OK arrives.
//   - OnTimeout(msgID) is called exactly once when the timeout expires.
//   - A message will NEVER trigger both events.
//   - Handlers are invoked synchronously in the caller goroutine.
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

// NewAckManager creates a new manager with a default timeout of 5 seconds.
func NewAckManager() *AckManager {
	return &AckManager{
		pending: make(map[int64]chan struct{}),
		timeout: 5 * time.Second,
	}
}

// RegisterHandler registers a UI or client component to receive ACK events.
// Handlers are invoked synchronously in the caller goroutine.
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
//   - Only one of these events will fire.
func (a *AckManager) AddPending(msgID int64) {
	ch := make(chan struct{})

	a.mu.Lock()
	a.pending[msgID] = ch
	a.mu.Unlock()

	go func() {
		select {
		case <-ch:
			// Delivered — nothing more to do.
			return

		case <-time.After(a.timeout):
			a.mu.Lock()
			if _, ok := a.pending[msgID]; ok {
				delete(a.pending, msgID)
				for _, h := range a.handlers {
					h.OnTimeout(msgID)
				}
			}
			a.mu.Unlock()
		}
	}()
}

// NotifyDelivered is called when MEOW_OK arrives.
// It removes the pending entry and notifies all handlers.
func (a *AckManager) NotifyDelivered(msgID int64) {
	a.mu.Lock()
	ch, ok := a.pending[msgID]
	if ok {
		delete(a.pending, msgID)
	}
	a.mu.Unlock()

	if ok {
		close(ch)
		for _, h := range a.handlers {
			h.OnDelivered(msgID)
		}
	}
}

// RegisterAckHandler is the public API exposed by KittyClient.
// It simply forwards the handler to the underlying AckManager.
func (c *KittyClient) RegisterAckHandler(h AckEventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ackMgr != nil {
		c.ackMgr.RegisterHandler(h)
	}
}
