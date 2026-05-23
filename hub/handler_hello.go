package main

import (
	"fmt"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handleHello(raw []byte) {
	if c.state != stateInit {
		sendError(c.stream, "ERR_02", "HELLO not allowed in current state")
		return
	}

	hello, err := protocol.ParseHelloFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	fmt.Println("[Hub] HELLO from client, version:", hello.Version)

	c.authTimer = handleHELLO(c.stream, c.conn)
	c.state = stateHelloReceived
}
