package protection

import (
	"time"

	"github.com/quic-go/quic-go"
)

// DefaultSessionRateLimit defines how many messages per second
// a single user session is allowed to send.
const DefaultSessionRateLimit = 10

// DefaultIdleCloseErrorCode is the application error code used when
// closing idle connections from the Hub side.
const DefaultIdleCloseErrorCode = 0x09

// Session represents a single authenticated user session on the Hub.
// It tracks last activity time, rate limiting, replay protection and
// provides a function to close the underlying QUIC connection.
type Session struct {
	ID         string
	LastActive time.Time
	Limiter    *RateLimiter
	CloseFunc  func()
	Conn       *quic.Conn
	Stream     *quic.Stream
	Target     string
	Replay     *ReplayDetector
}

// NewSession creates a new Session for the given user and connection.
func NewSession(user string, conn *quic.Conn, stream *quic.Stream) *Session {
	return &Session{
		ID:         user,
		LastActive: time.Now(),
		Limiter:    NewRateLimiter(DefaultSessionRateLimit),
		CloseFunc: func() {
			// Application error code 0x09 is used as "Idle Timeout" in this project.
			conn.CloseWithError(DefaultIdleCloseErrorCode, "Idle Timeout")
		},
		Conn:   conn,
		Stream: stream,
		Replay: NewReplayDetector(),
	}
}
