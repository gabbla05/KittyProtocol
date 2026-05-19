package main

import (
	"encoding/json"
	"fmt"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// handleData processes the DATA frame.
// Steps:
// 1. Parse and validate frame
// 2. Ensure session exists
// 3. Validate target
// 4. Apply rate limiting
// 5. Apply replay protection
// 6. Update activity timestamp
// 7. Route to target session
// 8. Send MEOW_OK ACK to sender
func (c *clientContext) handleData(raw []byte) {
	frame, err := protocol.ParseDataFrame(raw)
	if err != nil {
		sendError(c.stream, "ERR_02", err.Error())
		return
	}

	if frame.Target == "" {
		sendError(c.stream, "ERR_02", "Missing target")
		return
	}

	if c.session == nil {
		sendError(c.stream, "ERR_01", "DATA before AUTH")
		return
	}

	// Rate limiting
	if !c.session.Limiter.Allow() {
		sendError(c.stream, "ERR_07", "Rate limit exceeded")
		return
	}

	// Replay protection
	if c.session.Replay != nil && c.session.Replay.MarkAndCheck(frame.MsgID) {
		sendError(c.stream, "ERR_06", "Replay detected")
		return
	}

	// Update activity
	c.touch()

	// Route to target
	if !routeData(*frame, c.session, c.stream) {
		return
	}

	// ACK for sender
	ack := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "MEOW_OK",
			MsgID: frame.MsgID,
		},
		Status: "Delivered (mock)",
	}
	b, err := json.Marshal(ack)
	if err != nil {
		fmt.Println("[Hub: Data] Failed to marshal MEOW_OK ACK:", err)
		return
	}
	if _, err := c.stream.Write(b); err != nil {
		fmt.Println("[Hub: Data] Failed to send MEOW_OK ACK:", err)
	}
}
