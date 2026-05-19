package api

import (
	"sync"
	"time"
)

// AckEventHandler defines callbacks for ACK events.
// UI (CLI or Wails) can subscribe to these.
type AckEventHandler interface {
	OnDelivered(msgID int64)
	OnTimeout(msgID int64)
}

// AckManager tracks pending messages and their timers.
type AckManager struct {
	mu       sync.Mutex
	pending  map[int64]chan struct{}
	handlers []AckEventHandler
	timeout  time.Duration
}

// NewAckManager creates a new manager with default timeout (5s).
func NewAckManager() *AckManager {
	return &AckManager{
		pending: make(map[int64]chan struct{}),
		timeout: 5 * time.Second,
	}
}

// RegisterHandler allows UI or client code to receive ACK events.
func (a *AckManager) RegisterHandler(h AckEventHandler) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.handlers = append(a.handlers, h)
}

// AddPending registers a message ID and starts a timeout timer.
func (a *AckManager) AddPending(msgID int64) {
	ch := make(chan struct{})

	a.mu.Lock()
	a.pending[msgID] = ch
	a.mu.Unlock()

	// Start timeout goroutine
	go func() {
		select {
		case <-ch:
			// Delivered — nothing to do here
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
