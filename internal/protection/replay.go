package protection

import (
	"sync"
	"time"
)

const (
	maxHubReplayEntries = 100_000
	replayTTL           = 2 * time.Minute
)

type ReplayDetector struct {
	mu        sync.Mutex
	seen      map[int64]time.Time
	lastSweep time.Time
}

func NewReplayDetector() *ReplayDetector {
	return &ReplayDetector{
		seen: make(map[int64]time.Time),
	}
}

func (r *ReplayDetector) MarkAndCheck(msgID int64) bool {
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	// replay check
	if ts, ok := r.seen[msgID]; ok {
		// if still in TTL -> replay
		if now.Sub(ts) <= replayTTL {
			return true
		}
		// if defunct treat like new
	}

	// cleaning in time interval
	if len(r.seen) >= maxHubReplayEntries && now.Sub(r.lastSweep) > 5*time.Second {
		for id, ts := range r.seen {
			if now.Sub(ts) > replayTTL {
				delete(r.seen, id)
			}
		}
		r.lastSweep = now
	}

	// save new
	r.seen[msgID] = now
	return false
}
