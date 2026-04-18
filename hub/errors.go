package main

import (
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// sendError sends a standardized ERROR frame to the client.
func sendError(stream *quic.Stream, code, desc string) {
    errFrame := protocol.UniversalFrame{
        Type:  "ERROR",
        MsgID: time.Now().UnixMilli(),
        Code:  code,
        Desc:  desc,
    }
    stream.Write(errFrame.ToJSON())
}
