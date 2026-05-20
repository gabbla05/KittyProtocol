package api

import "github.com/gabbla05/KittyProtocol/internal/cryptoee"

// SetSharedSecret derives encryption and MAC keys from the shared secret
// and stores them in the client.
//
// SECURITY:
//   - Keys are derived using a KDF implemented in internal/cryptoee.
//   - Keys are zeroized and cleared in KittyClient.Close().
//   - Keys are kept only in memory and never written to disk.
func (c *KittyClient) SetSharedSecret(secret []byte) error {
	kEnc, kMac, err := cryptoee.DeriveKeysFromSecret(secret)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.kEnc = kEnc
	c.kMac = kMac
	return nil
}
