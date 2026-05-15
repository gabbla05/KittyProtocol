package protection

import (
	"testing"
	"time"
)

// Test that idle sessions are removed by the cleaner goroutine.
func TestSessionManagerIdleCleanup(t *testing.T) {
	sm := NewSessionManagerWithInterval(50*time.Millisecond, 100*time.Millisecond)

	sess := &Session{
		ID:         "alice",
		LastActive: time.Now().Add(-200 * time.Millisecond),
		CloseFunc:  func() {},
	}

	sm.Add("alice", sess)

	time.Sleep(200 * time.Millisecond)

	if _, ok := sm.Get("alice"); ok {
		t.Fatalf("expected idle session to be removed")
	}
}
