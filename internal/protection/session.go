package protection

import (
	"time"

	"github.com/quic-go/quic-go"
)

// Session represents a single authenticated user session on the Hub.
// It tracks last activity time, rate limiting, replay protection and
// provides a function to close the underlying QUIC connection.
type Session struct {
	ID         string
	LastActive time.Time
	Limiter    *RateLimiter
	CloseFunc  func()
	Conn       *quic.Conn
	Stream     Stream
	Replay     *ReplayDetector
}

// NewSession creates a new Session for the given user and connection.
func NewSession(user string, conn *quic.Conn, stream Stream) *Session {
	return &Session{
		ID:         user,
		LastActive: time.Now(),
		Limiter:    NewRateLimiter(DefaultSessionRateLimit),
		CloseFunc: func() {
			conn.CloseWithError(DefaultIdleCloseErrorCode, "Idle Timeout")
		},
		Conn:   conn,
		Stream: stream,
		Replay: NewReplayDetector(),
	}
}
