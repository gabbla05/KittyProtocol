package protection

import (
	"testing"
	"time"
)

// TestReplayDetector_FirstSeenIsNotReplay verifies that the first occurrence
// of a message ID is never treated as a replay.
func TestReplayDetector_FirstSeenIsNotReplay(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("first occurrence of msgID=%d should NOT be replay", id)
	}
}

// TestReplayDetector_SecondSeenIsReplay ensures that re-sending the same
// message ID within the TTL window is detected as a replay.
func TestReplayDetector_SecondSeenIsReplay(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	r.MarkAndCheck(id)
	if replay := r.MarkAndCheck(id); !replay {
		t.Fatalf("second occurrence of msgID=%d SHOULD be replay", id)
	}
}

// TestReplayDetector_TTLExpires verifies that replay entries expire after TTL.
func TestReplayDetector_TTLExpires(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	r.MarkAndCheck(id)
	r.MarkAndCheck(id)

	// Simulate TTL expiration
	r.mu.Lock()
	r.seen[id] = time.Now().Add(-ReplayTTL - time.Second)
	r.mu.Unlock()

	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("msgID should NOT be replay after TTL expiration")
	}
}

// TestReplayDetector_SweepRemovesOldEntries ensures that periodic sweeping
// removes expired entries.
func TestReplayDetector_SweepRemovesOldEntries(t *testing.T) {
	r := NewReplayDetector()

	oldID := int64(1)
	newID := int64(2)

	r.MarkAndCheck(oldID)

	// Simulate old timestamp + force sweep
	r.mu.Lock()
	r.seen[oldID] = time.Now().Add(-ReplayTTL - time.Second)
	r.lastSweep = time.Now().Add(-ReplaySweepInterval - time.Second)
	r.mu.Unlock()

	r.MarkAndCheck(newID)

	r.mu.Lock()
	_, exists := r.seen[oldID]
	r.mu.Unlock()

	if exists {
		t.Fatalf("expired entry should have been removed during sweep")
	}
}

// TestReplayDetector_MaxEntriesLimit ensures that the detector never grows
// beyond MaxReplayEntries.
func TestReplayDetector_MaxEntriesLimit(t *testing.T) {
	r := NewReplayDetector()

	for i := 0; i < MaxReplayEntries; i++ {
		r.MarkAndCheck(int64(i))
	}

	// Simulate all entries being old
	r.mu.Lock()
	cutoff := time.Now().Add(-ReplayTTL - time.Second)
	for id := range r.seen {
		r.seen[id] = cutoff
	}
	r.mu.Unlock()

	r.MarkAndCheck(999999)

	if len(r.seen) > MaxReplayEntries {
		t.Fatalf("ReplayDetector should not exceed MaxReplayEntries")
	}
}
