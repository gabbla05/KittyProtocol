package main

import (
	"encoding/json"
	"fmt"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handleAuth(raw []byte) {
	if c.state != stateHelloReceived {
		sendError(c.stream, "ERR_02", "AUTH not allowed before HELLO")
		return
	}

	frame, err := protocol.ParseAuthFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	if c.authTimer != nil {
		c.authTimer.Stop()
		c.authTimer = nil
	}

	if !globalAuth.CheckCredentials(frame.User, frame.Pass) {
		sendError(c.stream, "ERR_04", "Authentication failed")
		return
	}

	if globalSessions.IsOnline(frame.User) {
		sendError(c.stream, "ERR_05", "User already logged in")
		return
	}

	c.session = protection.NewSession(frame.User, c.conn, c.stream)
	globalSessions.Add(frame.User, c.session)
	c.username = frame.User
	c.state = stateAuthenticated

	ok := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "MEOW_OK",
			MsgID: frame.MsgID,
		},
		Status: "Logged in",
	}

	b, err := json.Marshal(ok)
	if err == nil {
		c.stream.Write(b)
	} else {
		fmt.Println("[Hub: Auth] Failed to marshal MEOW_OK:", err)
	}
}
