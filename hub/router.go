// router.go
// Implements DATA frame forwarding between authenticated sessions.
// This file contains no protocol parsing — only delivery logic.

package hub

import (
	"encoding/json"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// routeData forwards a DATA frame from sender → receiver.
// Returns true on success, false if delivery failed.
func routeData(frame protocol.DataFrame, sender *protection.Session, senderStream *quic.Stream) bool {
	targetSess, ok := globalSessions.Get(frame.Target)
	// router.go — poprawiony fragment
	if !ok {
		// ERR_15 — Unknown Target
		sendError(senderStream, protocol.ErrUnknownTarget, "Unknown target user")
		return false
	}

	if targetSess.Stream == nil {
		// ERR_05 — Session Error (receiver session corrupted)
		sendError(senderStream, protocol.ErrSessionError, "Receiver stream not available")
		return false
	}

	// Update activity timestamps
	now := time.Now()
	sender.LastActive = now
	targetSess.LastActive = now

	// Build forwarded DATA frame
	forward := protocol.DataFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeData,
			MsgID: frame.MsgID,
		},
		Sender:  sender.ID,
		Target:  frame.Target,
		Payload: frame.Payload,
		MAC:     frame.MAC,
	}

	fb, err := json.Marshal(forward)
	if err != nil {
		logError("[Router] Failed to marshal forwarded DATA: %v", err)
		sendError(senderStream, protocol.ErrFormatError, "Failed to marshal forwarded DATA")
		return false
	}

	if _, err := targetSess.Stream.Write(fb); err != nil {
		logError("[Router] Failed to deliver DATA: %v", err)
		sendError(senderStream, protocol.ErrSessionError, "Failed to deliver to receiver")
		return false
	}

	return true
}
