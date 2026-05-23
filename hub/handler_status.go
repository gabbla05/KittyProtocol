package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handleGetStatus(raw []byte) {
	if c.state != stateAuthenticated {
		sendError(c.stream, "ERR_02", "GET_STATUS not allowed before AUTH")
		return
	}

	frame, err := protocol.ParseGetStatusFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	target := strings.ToLower(strings.TrimSpace(frame.Target))

	online := globalSessions.IsOnline(target)
	status := "offline"
	if online {
		status = "online"
	}

	res := protocol.StatusResFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeStatusRes,
			MsgID: frame.MsgID,
		},
		Target: target,
		Status: status,
	}

	b, err := json.Marshal(res)
	if err != nil {
		sendError(c.stream, "ERR_02", "Failed to marshal STATUS_RES")
		return
	}

	if _, err := c.stream.Write(b); err != nil {
		fmt.Println("[Hub: Status] Failed to send STATUS_RES:", err)
	}
}
