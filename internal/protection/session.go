package protection

import (
	"time"

	"github.com/quic-go/quic-go"
)

// Session represents a single authenticated user session on the Hub.
// It tracks last activity time, rate limiting, and a function to close the connection.
type Session struct {
	ID         string
	LastActive time.Time
	Limiter    *RateLimiter
	CloseFunc  func()
	Conn       *quic.Conn
	Stream     *quic.Stream
}

// NewSession creates a new Session for the given user and connection.
func NewSession(user string, conn *quic.Conn, stream *quic.Stream) *Session {
	return &Session{
		ID:         user,
		LastActive: time.Now(),
		Limiter:    NewRateLimiter(10), // 10 messages per second per user
		CloseFunc: func() {
			conn.CloseWithError(0x09, "Idle Timeout")
		},
		Conn:   conn,
		Stream: stream,
	}
}
