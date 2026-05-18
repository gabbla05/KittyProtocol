// hub/auth_flow.go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// handleHELLO processes the initial HELLO frame and starts the AUTH timeout timer.
// It responds with MEOW_OK(status="Ready for auth").
func handleHELLO(stream *quic.Stream, conn *quic.Conn) *protection.AuthTimer {
	ok := protocol.MeowOkFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "MEOW_OK",
			MsgID: time.Now().UnixMilli(),
		},
		Status: "Ready for auth",
	}
	if b, err := json.Marshal(ok); err == nil {
		_, _ = stream.Write(b)
	} else {
		fmt.Println("[Hub] Failed to send MEOW_OK:", err)
	}

	// Start 20-second AUTH timeout.
	return protection.StartAuthTimer(func() {
		errFrame := protocol.ErrorFrame{
			BaseFrame: protocol.BaseFrame{
				Type:  "ERROR",
				MsgID: time.Now().UnixMilli(),
			},
			Code: "ERR_03",
			Desc: "Authorization timeout reached",
		}
		if eb, err := json.Marshal(errFrame); err == nil {
			_, _ = stream.Write(eb)
		}
		conn.CloseWithError(0x03, "ERR_03: Auth Timeout")
	})
}
