package cryptoee

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
)

// buildMACInput builds the input to HMAC as:
//
//	cipher || msg_id (big-endian uint64) || canonical_target
func buildMACInput(cipher []byte, msgID int64, target string) []byte {
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(msgID))

	canon := canonicalizeTarget(target)

	out := make([]byte, 0, len(cipher)+len(msg)+len(canon))
	out = append(out, cipher...)
	out = append(out, msg...)
	out = append(out, canon...)
	return out
}

// canonicalizeTarget normalizes the target string to a stable form
// to avoid MAC mismatches due to case or whitespace differences.
func canonicalizeTarget(t string) string {
	// Lowercase + trim is sufficient for this protocol.
	// If needed, Unicode normalization (NFC) can be added later.
	return strings.ToLower(strings.TrimSpace(t))
}

// EncryptAndMACWithKeys encrypts plaintext using AES-GCM and computes HMAC-SHA256
// over cipher || msg_id || canonical_target.
//
// The function returns:
//   - payloadB64: base64-encoded nonce || ciphertext
//   - macB64:     base64-encoded HMAC-SHA256
func EncryptAndMACWithKeys(msgID int64, target, plaintext string, kEnc, kMac []byte) (string, string, error) {
	block, err := aes.NewCipher(kEnc)
	if err != nil {
		return "", "", fmt.Errorf("[Encrypt]: AES cipher error: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", fmt.Errorf("[Encrypt]: AEAD error: %w", err)
	}

	nonceSize := aead.NonceSize()

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return "", "", fmt.Errorf("[Encrypt]: nonce generation error: %w", err)
	}

	// AAD binds msgID, canonical target and a format version to the ciphertext.
	aad := []byte(fmt.Sprintf("msgid=%d;target=%s;v=%d", msgID, canonicalizeTarget(target), aadFormatVersion))

	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), aad)

	macInput := buildMACInput(ciphertext, msgID, target)
	h := hmac.New(sha256.New, kMac)
	h.Write(macInput)
	mac := h.Sum(nil)

	raw := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(raw),
		base64.StdEncoding.EncodeToString(mac),
		nil
}
