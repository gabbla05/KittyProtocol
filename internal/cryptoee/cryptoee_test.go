package cryptoee

import (
	"bytes"
	"encoding/base64"
	"testing"
)

// mustKeys is a helper that derives deterministic test keys from a fixed secret.
// This ensures reproducible encryption/decryption results across test runs.
func mustKeys(t *testing.T) ([]byte, []byte) {
	t.Helper()

	secret := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA") // 32 bytes
	kEnc, kMac, err := DeriveKeysFromSecret(secret)
	if err != nil {
		t.Fatalf("DeriveKeysFromSecret failed: %v", err)
	}
	return kEnc, kMac
}

//
// ────────────────────────────────────────────────────────────────────────────────
//   CORE ROUNDTRIP TESTS
// ────────────────────────────────────────────────────────────────────────────────
//

// TestEncryptDecryptRoundtrip verifies that encryption followed by decryption
// returns the original plaintext without modification.
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

// TestFullCryptoFlow acts as an integration test verifying the end-to-end
// encryption, signing, verification, and decryption process for a clean data flow.
func TestFullCryptoFlow(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(999)
	target := "charlie"
	plaintext := "Integration test message"

	payload, mac, err := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	out, err := DecryptAndVerifyWithKeys(msgID, target, payload, mac, kEnc, kMac)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if out != plaintext {
		t.Fatalf("plaintext mismatch")
	}
}

//
// ────────────────────────────────────────────────────────────────────────────────
//   NONCE / RANDOMIZATION TESTS
// ────────────────────────────────────────────────────────────────────────────────
//

// TestDifferentNonceProducesDifferentCiphertext ensures that subsequent encryption
// operations of the same plaintext yield distinct ciphertexts and MACs due to unique nonces.
func TestDifferentNonceProducesDifferentCiphertext(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(777)
	target := "bob"
	plaintext := "Hello"

	p1, m1, _ := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	p2, m2, _ := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)

	if p1 == p2 || m1 == m2 {
		t.Fatalf("encryption must be randomized (nonce must differ)")
	}
}

//
// ────────────────────────────────────────────────────────────────────────────────
//   TAMPERING TESTS
// ────────────────────────────────────────────────────────────────────────────────
//

// TestTamperedCiphertext ensures that modifying even a single byte of the
// ciphertext results in decryption failure.
func TestTamperedCiphertext(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(42)
	target := "bob"
	plaintext := "Secret"

	payload, mac, err := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	raw, _ := base64.StdEncoding.DecodeString(payload)
	raw[len(raw)-1] ^= 0xFF // flip last byte
	tampered := base64.StdEncoding.EncodeToString(raw)

	if _, err := DecryptAndVerifyWithKeys(msgID, target, tampered, mac, kEnc, kMac); err == nil {
		t.Fatalf("expected decryption failure after ciphertext tampering")
	}
}

// TestTamperedMAC ensures that modifying the MAC results in HMAC verification failure.
func TestTamperedMAC(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(42)
	target := "bob"
	plaintext := "Secret"

	payload, mac, err := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	raw, _ := base64.StdEncoding.DecodeString(mac)
	raw[0] ^= 0xAA // flip first byte
	tampered := base64.StdEncoding.EncodeToString(raw)

	if _, err := DecryptAndVerifyWithKeys(msgID, target, payload, tampered, kEnc, kMac); err == nil {
		t.Fatalf("expected HMAC verification failure")
	}
}

//
// ────────────────────────────────────────────────────────────────────────────────
//   WRONG PARAMETERS (msgID / target)
// ────────────────────────────────────────────────────────────────────────────────
//

// TestWrongMsgID verifies that using a different msgID breaks HMAC verification.
func TestWrongMsgID(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(100)
	target := "bob"
	plaintext := "Hello"

	payload, mac, _ := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)

	if _, err := DecryptAndVerifyWithKeys(msgID+1, target, payload, mac, kEnc, kMac); err == nil {
		t.Fatalf("expected HMAC failure for wrong msgID")
	}
}

