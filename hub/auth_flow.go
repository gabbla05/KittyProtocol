// hub/auth_flow.go
// Implements the initial HELLO → AUTH flow and authorization timeout handling.

package main

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/gabbla05/KittyProtocol/internal/protection"
// 	"github.com/gabbla05/KittyProtocol/protocol"
// )

// // handleHello processes the initial HELLO frame and starts the AUTH timeout timer.
// // It responds with MEOW_OK(status="Ready for auth").
// func (c *clientContext) handleHello(raw []byte) {
// 	// Parse HELLO
// 	frame, err := protocol.ParseHelloFrame(raw)
// 	if err != nil {
// 		sendError(c.stream, "ERR_02", err.Error())
// 		return
// 	}

// 	// Set state
// 	c.state = stateHelloReceived

// 	// Send MEOW_OK
// 	ok := protocol.MeowOkFrame{
// 		BaseFrame: protocol.BaseFrame{
// 			Type:  protocol.FrameTypeMeowOK,
// 			MsgID: frame.MsgID,
// 		},
// 		Status: "Ready for auth",
// 	}

// 	b, err := json.Marshal(ok)
// 	if err == nil {
// 		_, _ = c.stream.Write(b)
// 	} else {
// 		fmt.Println("[Hub: HELLO] Failed to marshal MEOW_OK:", err)
// 	}

// 	// Start AUTH timeout
// 	c.authTimer = protection.StartAuthTimer(func() {
// 		sendError(c.stream, "ERR_03", "Authorization timeout reached")
// 		_ = c.conn.CloseWithError(0x03, "ERR_03: Auth Timeout")
// 	})
// }
