package api

import (
	"sync"
	"testing"
	"time"
)

// testHandler is a helper used to capture ACK events in tests.
type testHandler struct {
	mu        sync.Mutex
	delivered []int64
	timeout   []int64
}

func (h *testHandler) OnDelivered(msgID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.delivered = append(h.delivered, msgID)
}

func (h *testHandler) OnTimeout(msgID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.timeout = append(h.timeout, msgID)
}

func TestAckManager_Delivered(t *testing.T) {
	ack := NewAckManager()
	ack.timeout = 200 * time.Millisecond // speed up tests

	h := &testHandler{}
	ack.RegisterHandler(h)

	msgID := int64(123)
	ack.AddPending(msgID)

	// Simulate MEOW_OK
	ack.NotifyDelivered(msgID)

	time.Sleep(50 * time.Millisecond)

	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.delivered) != 1 || h.delivered[0] != msgID {
		t.Fatalf("expected delivered=%d, got %v", msgID, h.delivered)
	}
	if len(h.timeout) != 0 {
		t.Fatalf("unexpected timeout event: %v", h.timeout)
	}
}

func TestAckManager_Timeout(t *testing.T) {
	ack := NewAckManager()
	ack.timeout = 100 * time.Millisecond

	h := &testHandler{}
	ack.RegisterHandler(h)

	msgID := int64(456)
	ack.AddPending(msgID)

	time.Sleep(200 * time.Millisecond)

	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.timeout) != 1 || h.timeout[0] != msgID {
		t.Fatalf("expected timeout=%d, got %v", msgID, h.timeout)
	}
	if len(h.delivered) != 0 {
		t.Fatalf("unexpected delivered event: %v", h.delivered)
	}
}

func TestAckManager_NoDoubleEvents(t *testing.T) {
	ack := NewAckManager()
	ack.timeout = 200 * time.Millisecond

	h := &testHandler{}
	ack.RegisterHandler(h)

	msgID := int64(789)
	ack.AddPending(msgID)

	// Delivered before timeout
	ack.NotifyDelivered(msgID)

	time.Sleep(300 * time.Millisecond)

	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.delivered) != 1 {
		t.Fatalf("expected exactly one delivered event, got %v", h.delivered)
	}
	if len(h.timeout) != 0 {
		t.Fatalf("timeout should NOT fire after delivered, got %v", h.timeout)
	}
}

func TestAckManager_MultipleHandlers(t *testing.T) {
	ack := NewAckManager()
	ack.timeout = 200 * time.Millisecond

	h1 := &testHandler{}
	h2 := &testHandler{}
	ack.RegisterHandler(h1)
	ack.RegisterHandler(h2)

	msgID := int64(111)
	ack.AddPending(msgID)
	ack.NotifyDelivered(msgID)

	time.Sleep(50 * time.Millisecond)

	h1.mu.Lock()
	h2.mu.Lock()
	defer h1.mu.Unlock()
	defer h2.mu.Unlock()

	if len(h1.delivered) != 1 || len(h2.delivered) != 1 {
		t.Fatalf("both handlers must receive delivered event")
	}
}

func TestAckManager_MultiplePending(t *testing.T) {
	ack := NewAckManager()
	ack.timeout = 150 * time.Millisecond

	h := &testHandler{}
	ack.RegisterHandler(h)

	ack.AddPending(1)
	ack.AddPending(2)

	ack.NotifyDelivered(1)

	time.Sleep(250 * time.Millisecond)

	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.delivered) != 1 || h.delivered[0] != 1 {
		t.Fatalf("expected delivered=1, got %v", h.delivered)
	}
	if len(h.timeout) != 1 || h.timeout[0] != 2 {
		t.Fatalf("expected timeout=2, got %v", h.timeout)
	}
}
