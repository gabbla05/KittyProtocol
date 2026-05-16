package protection

import (
	"sync"
	"time"
)

const (
	// maxHubReplayEntries defines the maximum number of tracked message IDs
	// before a cleanup sweep is triggered.
	maxHubReplayEntries = 100_000

	// replayTTL defines how long a message ID is considered "recent" and
	// thus subject to replay detection.
	replayTTL = 2 * time.Minute

	// replaySweepInterval defines the minimum time between cleanup sweeps.
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
		seen: make(map[int64]time.Time),
	}
}

// MarkAndCheck records the given msgID and returns true if it is a replay
// (i.e. the same ID was seen within the replayTTL window).
func (r *ReplayDetector) MarkAndCheck(msgID int64) bool {
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Replay check: if we've seen this msgID recently, treat it as replay.
	if ts, ok := r.seen[msgID]; ok {
		if now.Sub(ts) <= replayTTL {
			return true
		}
		// If the entry is older than TTL, treat it as new and overwrite below.
	}

	// Periodic cleanup when the map grows large and enough time has passed.
	if len(r.seen) >= maxHubReplayEntries && now.Sub(r.lastSweep) > replaySweepInterval {
		for id, ts := range r.seen {
			if now.Sub(ts) > replayTTL {
				delete(r.seen, id)
			}
		}
		r.lastSweep = now
	}

	// Save the new (or refreshed) entry.
	r.seen[msgID] = now
	return false
}
