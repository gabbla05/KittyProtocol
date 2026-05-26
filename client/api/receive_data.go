package api

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// handleDataFrame processes a DATA frame: verifies replay protection,
// decrypts and authenticates the payload, then dispatches either chat
// control events or generic application payloads.
func (c *KittyClient) handleDataFrame(frameBytes []byte) {
	var df protocol.DataFrame
	if json.Unmarshal(frameBytes, &df) != nil {
		log(LogError, "failed to parse DATA frame")
		return
	}

	c.mu.Lock()
	replay := c.replay
	kEnc, kMac, ok := c.getKeysForPeer(df.Sender)
	appHandler := c.appHandler
	c.mu.Unlock()

	if replay != nil && replay.MarkAndCheck(df.MsgID) {
		// Replay detected — silently drop.
		return
	}

	if !ok {
		log(LogWarn, "no shared secret for %s", df.Sender)
		return
	}

	plaintext, err := cryptoee.DecryptAndVerifyWithKeys(
		df.MsgID, df.Target, df.Payload, df.MAC, kEnc, kMac,
	)
	if err != nil {
		log(LogError, "E2EE error: %v", err)
		return
	}

	// Try to interpret as chat control frame first.
	if handled := c.handleChatPayload(df.Sender, []byte(plaintext)); handled {
		return
	}

	// Fallback: generic application payload.
	if appHandler != nil {
		appHandler(df.Sender, []byte(plaintext))
	}
}
