package protection

import (
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket limiter.
// It allows up to maxTokens operations per second.
type RateLimiter struct {
	tokens     int
	maxTokens  int
	lastUpdate time.Time
	mu         sync.Mutex
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
// It refills tokens proportionally to elapsed time since the last check.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate)
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

// AuthTimer wraps a time.Timer used for the 20-second AUTH timeout.
type AuthTimer struct {
	timer *time.Timer
}

// StartAuthTimer starts a 20-second timer that calls onTimeout when it fires.
func StartAuthTimer(onTimeout func()) *AuthTimer {
	return &AuthTimer{
		timer: time.AfterFunc(20*time.Second, onTimeout),
	}
}

// Stop cancels the AUTH timer if it is still running.
func (at *AuthTimer) Stop() {
	if at.timer != nil {
		at.timer.Stop()
	}
}
