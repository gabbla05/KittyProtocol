package protection

import (
	"testing"
	"time"
)

// Test that limiter allows up to N operations immediately.
func TestRateLimiterInitialTokens(t *testing.T) {
    rl := NewRateLimiter(10)

    for i := 0; i < 10; i++ {
        if !rl.Allow() {
            t.Fatalf("expected token %d to be allowed", i)
        }
    }

    if rl.Allow() {
        t.Fatalf("expected limiter to block after tokens exhausted")
    }
}

// Test that tokens refill over time.
func TestRateLimiterRefill(t *testing.T) {
    rl := NewRateLimiter(10)

    // Exhaust tokens
    for i := 0; i < 10; i++ {
        rl.Allow()
    }

    if rl.Allow() {
        t.Fatalf("expected limiter to block after exhaustion")
    }

    // Wait long enough for refill
    time.Sleep(150 * time.Millisecond)

    if !rl.Allow() {
        t.Fatalf("expected limiter to refill tokens after time")
    }
}
