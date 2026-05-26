package protection

import (
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket limiter.
// It allows up to maxTokens operations per second for a given session.
// This protects the Hub from message flooding.
type RateLimiter struct {
	mu         sync.Mutex
	tokens     int
	maxTokens  int
	lastUpdate time.Time
}

// NewRateLimiter creates a new RateLimiter with the given per-second limit.
func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		tokens:     limit,
		maxTokens:  limit,
		lastUpdate: time.Now(),
	}
}

// Allow returns true if the operation is allowed at this moment.
// Tokens are refilled proportionally to elapsed time since the last check.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate)

	// Refill tokens based on elapsed time.
	refill := int(elapsed.Seconds() * float64(rl.maxTokens))
	if refill > 0 {
		rl.tokens += refill
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastUpdate = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}
