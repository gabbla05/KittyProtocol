// hub/router.go
package main

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func routeData(frame protocol.DataFrame, sender *protection.Session, senderStream *quic.Stream) bool {
	targetSess, ok := globalSessions.Get(frame.Target)
	if !ok {
		// Odbiorca nie ma aktywnej sesji
		sendError(senderStream, "ERR_15", "Receiver offline")
		return false
	}

	if targetSess.Stream == nil {
		sendError(senderStream, "ERR_10", "Receiver stream not available")
		return false
	}

	// NOWE: odbiorca jest zalogowany, ale jeszcze nie „ wszedł w rozmowę”
	// (np. nie wybrał targetu / nie wysłał żadnego DATA).
	if !targetSess.ReadyForChat {
		sendError(senderStream, "ERR_16", "Receiver not ready for chat")
		return false
	}

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
