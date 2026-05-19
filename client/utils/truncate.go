package clientutils

// MaxPlaintextSize is a conservative limit for plaintext message length.
// It is chosen to stay safely below the encrypted payload limit after AEAD overhead.
const MaxPlaintextSize = 1500

// TruncateMessage trims the input string so that it does not exceed MaxPlaintextSize.
// It returns the truncated string. No logging is performed here to keep the utility pure.
func TruncateMessage(input string) string {
	if len(input) > MaxPlaintextSize {
		return input[:MaxPlaintextSize]
	}
	return input
}
