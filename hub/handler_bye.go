package main

import (
	"fmt"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handleBye(raw []byte) {
	if c.state != stateAuthenticated {
		sendError(c.stream, "ERR_02", "BYE not allowed before AUTH")
		return
	}

	_, err := protocol.ParseByeFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	fmt.Println("[Handler: Bye] Cleaning up session for:", c.username)

	globalSessions.Remove(c.username)

	if c.session != nil && c.session.CloseFunc != nil {
		c.session.CloseFunc()
	}

	c.session = nil
	c.state = stateInit
}
