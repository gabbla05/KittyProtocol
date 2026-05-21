package api

import (
	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/internal/protection"
)

// Close gracefully shuts down the client:
//
//   - stops ping and receiver loops,
//   - closes the QUIC stream and connection,
//   - cancels the internal context,
//   - zeroizes encryption keys,
//   - resets replay detector and ACK manager,
//   - clears session state (target, lastFrame).
//
// This method is idempotent: calling it multiple times is safe.
func (c *KittyClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Stop background loops
	select {
	case <-c.stopPing:
	default:
		close(c.stopPing)
	}
	select {
	case <-c.stopRecv:
	default:
		close(c.stopRecv)
	}

	// Cancel context
	if c.cancel != nil {
		c.cancel()
	}

	// Close stream and connection
	if c.stream != nil {
		_ = c.stream.Close()
		c.stream = nil
	}
	if c.conn != nil {
		_ = c.conn.CloseWithError(0, "client closed")
		c.conn = nil
	}

	// Zeroize keys
	if c.kEnc != nil {
		cryptoee.Zeroize(c.kEnc)
		c.kEnc = nil
	}
	if c.kMac != nil {
		cryptoee.Zeroize(c.kMac)
		c.kMac = nil
	}

	// Reset session state
	c.target = ""
	c.lastFrame = nil
	c.replay = protection.NewReplayDetector()
	c.ackMgr = NewAckManager()

	c.state = StateDisconnected
}
