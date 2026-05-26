package cryptoee

// KeySizeBytes defines the size (in bytes) of derived encryption and MAC keys.
// AES-256 + HMAC-SHA256 both use 32-byte keys.
const KeySizeBytes = 32

// aadFormatVersion is a protocol version tag embedded in the AAD string.
// Bumping this value invalidates old ciphertexts at the AAD level.
const aadFormatVersion = 1
