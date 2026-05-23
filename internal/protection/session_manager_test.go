package protection

import (
	"testing"
	"time"
)

// TestSessionManagerIdleCleanup verifies that idle sessions are removed
// by the background cleaner goroutine.
func TestSessionManagerIdleCleanup(t *testing.T) {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
		stopChan: make(chan struct{}),
	}

	// Start cleaner with very short intervals for testing.
	go sm.startCleaner(30*time.Millisecond, 50*time.Millisecond)

	// Create a session that is already idle for >50ms.
	sess := &Session{
		ID:         "alice",
		LastActive: time.Now().Add(-200 * time.Millisecond),
		CloseFunc:  func() {},
	}

	sm.Add("alice", sess)

	// Wait long enough for cleaner to run.
	time.Sleep(120 * time.Millisecond)

	sm.Stop()

	if _, ok := sm.Get("alice"); ok {
		t.Fatalf("expected idle session to be removed by cleaner")
	}
}
