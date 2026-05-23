package protection

import (
	"testing"
	"time"
)

func TestReplayDetector_FirstSeenIsNotReplay(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("first time msgID=%d should NOT be replay", id)
	}
}

func TestReplayDetector_SecondSeenIsReplay(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("first time msgID=%d should NOT be replay", id)
	}
	if replay := r.MarkAndCheck(id); !replay {
		t.Fatalf("second time msgID=%d SHOULD be replay", id)
	}
}

func TestReplayDetector_TTLExpires(t *testing.T) {
	r := NewReplayDetector()
	id := int64(123)

	// First time → not replay
	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("first time should NOT be replay")
	}

	// Second time immediately → replay
	if replay := r.MarkAndCheck(id); !replay {
		t.Fatalf("second time SHOULD be replay")
	}

	// --- symulacja upływu czasu ---
	r.mu.Lock()
	r.seen[id] = time.Now().Add(-replayTTL - time.Second)
	r.mu.Unlock()

	// Now TTL expired → should NOT be replay
	if replay := r.MarkAndCheck(id); replay {
		t.Fatalf("after TTL msgID should NOT be replay")
	}
}

func TestReplayDetector_SweepRemovesOldEntries(t *testing.T) {
	r := NewReplayDetector()

	oldID := int64(1)
	newID := int64(2)

	// Insert old entry
	r.MarkAndCheck(oldID)

	// Cofamy czas starego wpisu
	r.mu.Lock()
	r.seen[oldID] = time.Now().Add(-replayTTL - time.Second)
	r.lastSweep = time.Now().Add(-replaySweepInterval - time.Second)
	r.mu.Unlock()

	// Trigger sweep
	r.MarkAndCheck(newID)

	r.mu.Lock()
	_, exists := r.seen[oldID]
	r.mu.Unlock()

	if exists {
		t.Fatalf("old entry should have been swept out")
	}
}

func TestReplayDetector_MaxEntriesLimit(t *testing.T) {
	r := NewReplayDetector()

	// Wypełniamy mapę do limitu
	for i := 0; i < maxReplayEntries; i++ {
		r.MarkAndCheck(int64(i))
	}

	// Cofamy czas części wpisów, aby mogły zostać usunięte
	r.mu.Lock()
	cutoff := time.Now().Add(-replayTTL - time.Second)
	for id := range r.seen {
		r.seen[id] = cutoff
	}
	r.mu.Unlock()

	// Dodanie nowego wpisu powinno wywołać cleanup
	r.MarkAndCheck(999999)

	if len(r.seen) > maxReplayEntries {
		t.Fatalf("map should not exceed maxReplayEntries after cleanup")
	}
}
