package api

import (
	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

// Close gracefully shuts down the client and securely clears all sensitive data.
//
// Behavior:
//   - stops ping and receiver loops,
//   - closes the QUIC stream and connection,
//   - cancels the internal context,
//   - zeroizes all per‑peer E2EE keys,
//   - resets replay detector and ACK manager,
//   - clears session metadata.
//
// This method is idempotent.
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

	// Forcefully interrupt any blocking Read/Write
	if c.stream != nil {
		// QUIC-specific error code 0 is fine here; semantics are "no specific error".
		c.stream.CancelRead(quic.StreamErrorCode(0))
		c.stream.CancelWrite(quic.StreamErrorCode(0))
	}

	// Close stream
	if c.stream != nil {
		_ = c.stream.Close()
		c.stream = nil
	}

	// Close connection
	if c.conn != nil {
		_ = c.conn.CloseWithError(0, "client closed")
		c.conn = nil
	}

	// Zeroize all per‑peer E2EE keys
	for peer, pk := range c.peerKeys {
		if pk.kEnc != nil {
			cryptoee.Zeroize(pk.kEnc)
		}
		if pk.kMac != nil {
			cryptoee.Zeroize(pk.kMac)
		}
		delete(c.peerKeys, peer)
	}
	c.peerKeys = nil

	// Reset session state
	c.user = ""
	c.target = ""
	c.lastFrame = nil
	c.replay = protection.NewReplayDetector()
	c.ackMgr = NewAckManager()

	c.state = StateDisconnected
}
