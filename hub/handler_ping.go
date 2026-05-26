// handler_ping.go
// Handles PING frames — keeps the session alive.

package hub

import (
	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handlePing(raw []byte) {
	if c.state != stateAuthenticated {
		sendError(c.stream, protocol.ErrProtocolViolation, "PING not allowed before AUTH")
		return
	}

	if _, err := protocol.ParsePingFrame(raw); err != nil {
		sendError(c.stream, protocol.ErrFormatError, err.Error())
		return
	}

	c.touch()
}
