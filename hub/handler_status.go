// handler_status.go
// Handles GET_STATUS — checks whether a target user is online.

package hub

import (
	"encoding/json"
	"strings"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handleGetStatus(raw []byte) {
	if c.state != stateAuthenticated {
		sendError(c.stream, protocol.ErrProtocolViolation, "GET_STATUS not allowed before AUTH")
		return
	}

	frame, err := protocol.ParseGetStatusFrame(raw)
	if err != nil {
		sendError(c.stream, protocol.ErrFormatError, err.Error())
		return
	}

	target := strings.ToLower(strings.TrimSpace(frame.Target))

	status := "offline"
	if globalSessions.IsOnline(target) {
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
		sendError(c.stream, protocol.ErrFormatError, "Failed to marshal STATUS_RES")
		return
	}

	if _, err := c.stream.Write(b); err != nil {
		logError("[STATUS] Failed to send STATUS_RES: %v", err)
	}
}
