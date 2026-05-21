// hub/errors.go
// Centralized helpers for sending protocol-level ERROR frames from the Hub.

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// sendError sends a standardized ERROR frame to the client.
// Serialization or write failures are logged but do not panic.
func sendError(stream *quic.Stream, code, desc string) {
	errFrame := protocol.ErrorFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "ERROR",
			MsgID: time.Now().UnixMilli(),
		},
		Code: code,
		Desc: desc,
	}

	b, err := json.Marshal(errFrame)
	if err != nil {
		fmt.Println("[Hub: Errors] Failed to marshal ERROR frame:", err)
		return
	}

	if _, err := stream.Write(b); err != nil {
		fmt.Println("[Hub: Errors] Failed to send ERROR frame:", err)
	}
}
