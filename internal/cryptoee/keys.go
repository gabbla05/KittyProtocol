package cryptoee

import (
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/hkdf"
)

// NOTE:
// For now we assume that Alice and Bob have agreed on a shared secret K_AB
// out-of-band (OOB). In this prototype we use a static 32-byte value as K_AB.
// In a real deployment this must be replaced with a proper key agreement
// mechanism and per-conversation secrets.
// ==========================================================================
var sharedSecretKAB = []byte("0123456789ABCDEF0123456789ABCDEF") // 32 bytes
// ==========================================================================

// DeriveKeys derives K_enc and K_mac from K_AB using HKDF-SHA256.
func DeriveKeys() (kEnc, kMac []byte, err error) {
	hkdfEnc := hkdf.New(sha256.New, sharedSecretKAB, nil, []byte("encryption"))
	hkdfMac := hkdf.New(sha256.New, sharedSecretKAB, nil, []byte("authentication"))

	kEnc = make([]byte, 32)
	kMac = make([]byte, 32)

	if _, err = io.ReadFull(hkdfEnc, kEnc); err != nil {
		return nil, nil, err
	}
	if _, err = io.ReadFull(hkdfMac, kMac); err != nil {
		return nil, nil, err
	}

	return kEnc, kMac, nil
}
