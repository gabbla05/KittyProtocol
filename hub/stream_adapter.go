package hub

import (
	"github.com/quic-go/quic-go"
)

// quicStreamAdapter adapts *quic.Stream to protection.Stream.
type quicStreamAdapter struct {
	s *quic.Stream
}

func (q *quicStreamAdapter) Write(b []byte) (int, error) { return q.s.Write(b) }
func (q *quicStreamAdapter) Read(b []byte) (int, error)  { return q.s.Read(b) }
func (q *quicStreamAdapter) Close() error                { return q.s.Close() }
