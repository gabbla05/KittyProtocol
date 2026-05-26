package protection

import (
	"sync"
	"time"
)

// ReplayDetector tracks recently seen message IDs to detect replay attacks.
// It is used per-session to ensure that clients cannot resend old DATA frames.
type ReplayDetector struct {
	mu        sync.Mutex
	seen      map[int64]time.Time
	lastSweep time.Time
}

// NewReplayDetector creates a new ReplayDetector instance.
func NewReplayDetector() *ReplayDetector {
	return &ReplayDetector{
		seen:      make(map[int64]time.Time),
		lastSweep: time.Now(),
	}
}

// MarkAndCheck records the given msgID and returns true if it is a replay.
func (r *ReplayDetector) MarkAndCheck(msgID int64) bool {
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Replay check
	if ts, ok := r.seen[msgID]; ok {
		if now.Sub(ts) <= ReplayTTL {
			return true
		}
	}

	// Periodic sweep
	if now.Sub(r.lastSweep) >= ReplaySweepInterval {
		for id, ts := range r.seen {
			if now.Sub(ts) > ReplayTTL {
				delete(r.seen, id)
			}
		}
		r.lastSweep = now
	}

	// Enforce memory limit
	if len(r.seen) >= MaxReplayEntries {
		cutoff := now.Add(-ReplayTTL)
		for id, ts := range r.seen {
			if ts.Before(cutoff) {
				delete(r.seen, id)
			}
			if len(r.seen) < MaxReplayEntries {
				break
			}
		}
	}

	// Save entry
	r.seen[msgID] = now
	return false
}
