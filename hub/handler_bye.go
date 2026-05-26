// handler_bye.go
// Handles the BYE frame — clean session termination requested by the client.

package hub

import (
	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handleBye(raw []byte) {
	if c.state != stateAuthenticated {
		sendError(c.stream, protocol.ErrProtocolViolation, "BYE not allowed before AUTH")
		return
	}

	if _, err := protocol.ParseByeFrame(raw); err != nil {
		sendError(c.stream, protocol.ErrFormatError, err.Error())
		return
	}

	logInfo("[BYE] Cleaning up session for user: %s", c.username)

	globalSessions.Remove(c.username)

	if c.session != nil && c.session.CloseFunc != nil {
		c.session.CloseFunc()
	}

	c.session = nil
	c.state = stateInit
}
