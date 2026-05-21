// hub/auth_flow.go
// Implements the initial HELLO → AUTH flow and authorization timeout handling.

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
		if _, err := stream.Write(b); err != nil {
			fmt.Println("[Hub: AuthFlow] Failed to send MEOW_OK:", err)
		}
	} else {
		fmt.Println("[Hub: AuthFlow] Failed to marshal MEOW_OK:", err)
	}

	// Start 20-second AUTH timeout.
	return protection.StartAuthTimer(func() {
		// On timeout, send ERR_03 and close the connection.
		sendError(stream, "ERR_03", "Authorization timeout reached")
		_ = conn.CloseWithError(0x03, "ERR_03: Auth Timeout")
	})
}
