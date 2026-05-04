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

// EncryptAndMAC encrypts plaintext using AEAD (AES-GCM) and computes HMAC.
// Returns BASE64(nonce||ciphertext) and BASE64(mac).
func EncryptAndMAC(msgID int64, target, plaintext string) (string, string, error) {
	kEnc, kMac, err := DeriveKeys()
	if err != nil {
		return "", "", err
	}

	block, err := aes.NewCipher(kEnc)
	if err != nil {
		return "", "", fmt.Errorf("AES cipher error: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", fmt.Errorf("AEAD error: %w", err)
	}

	nonceSize := aead.NonceSize()
	nonce := make([]byte, nonceSize)

	// Encode msg_id into last 8 bytes of nonce
	binary.BigEndian.PutUint64(nonce[nonceSize-8:], uint64(msgID))

	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), nil)

	// HMAC(cipher || msg_id || target)
	macInput := buildMACInput(ciphertext, msgID, target)
	h := hmac.New(sha256.New, kMac)
	h.Write(macInput)
	mac := h.Sum(nil)

	// payload = BASE64(nonce||cipher)
	raw := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(raw),
		base64.StdEncoding.EncodeToString(mac),
		nil
}

// buildMACInput = cipher || msg_id || target
func buildMACInput(cipher []byte, msgID int64, target string) []byte {
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(msgID))

	out := make([]byte, 0, len(cipher)+8+len(target))
	out = append(out, cipher...)
	out = append(out, msg...)
	out = append(out, []byte(target)...)
	return out
}
