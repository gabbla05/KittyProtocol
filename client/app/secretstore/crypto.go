package secretstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// deriveKey normalizes masterKey to 32 bytes (AES‑256) using SHA‑256.
// In the future this can be replaced with PBKDF2/Argon2 without changing callers.
func deriveKey(masterKey []byte) []byte {
	sum := sha256.Sum256(masterKey)
	return sum[:]
}

// encrypt encrypts plaintext using AES‑GCM(deriveKey(masterKey)).
// Returns base64(nonce || ciphertext).
func encrypt(masterKey, plaintext []byte) (string, error) {
	if len(masterKey) == 0 {
		return "", ErrEmptyMasterKey
	}

	key := deriveKey(masterKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	out := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(out), nil
}

// decrypt decrypts base64(nonce || ciphertext) using AES‑GCM(deriveKey(masterKey)).
func decrypt(masterKey []byte, enc string) ([]byte, error) {
	if len(masterKey) == 0 {
		return nil, ErrEmptyMasterKey
	}

	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	key := deriveKey(masterKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(raw) < gcm.NonceSize() {
		return nil, ErrShortCipher
	}

	nonce := raw[:gcm.NonceSize()]
	ciphertext := raw[gcm.NonceSize():]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
