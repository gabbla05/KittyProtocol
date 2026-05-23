package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func canonicalTarget(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}

func (c *clientContext) handleData(raw []byte) {
	if c.state != stateAuthenticated {
		sendError(c.stream, "ERR_02", "DATA not allowed before AUTH")
		return
	}

	frame, err := protocol.ParseDataFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	frame.Target = canonicalTarget(frame.Target)

	if c.session == nil {
		sendError(c.stream, "ERR_01", "DATA before AUTH")
		return
	}

	if !c.session.Limiter.Allow() {
		sendError(c.stream, "ERR_07", "Rate limit exceeded")
		return
	}

	if c.session.Replay.MarkAndCheck(frame.MsgID) {
		sendError(c.stream, "ERR_06", "Replay detected")
		return
	}

	c.touch()

	if !routeData(*frame, c.session, c.stream) {
		return
	}

	ack := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "MEOW_OK",
			MsgID: frame.MsgID,
		},
		Status: "Delivered",
	}

	b, err := json.Marshal(ack)
	if err == nil {
		c.stream.Write(b)
	} else {
		fmt.Println("[Hub: Data] Failed to marshal ACK:", err)
	}
}
