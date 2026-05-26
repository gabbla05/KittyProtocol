package protection

// Stream is a minimal abstraction over a bidirectional transport stream.
// It is intentionally small to keep Hub logic independent from quic-go.
type Stream interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Close() error
}
