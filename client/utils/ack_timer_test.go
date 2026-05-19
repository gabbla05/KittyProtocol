package clientutils

import (
	"testing"
	"time"
)

func TestAckTimerTimeout(t *testing.T) {
	done := make(chan struct{})
	called := make(chan struct{})

	StartAckTimer(123, done, func() {
		close(called)
	})

	select {
	case <-called:
		// OK
	case <-time.After(DefaultAckTimeout + 200*time.Millisecond):
		t.Fatalf("expected timeout callback to be called")
	}
}

func TestAckTimerCancelled(t *testing.T) {
	done := make(chan struct{})
	called := make(chan struct{})

	StartAckTimer(123, done, func() {
		close(called)
	})

	// Cancel immediately
	close(done)

	select {
	case <-called:
		t.Fatalf("timeout callback should NOT be called when done is closed")
	case <-time.After(200 * time.Millisecond):
		// OK
	}
}
