package protection

import (
	"testing"
	"time"
)

// TestRateLimiterInitialTokens verifies that a new limiter starts with
// maxTokens available and allows exactly that many operations.
func TestRateLimiterInitialTokens(t *testing.T) {
	rl := NewRateLimiter(10)

	for i := 0; i < 10; i++ {
		if !rl.Allow() {
			t.Fatalf("expected token %d to be allowed", i)
		}
	}

	if rl.Allow() {
		t.Fatalf("expected limiter to block after tokens are exhausted")
	}
}

// TestRateLimiterRefill ensures that tokens are refilled proportionally
// to elapsed time since the last Allow() call.
func TestRateLimiterRefill(t *testing.T) {
	rl := NewRateLimiter(10)

	// Exhaust tokens
	for i := 0; i < 10; i++ {
		rl.Allow()
	}

	if rl.Allow() {
		t.Fatalf("expected limiter to block after exhaustion")
	}

	// Wait long enough for at least one token to refill
	time.Sleep(150 * time.Millisecond)

	if !rl.Allow() {
		t.Fatalf("expected limiter to allow after refill")
	}
}
