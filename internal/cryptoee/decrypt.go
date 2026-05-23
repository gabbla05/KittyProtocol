package cryptoee

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// DecryptAndVerifyWithKeys verifies HMAC and decrypts the payload using AES-GCM.
// It expects:
//   - payloadB64: base64-encoded nonce || ciphertext
//   - macB64:     base64-encoded HMAC-SHA256 over cipher || msg_id || canonical_target.
func DecryptAndVerifyWithKeys(msgID int64, target, payloadB64, macB64 string, kEnc, kMac []byte) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(payloadB64)
	if err != nil {
		return "", fmt.Errorf("payload decode error: %w", err)
	}

	macRaw, err := base64.StdEncoding.DecodeString(macB64)
	if err != nil {
		return "", fmt.Errorf("mac decode error: %w", err)
	}

	block, err := aes.NewCipher(kEnc)
	if err != nil {
		return "", fmt.Errorf("AES cipher error: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("AEAD error: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(raw) < nonceSize {
		return "", fmt.Errorf("payload too short")
	}

	nonce := raw[:nonceSize]
	ciphertext := raw[nonceSize:]

	aad := []byte(fmt.Sprintf("msgid=%d;target=%s;v=%d", msgID, canonicalizeTarget(target), aadFormatVersion))

	macInput := buildMACInput(ciphertext, msgID, target)
	h := hmac.New(sha256.New, kMac)
	h.Write(macInput)
	expected := h.Sum(nil)

	if !hmac.Equal(macRaw, expected) {
		return "", fmt.Errorf("HMAC verification failed")
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}
