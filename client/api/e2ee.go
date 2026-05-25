package api

import (
	"errors"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
)

// SetSharedSecretForPeer derives encryption and MAC keys from the shared secret
// and stores them for a specific peer (logical username).
//
// SECURITY:
//   - Keys are derived using HKDF-SHA256 in internal/cryptoee.
//   - Keys are zeroized and cleared in KittyClient.Close().
//   - Keys are kept only in memory and never written to disk.
func (c *KittyClient) SetSharedSecretForPeer(peer string, secret []byte) error {
	if peer == "" {
		return errors.New("peer cannot be empty")
	}
	if len(secret) == 0 {
		return errors.New("secret cannot be empty")
	}

	kEnc, kMac, err := cryptoee.DeriveKeysFromSecret(secret)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.peerKeys == nil {
		c.peerKeys = make(map[string]peerKeys)
	}
	c.peerKeys[peer] = peerKeys{
		kEnc: kEnc,
		kMac: kMac,
	}
	return nil
}

func (c *KittyClient) HasSharedSecret(peer string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.peerKeys == nil {
		return false
	}
	_, ok := c.peerKeys[peer]
	return ok
}
