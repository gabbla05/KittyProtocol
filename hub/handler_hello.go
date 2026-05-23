package main

import (
	"encoding/json"
	"fmt"

	"github.com/gabbla05/KittyProtocol/internal/protection"
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

	// Set state
	c.state = stateHelloReceived

	// Send MEOW_OK
	ok := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeMeowOK,
			MsgID: hello.MsgID,
		},
		Status: "Ready for auth",
	}

	if b, err := json.Marshal(ok); err == nil {
		_, _ = c.stream.Write(b)
	} else {
		fmt.Println("[Hub: HELLO] Failed to marshal MEOW_OK:", err)
	}

	// Start AUTH timeout
	c.authTimer = protection.StartAuthTimer(func() {
		sendError(c.stream, "ERR_03", "Authorization timeout reached")
		_ = c.conn.CloseWithError(0x03, "ERR_03: Auth Timeout")
	})
}
