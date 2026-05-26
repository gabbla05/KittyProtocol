package api

import (
	"context"

	"github.com/quic-go/quic-go"
)

// StreamAdapter abstracts a bidirectional QUIC stream.
//
// This indirection allows KittyClient to remain transport-agnostic and
// easily testable (e.g. with in-memory or mock streams).
type StreamAdapter interface {
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Close() error

	// QUIC-specific cancellation hooks are kept here because the current
	// transport is QUIC. If we ever swap transport, a new adapter can
	// implement these as no-ops or map them to equivalent semantics.
	CancelRead(code quic.StreamErrorCode)
	CancelWrite(code quic.StreamErrorCode)
}

// ConnAdapter abstracts a QUIC connection.
type ConnAdapter interface {
	ConnectionState() quic.ConnectionState
	OpenStreamSync(ctx context.Context) (StreamAdapter, error)
	CloseWithError(code quic.ApplicationErrorCode, msg string) error
}

// quicConnAdapter is the real QUIC implementation of ConnAdapter.
type quicConnAdapter struct {
	conn *quic.Conn
}

func newQuicConnAdapter(conn *quic.Conn) *quicConnAdapter {
	return &quicConnAdapter{conn: conn}
}

func (a *quicConnAdapter) ConnectionState() quic.ConnectionState {
	return a.conn.ConnectionState()
}

func (a *quicConnAdapter) OpenStreamSync(ctx context.Context) (StreamAdapter, error) {
	s, err := a.conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	return &quicStreamAdapter{s: s}, nil
}

func (a *quicConnAdapter) CloseWithError(code quic.ApplicationErrorCode, msg string) error {
	return a.conn.CloseWithError(code, msg)
}

// quicStreamAdapter is the real QUIC implementation of StreamAdapter.
type quicStreamAdapter struct {
	s *quic.Stream
}

func (q *quicStreamAdapter) Read(p []byte) (int, error)  { return q.s.Read(p) }
func (q *quicStreamAdapter) Write(p []byte) (int, error) { return q.s.Write(p) }
func (q *quicStreamAdapter) Close() error                { return q.s.Close() }

func (q *quicStreamAdapter) CancelRead(code quic.StreamErrorCode) {
	q.s.CancelRead(code)
}

func (q *quicStreamAdapter) CancelWrite(code quic.StreamErrorCode) {
	q.s.CancelWrite(code)
}
