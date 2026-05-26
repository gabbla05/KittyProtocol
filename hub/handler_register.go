// handler_register.go
// Handles the REGISTER frame — creates a new user account.
// REGISTER does NOT authenticate the user; AUTH must follow.

package hub

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handleRegister(raw []byte) {
	if c.state != stateHelloReceived {
		sendError(c.stream, protocol.ErrProtocolViolation, "REGISTER not allowed before HELLO")
		return
	}

	frame, err := protocol.ParseRegisterFrame(raw)
	if err != nil {
		sendError(c.stream, protocol.ErrFormatError, err.Error())
		return
	}

	if err := globalAuth.Register(frame.User, frame.Pass); err != nil {
		sendError(c.stream, protocol.ErrSessionError, err.Error())
		return
	}

	logInfo("[REGISTER] User '%s' registered successfully", frame.User)

	ok := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeMeowOK,
			MsgID: frame.MsgID,
		},
		Status: "Registered",
	}

	b, err := json.Marshal(ok)
	if err != nil {
		logError("[REGISTER] Failed to marshal MEOW_OK: %v", err)
		return
	}

	if _, err := c.stream.Write(b); err != nil {
		logError("[REGISTER] Failed to send MEOW_OK: %v", err)
	}
}
