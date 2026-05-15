// hub/router.go
package main

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// routeData forwards a DATA frame from the sender to the target session.
// It performs only protocol-level checks (session existence, stream availability),
// without imposing any application-level "chat state" semantics.
func routeData(frame protocol.DataFrame, sender *protection.Session, senderStream *quic.Stream) bool {
	targetSess, ok := globalSessions.Get(frame.Target)
	if !ok {
		// Receiver has no active session.
		sendError(senderStream, "ERR_15", "Receiver offline")
		return false
	}

	if targetSess.Stream == nil {
		sendError(senderStream, "ERR_10", "Receiver stream not available")
		return false
	}

	// Pure routing: Hub does not enforce whether the receiver "accepted" the chat.
	// That logic belongs to the application layer (e.g. Meowssenger).

	forward := protocol.DataFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "DATA",
			MsgID: frame.MsgID,
		},
		Sender:  sender.ID,
		Target:  frame.Target,
		Payload: frame.Payload,
		MAC:     frame.MAC,
	}

	fb, err := json.Marshal(forward)
	if err != nil {
		sendError(senderStream, "ERR_02", "Failed to marshal forwarded DATA")
		return false
	}

	if _, err := targetSess.Stream.Write(fb); err != nil {
		sendError(senderStream, "ERR_10", "Failed to deliver to receiver")
		return false
	}

	return true
}
