package main

import (
	"encoding/json"
	"fmt"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// handleGetStatus processes GET_STATUS:
//   - parses the frame
//   - checks if target user is online
//   - sends STATUS_RES with "online" or "offline"
func (c *clientContext) handleGetStatus(raw []byte) {
	frame, err := protocol.ParseGetStatusFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	online := globalSessions.IsOnline(frame.Target)
	status := "offline"
	if online {
		status = "online"
	}

	res := protocol.StatusResFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "STATUS_RES",
			MsgID: frame.MsgID,
		},
		Target: frame.Target,
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
