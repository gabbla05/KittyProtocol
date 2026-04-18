package clientutils

import (
	"fmt"
)

// MaxPlaintextSize is a conservative limit for plaintext message length.
// It is chosen to stay safely below the 2048-byte encrypted payload limit.
const MaxPlaintextSize = 1500

// TruncateMessage trims the input string so that it does not exceed MaxPlaintextSize.
// This is a temporary approximation until AEAD encryption is implemented.
func TruncateMessage(input string) string {
    if len(input) > MaxPlaintextSize {
        fmt.Printf("[System] Message too long (%d bytes). Truncating to %d bytes...\n", len(input), MaxPlaintextSize)
        return input[:MaxPlaintextSize]
    }
    return input
}
