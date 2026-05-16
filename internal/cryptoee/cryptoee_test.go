package cryptoee

import (
	"bytes"
	"encoding/base64"
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

func TestWrongMsgID(t *testing.T) {
	msgID := int64(100)
	target := "bob"
	plaintext := "Hello"

	payload, mac, err := EncryptAndMAC(msgID, target, plaintext)
	if err != nil {
		t.Fatalf("EncryptAndMAC failed: %v", err)
	}

	// używamy innego msgID przy deszyfrowaniu
	_, err = DecryptAndVerify(msgID+1, target, payload, mac)
	if err == nil {
		t.Fatalf("Expected HMAC failure for wrong msgID")
	}
}

func TestWrongTarget(t *testing.T) {
	msgID := int64(200)
	target := "bob"
	plaintext := "Hello"

	payload, mac, err := EncryptAndMAC(msgID, target, plaintext)
	if err != nil {
		t.Fatalf("EncryptAndMAC failed: %v", err)
	}

	// zmieniamy target przy deszyfrowaniu
	_, err = DecryptAndVerify(msgID, "alice", payload, mac)
	if err == nil {
		t.Fatalf("Expected HMAC failure for wrong target")
	}
}

func TestPayloadTooShort(t *testing.T) {
	msgID := int64(300)
	target := "bob"

	// payload krótszy niż nonce
	shortPayload := base64.StdEncoding.EncodeToString([]byte{1, 2, 3})

	_, err := DecryptAndVerify(msgID, target, shortPayload, "AAAA")
	if err == nil {
		t.Fatalf("Expected error for too short payload")
	}
}

func TestInvalidBase64Payload(t *testing.T) {
	msgID := int64(400)
	target := "bob"

	_, err := DecryptAndVerify(msgID, target, "!!!notbase64!!!", "AAAA")
	if err == nil {
		t.Fatalf("Expected base64 decode error")
	}
}

func TestInvalidBase64MAC(t *testing.T) {
	msgID := int64(500)
	target := "bob"

	payload, _, err := EncryptAndMAC(msgID, target, "Hello")
	if err != nil {
		t.Fatalf("EncryptAndMAC failed: %v", err)
	}

	_, err = DecryptAndVerify(msgID, target, payload, "!!!notbase64!!!")
	if err == nil {
		t.Fatalf("Expected base64 decode error for MAC")
	}
}

func TestDeriveKeysDeterministic(t *testing.T) {
	k1Enc, k1Mac, err := DeriveKeys()
	if err != nil {
		t.Fatalf("DeriveKeys failed: %v", err)
	}

	k2Enc, k2Mac, err := DeriveKeys()
	if err != nil {
		t.Fatalf("DeriveKeys failed: %v", err)
	}

	if !bytes.Equal(k1Enc, k2Enc) || !bytes.Equal(k1Mac, k2Mac) {
		t.Fatalf("DeriveKeys must be deterministic for static secret")
	}
}

func TestDeriveKeysFromSecret(t *testing.T) {
	secret := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA") // 32 bytes

	kEnc, kMac, err := DeriveKeysFromSecret(secret)
	if err != nil {
		t.Fatalf("DeriveKeysFromSecret failed: %v", err)
	}

	if len(kEnc) != KeySizeBytes || len(kMac) != KeySizeBytes {
		t.Fatalf("Derived keys must be %d bytes", KeySizeBytes)
	}
}
