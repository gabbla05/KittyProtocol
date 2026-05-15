package cryptoee

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	msgID := int64(123456789)
	target := "bob"
	plaintext := "Hello Bob, this is a secret message."

	payload, mac, err := EncryptAndMAC(msgID, target, plaintext)
	if err != nil {
		t.Fatalf("EncryptAndMAC failed: %v", err)
	}

	out, err := DecryptAndVerify(msgID, target, payload, mac)
	if err != nil {
		t.Fatalf("DecryptAndVerify failed: %v", err)
	}

	if out != plaintext {
		t.Fatalf("Decrypted plaintext mismatch.\nExpected: %s\nGot: %s", plaintext, out)
	}
}

func TestTamperedCiphertext(t *testing.T) {
	msgID := int64(42)
	target := "bob"
	plaintext := "Secret"

	payload, mac, err := EncryptAndMAC(msgID, target, plaintext)
	if err != nil {
		t.Fatalf("EncryptAndMAC failed: %v", err)
	}

	// Tamper with payload
	payloadBytes := []byte(payload)
	payloadBytes[len(payloadBytes)-1] ^= 0xFF
	tampered := string(payloadBytes)

	_, err = DecryptAndVerify(msgID, target, tampered, mac)
	if err == nil {
		t.Fatalf("Expected decryption failure after tampering, got nil error")
	}
}

func TestTamperedMAC(t *testing.T) {
	msgID := int64(42)
	target := "bob"
	plaintext := "Secret"

	payload, mac, err := EncryptAndMAC(msgID, target, plaintext)
	if err != nil {
		t.Fatalf("EncryptAndMAC failed: %v", err)
	}

	// Tamper with MAC
	macBytes := []byte(mac)
	macBytes[0] ^= 0xAA
	tamperedMAC := string(macBytes)

	_, err = DecryptAndVerify(msgID, target, payload, tamperedMAC)
	if err == nil {
		t.Fatalf("Expected HMAC verification failure, got nil error")
	}
}
