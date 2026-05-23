package cryptoee

import (
	"bytes"
	"encoding/base64"
	"testing"
)

// --- Helpers ---

func mustKeys(t *testing.T) ([]byte, []byte) {
	secret := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA") // 32 bytes
	kEnc, kMac, err := DeriveKeysFromSecret(secret)
	if err != nil {
		t.Fatalf("DeriveKeysFromSecret failed: %v", err)
	}
	return kEnc, kMac
}

// --- Core encryption/decryption tests ---

func TestEncryptDecryptRoundtrip(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(123456)
	target := "Bob"
	plaintext := "Hello Bob, this is a secret message."

	payload, mac, err := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	out, err := DecryptAndVerifyWithKeys(msgID, target, payload, mac, kEnc, kMac)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if out != plaintext {
		t.Fatalf("plaintext mismatch: expected %q, got %q", plaintext, out)
	}
}

// --- Tampering tests ---

func TestTamperedCiphertext(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(42)
	target := "bob"
	plaintext := "Secret"

	payload, mac, err := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Tamper payload
	raw, _ := base64.StdEncoding.DecodeString(payload)
	raw[len(raw)-1] ^= 0xFF
	tampered := base64.StdEncoding.EncodeToString(raw)

	_, err = DecryptAndVerifyWithKeys(msgID, target, tampered, mac, kEnc, kMac)
	if err == nil {
		t.Fatalf("expected decryption failure after tampering")
	}
}

func TestTamperedMAC(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(42)
	target := "bob"
	plaintext := "Secret"

	payload, mac, err := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Tamper MAC
	raw, _ := base64.StdEncoding.DecodeString(mac)
	raw[0] ^= 0xAA
	tampered := base64.StdEncoding.EncodeToString(raw)

	_, err = DecryptAndVerifyWithKeys(msgID, target, payload, tampered, kEnc, kMac)
	if err == nil {
		t.Fatalf("expected HMAC verification failure")
	}
}

// --- Wrong msgID / wrong target ---

func TestWrongMsgID(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(100)
	target := "bob"
	plaintext := "Hello"

	payload, mac, err := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = DecryptAndVerifyWithKeys(msgID+1, target, payload, mac, kEnc, kMac)
	if err == nil {
		t.Fatalf("expected HMAC failure for wrong msgID")
	}
}

func TestWrongTarget(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(200)
	target := "bob"
	plaintext := "Hello"

	payload, mac, err := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = DecryptAndVerifyWithKeys(msgID, "alice", payload, mac, kEnc, kMac)
	if err == nil {
		t.Fatalf("expected HMAC failure for wrong target")
	}
}

// --- Base64 / payload errors ---

func TestPayloadTooShort(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(300)
	target := "bob"

	shortPayload := base64.StdEncoding.EncodeToString([]byte{1, 2, 3})

	_, err := DecryptAndVerifyWithKeys(msgID, target, shortPayload, "AAAA", kEnc, kMac)
	if err == nil {
		t.Fatalf("expected error for too short payload")
	}
}

func TestInvalidBase64Payload(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(400)
	target := "bob"

	_, err := DecryptAndVerifyWithKeys(msgID, target, "!!!notbase64!!!", "AAAA", kEnc, kMac)
	if err == nil {
		t.Fatalf("expected base64 decode error")
	}
}

func TestInvalidBase64MAC(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(500)
	target := "bob"

	payload, _, err := EncryptAndMACWithKeys(msgID, target, "Hello", kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = DecryptAndVerifyWithKeys(msgID, target, payload, "!!!notbase64!!!", kEnc, kMac)
	if err == nil {
		t.Fatalf("expected base64 decode error for MAC")
	}
}

// --- HKDF tests ---

func TestDeriveKeysDeterministic(t *testing.T) {
	secret := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	k1Enc, k1Mac, err := DeriveKeysFromSecret(secret)
	if err != nil {
		t.Fatalf("DeriveKeysFromSecret failed: %v", err)
	}

	k2Enc, k2Mac, err := DeriveKeysFromSecret(secret)
	if err != nil {
		t.Fatalf("DeriveKeysFromSecret failed: %v", err)
	}

	if !bytes.Equal(k1Enc, k2Enc) || !bytes.Equal(k1Mac, k2Mac) {
		t.Fatalf("HKDF must be deterministic for same secret")
	}
}

func TestCanonicalization(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(600)
	plaintext := "Hello"

	// Different forms of same target
	payload1, mac1, _ := EncryptAndMACWithKeys(msgID, " Bob ", plaintext, kEnc, kMac)
	payload2, mac2, _ := EncryptAndMACWithKeys(msgID, "bob", plaintext, kEnc, kMac)

	if payload1 == payload2 && mac1 == mac2 {
		// This is OK — canonicalization makes them equivalent
		return
	}

	// But decryption must work for both
	_, err := DecryptAndVerifyWithKeys(msgID, "bob", payload1, mac1, kEnc, kMac)
	if err != nil {
		t.Fatalf("canonicalization failed: %v", err)
	}
}
