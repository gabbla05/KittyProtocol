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
// We assume that Alice and Bob have agreed on a shared secret K_AB
// through some out-of-band mechanism (e.g., QR code, password, or prior exchange).
// This secret is then used to derive separate keys for encryption and authentication using HKDF.

// DeriveKeysFromSecret derives K_enc and K_mac from the provided shared secret
// using HKDF-SHA256.
//
// IMPORTANT:
// - 'secret' MUST already be a high-entropy shared secret K_AB agreed OOB
//   between the two parties (e.g. 32 random bytes).
// - This function does NOT perform password hashing. If you want to use
//   human-memorable passphrases, you MUST first run them through a KDF
//   suitable for passwords (e.g. Argon2, PBKDF2) and only then call this
//   function with the resulting key material.

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
