package main

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// handleAuth processes the AUTH frame.
// Steps:
// 1. Parse and validate the frame
// 2. Stop AUTH timer
// 3. Verify credentials using globalAuth
// 4. Create a session and register it in SessionManager
// 5. Send MEOW_OK("Logged in")
func (c *clientContext) handleAuth(raw []byte) {
	frame, err := protocol.ParseAuthFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	// Stop AUTH timer
	if c.authTimer != nil {
		c.authTimer.Stop()
		c.authTimer = nil
	}

	// Verify credentials
	if !globalAuth.CheckCredentials(frame.User, frame.Pass) {
		sendError(c.stream, "ERR_04", "Authentication failed")
		return
	}

	// Create session
	c.session = protection.NewSession(frame.User, c.conn, c.stream)
	globalSessions.Add(frame.User, c.session)
	c.username = frame.User

	// Send MEOW_OK
	ok := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "MEOW_OK",
			MsgID: frame.MsgID,
		},
		Status: "Logged in",
	}
	if b, err := json.Marshal(ok); err == nil {
		_, _ = c.stream.Write(b)
	}
}