// TestWrongTarget verifies that canonicalized target mismatch breaks HMAC verification.
func TestWrongTarget(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(200)
	target := "bob"
	plaintext := "Hello"

	payload, mac, _ := EncryptAndMACWithKeys(msgID, target, plaintext, kEnc, kMac)

	if _, err := DecryptAndVerifyWithKeys(msgID, "alice", payload, mac, kEnc, kMac); err == nil {
		t.Fatalf("expected HMAC failure for wrong target")
	}
}

//
// ────────────────────────────────────────────────────────────────────────────────
//   BASE64 / PAYLOAD VALIDATION
// ────────────────────────────────────────────────────────────────────────────────
//

// TestPayloadTooShort ensures that payloads shorter than nonce size are rejected.
func TestPayloadTooShort(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(300)
	target := "bob"

	shortPayload := base64.StdEncoding.EncodeToString([]byte{1, 2, 3})

	if _, err := DecryptAndVerifyWithKeys(msgID, target, shortPayload, "AAAA", kEnc, kMac); err == nil {
		t.Fatalf("expected error for too short payload")
	}
}

// TestInvalidBase64Payload ensures that invalid base64 payloads are rejected.
func TestInvalidBase64Payload(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(400)
	target := "bob"

	if _, err := DecryptAndVerifyWithKeys(msgID, target, "!!!notbase64!!!", "AAAA", kEnc, kMac); err == nil {
		t.Fatalf("expected base64 decode error")
	}
}

// TestInvalidBase64MAC ensures that invalid base64 MACs are rejected.
func TestInvalidBase64MAC(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(500)
	target := "bob"

	payload, _, _ := EncryptAndMACWithKeys(msgID, target, "Hello", kEnc, kMac)

	if _, err := DecryptAndVerifyWithKeys(msgID, target, payload, "!!!notbase64!!!", kEnc, kMac); err == nil {
		t.Fatalf("expected base64 decode error for MAC")
	}
}

//
// ────────────────────────────────────────────────────────────────────────────────
//   HKDF TESTS
// ────────────────────────────────────────────────────────────────────────────────
//

// TestDeriveKeysDeterministic verifies that HKDF produces deterministic output
// for the same input secret.
func TestDeriveKeysDeterministic(t *testing.T) {
	secret := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	k1Enc, k1Mac, _ := DeriveKeysFromSecret(secret)
	k2Enc, k2Mac, _ := DeriveKeysFromSecret(secret)

	if !bytes.Equal(k1Enc, k2Enc) || !bytes.Equal(k1Mac, k2Mac) {
		t.Fatalf("HKDF must be deterministic for same secret")
	}
}

//
// ────────────────────────────────────────────────────────────────────────────────
//   CANONICALIZATION TESTS
// ────────────────────────────────────────────────────────────────────────────────
//

// TestCanonicalization verifies that different textual forms of the same target
// decrypt correctly due to canonicalization.
func TestCanonicalization(t *testing.T) {
	kEnc, kMac := mustKeys(t)

	msgID := int64(600)
	plaintext := "Hello"

	payload1, mac1, _ := EncryptAndMACWithKeys(msgID, " Bob ", plaintext, kEnc, kMac)
	payload2, mac2, _ := EncryptAndMACWithKeys(msgID, "bob", plaintext, kEnc, kMac)

	// Decryption must succeed regardless of whitespace/case differences.
	if _, err := DecryptAndVerifyWithKeys(msgID, "bob", payload1, mac1, kEnc, kMac); err != nil {
		t.Fatalf("canonicalization failed for payload1: %v", err)
	}

	if _, err := DecryptAndVerifyWithKeys(msgID, " Bob ", payload2, mac2, kEnc, kMac); err != nil {
		t.Fatalf("canonicalization failed for payload2: %v", err)
	}
}

//
// ────────────────────────────────────────────────────────────────────────────────
//   MEMORY MANAGEMENT TESTS
// ────────────────────────────────────────────────────────────────────────────────
//

// TestZeroize verifies that the Zeroize function securely overwrites the provided
// byte slice with zeroes to clear sensitive data from memory.
func TestZeroize(t *testing.T) {
	b := []byte("sensitive-data")
	Zeroize(b)

	for _, v := range b {
		if v != 0 {
			t.Fatalf("Zeroize failed: memory not overwritten")
		}
	}
}
