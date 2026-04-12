package protection

import (
	"sync"
	"time"
)

type RateLimiter struct {
	tokens     int
	maxTokens  int
	lastUpdate time.Time
	mu         sync.Mutex
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		tokens:     limit,
		maxTokens:  limit,
		lastUpdate: time.Now(),
	}
}

// Allow sprawdza, czy użytkownik może wysłać wiadomość.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// Regeneracja tokenów (1 na 100ms dla limitu 10/s)
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
