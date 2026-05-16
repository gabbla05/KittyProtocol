package cryptoee

import (
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/hkdf"
)

// Key sizes used for both encryption and MAC keys.
const (
	KeySizeBytes = 32
)

// NOTE:
// For now we assume that Alice and Bob have agreed on a shared secret K_AB
// out-of-band (OOB). In this prototype we use a static 32-byte value as K_AB.
// In a real deployment this MUST be replaced with a proper key agreement
// mechanism (e.g. X25519) and per-conversation secrets.
// ============================================
// DO NOT USE THIS STATIC SECRET IN PRODUCTION.
// ============================================
var sharedSecretKAB = []byte("0123456789ABCDEF0123456789ABCDEF") // 32 bytes
// ============================================================================

// DeriveKeys derives K_enc and K_mac from the static K_AB using HKDF-SHA256.
// This is a convenience wrapper for the prototype. In a real system you should
// use DeriveKeysFromSecret with a per-session/per-conversation secret.
func DeriveKeys() (kEnc, kMac []byte, err error) {
	return DeriveKeysFromSecret(sharedSecretKAB)
}

// DeriveKeysFromSecret derives K_enc and K_mac from the provided secret using HKDF-SHA256.
func DeriveKeysFromSecret(secret []byte) (kEnc, kMac []byte, err error) {
	hkdfEnc := hkdf.New(sha256.New, secret, nil, []byte("encryption"))
	hkdfMac := hkdf.New(sha256.New, secret, nil, []byte("authentication"))

	kEnc = make([]byte, KeySizeBytes)
	kMac = make([]byte, KeySizeBytes)

	if _, err = io.ReadFull(hkdfEnc, kEnc); err != nil {
		return nil, nil, err
	}
	if _, err = io.ReadFull(hkdfMac, kMac); err != nil {
		return nil, nil, err
	}

	return kEnc, kMac, nil
}
