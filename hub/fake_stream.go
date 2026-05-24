package hub

import "io"

// Minimal fake stream implementing protection.Stream
type fakeStream struct {
	written [][]byte
}

func (f *fakeStream) Write(b []byte) (int, error) {
	f.written = append(f.written, append([]byte(nil), b...))
	return len(b), nil
}

func (f *fakeStream) Read(b []byte) (int, error) { return 0, io.EOF }
func (f *fakeStream) Close() error               { return nil }
