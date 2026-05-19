package cryptoee

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

// buildMACInput = cipher || msg_id || target.
func buildMACInput(cipher []byte, msgID int64, target string) []byte {
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(msgID))

	out := make([]byte, 0, len(cipher)+len(msg)+len(target))
	out = append(out, cipher...)
	out = append(out, msg...)
	out = append(out, []byte(target)...)
	return out
}

// EncryptAndMACWithKeys encrypts plaintext using AES-GCM and computes HMAC-SHA256
// over cipher || msg_id || target.
//
// SECURITY NOTE:
// - For a given K_enc, (msgID) MUST NOT repeat. Reusing the same (K_enc, nonce)
//   pair breaks AES-GCM security guarantees.
// - In this prototype, msgID is derived from time.Now().UnixMilli() on the client
//   side. In a production system, you should use a strictly monotonic counter
//   or another mechanism that guarantees uniqueness per key.

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
	binary.BigEndian.PutUint64(nonce[nonceSize-8:], uint64(msgID))

	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), nil)

	macInput := buildMACInput(ciphertext, msgID, target)
	h := hmac.New(sha256.New, kMac)
	h.Write(macInput)
	mac := h.Sum(nil)

	raw := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(raw),
		base64.StdEncoding.EncodeToString(mac),
		nil
}
