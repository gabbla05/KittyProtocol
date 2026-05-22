package app

import (
	"encoding/json"
)

// HandleIncomingPayload is called by KittyClient (via callback)
// whenever a decrypted DATA payload arrives.
func (a *App) HandleIncomingPayload(sender string, payload []byte) {
	// Try to decode as ChatFrame
	var cf ChatFrame
	if err := json.Unmarshal(payload, &cf); err == nil && cf.Type != "" {
		a.HandleIncomingChatFrame(cf)
		return
	}

	// Fallback: plain text message
	a.ui.Printf("\n[%s]: %s\n> ", sender, string(payload))
}
