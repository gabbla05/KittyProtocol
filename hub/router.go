package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func routeData(frame protocol.DataFrame, sender *protection.Session, senderStream *quic.Stream) bool {
	targetSess, ok := globalSessions.Get(frame.Target)
	if !ok {
		sendError(senderStream, "ERR_15", "Receiver offline")
		return false
	}

	if targetSess.Stream == nil {
		sendError(senderStream, "ERR_10", "Receiver stream not available")
		return false
	}

	sender.LastActive = time.Now()
	targetSess.LastActive = time.Now()

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
		fmt.Println("[Hub: Router] Failed to marshal forwarded DATA:", err)
		sendError(senderStream, "ERR_02", "Failed to marshal forwarded DATA")
		return false
	}

	if _, err := targetSess.Stream.Write(fb); err != nil {
		fmt.Println("[Hub: Router] Failed to deliver DATA:", err)
		sendError(senderStream, "ERR_10", "Failed to deliver to receiver")
		return false
	}

	return true
}
