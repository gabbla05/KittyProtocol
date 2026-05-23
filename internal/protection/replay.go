package protection

import (
	"sync"
	"time"
)

const (
	// Maximum number of tracked message IDs before forced cleanup.
	maxReplayEntries = 10_000

	// TTL for replay entries.
	replayTTL = 2 * time.Minute

	// Sweep interval (always performed, not only when map is large).
	replaySweepInterval = 5 * time.Second
)

// ReplayDetector tracks recently seen message IDs to detect replays.
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
		if now.Sub(ts) <= replayTTL {
			return true
		}
	}

	// Always sweep periodically
	if now.Sub(r.lastSweep) >= replaySweepInterval {
		for id, ts := range r.seen {
			if now.Sub(ts) > replayTTL {
				delete(r.seen, id)
			}
		}
		r.lastSweep = now
	}

	// Enforce memory limit
	if len(r.seen) >= maxReplayEntries {
		// Remove oldest entries
		cutoff := now.Add(-replayTTL)
		for id, ts := range r.seen {
			if ts.Before(cutoff) {
				delete(r.seen, id)
			}
			if len(r.seen) < maxReplayEntries {
				break
			}
		}
	}

	// Save entry
	r.seen[msgID] = now
	return false
}
