// handler_data.go
// Handles DATA frames — encrypted message delivery between authenticated users.

package hub

import (
	"encoding/json"
	"strings"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func canonicalTarget(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}

func (c *clientContext) handleData(raw []byte) {
	if c.state != stateAuthenticated {
		sendError(c.stream, protocol.ErrProtocolViolation, "DATA not allowed before AUTH")
		return
	}

	frame, err := protocol.ParseDataFrame(raw)
	if err != nil {
		sendError(c.stream, protocol.ErrFormatError, err.Error())
		return
	}

	// Normalize target username
	frame.Target = canonicalTarget(frame.Target)

	// Sanity check
	if c.session == nil {
		sendError(c.stream, protocol.ErrProtocolViolation, "DATA before AUTH")
		return
	}

	// Rate limit
	if !c.session.Limiter.Allow() {
		sendError(c.stream, protocol.ErrRateLimitExceeded, "Rate limit exceeded")
		return
	}

	// Replay protection
	if c.session.Replay.MarkAndCheck(frame.MsgID) {
		sendError(c.stream, protocol.ErrReplayDetected, "Replay detected")
		return
	}

	// Update activity
	c.touch()

	// Route message
	if !routeData(*frame, c.session, c.stream) {
		return
	}

	// Send ACK
	ack := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeMeowOK,
			MsgID: frame.MsgID,
		},
		Status: "Delivered",
	}

	b, err := json.Marshal(ack)
	if err != nil {
		logError("[DATA] Failed to marshal ACK: %v", err)
		return
	}

	if _, err := c.stream.Write(b); err != nil {
		logError("[DATA] Failed to send ACK: %v", err)
	}
}
