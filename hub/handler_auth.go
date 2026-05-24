// handler_auth.go
// Handles the AUTH frame — the second step of the KittyProtocol handshake.
// After successful authentication, a session is created and the client enters
// the stateAuthenticated state.

package hub

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// handleAuth processes an AUTH frame after a successful HELLO.
// Expected state: stateHelloReceived → stateAuthenticated.
func (c *clientContext) handleAuth(raw []byte) {
	if c.state != stateHelloReceived {
		sendError(c.stream, protocol.ErrProtocolViolation, "AUTH not allowed before HELLO")
		return
	}

	frame, err := protocol.ParseAuthFrame(raw)
	if err != nil {
		sendError(c.stream, protocol.ErrFormatError, err.Error())
		return
	}

	// Stop AUTH timeout
	if c.authTimer != nil {
		c.authTimer.Stop()
		c.authTimer = nil
	}

	// Validate credentials
	if !globalAuth.CheckCredentials(frame.User, frame.Pass) {
		sendError(c.stream, protocol.ErrAuthenticationFailed, "Authentication failed")
		return
	}

	// Prevent duplicate logins
	if globalSessions.IsOnline(frame.User) {
		sendError(c.stream, protocol.ErrSessionError, "User already logged in")
		return
	}

	// Create session
	c.session = protection.NewSession(frame.User, c.conn, c.stream)
	globalSessions.Add(frame.User, c.session)
	c.username = frame.User
	c.state = stateAuthenticated

	logInfo("[AUTH] User '%s' authenticated successfully", frame.User)

	// Send MEOW_OK
	ok := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeMeowOK,
			MsgID: frame.MsgID,
		},
		Status: "Logged in",
	}

	b, err := json.Marshal(ok)
	if err != nil {
		logError("[AUTH] Failed to marshal MEOW_OK: %v", err)
		return
	}

	if _, err := c.stream.Write(b); err != nil {
		logError("[AUTH] Failed to send MEOW_OK: %v", err)
	}
}
