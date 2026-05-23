package protection

import (
	"testing"
	"time"
)

// Test that idle sessions are removed by the cleaner goroutine.
func TestSessionManagerIdleCleanup(t *testing.T) {
	// Create SessionManager with short intervals for testing.
	sm := &SessionManager{
		sessions: make(map[string]*Session),
		stopChan: make(chan struct{}),
	}

	// Start cleaner manually with short intervals.
	go sm.startCleaner(30*time.Millisecond, 50*time.Millisecond)

	// Create idle session (already idle for >50ms)
	sess := &Session{
		ID:         "alice",
		LastActive: time.Now().Add(-200 * time.Millisecond),
		CloseFunc:  func() {},
	}

	sm.Add("alice", sess)

	// Wait long enough for cleaner to run
	time.Sleep(120 * time.Millisecond)

	// Stop cleaner to avoid goroutine leak
	sm.Stop()

	if _, ok := sm.Get("alice"); ok {
		t.Fatalf("expected idle session to be removed")
	}
}
