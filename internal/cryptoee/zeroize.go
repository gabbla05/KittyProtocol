package cryptoee

// Zeroize attempts to securely overwrite the contents of the byte slice.
// This is a best-effort mitigation to reduce the lifetime of sensitive data
// in memory. Go's runtime and garbage collector may still keep copies, so
// this is not a hard security guarantee.
func Zeroize(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
