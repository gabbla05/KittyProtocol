package main

import (
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// startPingLoop periodically sends PING frames to keep the session alive.
func startPingLoop(stream *quic.Stream) {
    go func() {
        for {
            time.Sleep(30 * time.Second)
            ping := protocol.UniversalFrame{
                Type:  "PING",
                MsgID: time.Now().UnixMilli(),
            }
            stream.Write(ping.ToJSON())
        }
    }()
}
