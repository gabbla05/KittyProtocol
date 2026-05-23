package main

import (
	"encoding/json"
	"fmt"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *clientContext) handlePing(raw []byte) {
	if c.state != stateAuthenticated {
		sendError(c.stream, "ERR_02", "PING not allowed before AUTH")
		return
	}

	_, err := protocol.ParsePingFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	c.touch()
}

func ParsePingFrame(data []byte) (*protocol.PingFrame, error) {
	var f protocol.PingFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", protocol.ErrCodeInvalidFrame)
	}
	if f.Type != protocol.FrameTypePing {
		return nil, fmt.Errorf("%s: Invalid type for PING frame", protocol.ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in PING frame", protocol.ErrCodeInvalidFrame)
	}
	return &f, nil
}
