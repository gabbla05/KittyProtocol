package cryptoee

// Zeroize securely zeroes out the contents of the byte slice.
// This is a best-effort function to reduce the risk of sensitive data
// lingering in memory. Note that Go's garbage collector and memory management
// may still keep copies of the data, so this is not a guarantee.
func Zeroize(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
