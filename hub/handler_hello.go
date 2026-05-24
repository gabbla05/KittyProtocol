// handler_hello.go
// Handles the HELLO frame — the first step of the KittyProtocol handshake.

package hub

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handleHello(raw []byte) {
	if c.state != stateInit {
		sendError(c.stream, protocol.ErrProtocolViolation, "HELLO not allowed in current state")
		return
	}

	hello, err := protocol.ParseHelloFrame(raw)
	if err != nil {
		sendError(c.stream, protocol.ErrFormatError, err.Error())
		return
	}

	logInfo("[HELLO] Client version: %s", hello.Version)

	c.state = stateHelloReceived

	// Respond with MEOW_OK
	ok := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeMeowOK,
			MsgID: hello.MsgID,
		},
		Status: "Ready for auth",
	}

	b, err := json.Marshal(ok)
	if err != nil {
		logError("[HELLO] Failed to marshal MEOW_OK: %v", err)
		return
	}

	if _, err := c.stream.Write(b); err != nil {
		logError("[HELLO] Failed to send MEOW_OK: %v", err)
	}

	// Start AUTH timeout (20s from protection.DefaultAuthTimeout)
	c.authTimer = protection.StartAuthTimer(func() {
		sendError(c.stream, protocol.ErrAuthorizationTimeout, "Authorization timeout reached")
		_ = c.conn.CloseWithError(0x03, "ERR_03: Auth Timeout")
	})
}
