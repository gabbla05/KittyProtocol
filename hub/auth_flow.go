package main

import (
	"time"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// handleHELLO processes the initial HELLO frame and starts the AUTH timeout timer.
// It responds with MEOW_OK(status="Ready for auth").
func handleHELLO(stream *quic.Stream, conn *quic.Conn) *protection.AuthTimer {
    ok := protocol.UniversalFrame{
        Type:   "MEOW_OK",
        MsgID:  time.Now().UnixMilli(),
        Status: "Ready for auth",
    }
    stream.Write(ok.ToJSON())

    // Start 20-second AUTH timeout
    return protection.StartAuthTimer(func() {
        errFrame := protocol.UniversalFrame{
            Type:  "ERROR",
            MsgID: time.Now().UnixMilli(),
            Code:  "ERR_03",
            Desc:  "Authorization timeout reached",
        }
        stream.Write(errFrame.ToJSON())
        conn.CloseWithError(0x03, "ERR_03: Auth Timeout")
    })
}
