package main

import (
	"encoding/json" // Dodano dla json.Marshal
	"time"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// handleHELLO processes the initial HELLO frame and starts the AUTH timeout timer.
// It responds with MEOW_OK(status="Ready for auth").
func handleHELLO(stream *quic.Stream, conn *quic.Conn) *protection.AuthTimer {
	// TASK 10: Użycie MeowOkFrame zamiast UniversalFrame
	ok := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "MEOW_OK",
			MsgID: time.Now().UnixMilli(),
		},
		Status: "Ready for auth",
	}
	b, _ := json.Marshal(ok)
	stream.Write(b)

	// Start 20-second AUTH timeout
	return protection.StartAuthTimer(func() {
		// TASK 10: Użycie ErrorFrame zamiast UniversalFrame
		errFrame := protocol.ErrorFrame{
			BaseFrame: protocol.BaseFrame{
				Type:  "ERROR",
				MsgID: time.Now().UnixMilli(),
			},
			Code: "ERR_03",
			Desc: "Authorization timeout reached",
		}
		eb, _ := json.Marshal(errFrame)
		stream.Write(eb)
		conn.CloseWithError(0x03, "ERR_03: Auth Timeout")
	})
}
