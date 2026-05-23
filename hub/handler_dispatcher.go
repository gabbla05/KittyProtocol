package main

import (
	"context"
	"io"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func handleClient(conn *quic.Conn) {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		logError("Stream accept error: %v", err)
		return
	}
	defer stream.Close()

	ctx := &clientContext{
		conn:   conn,
		stream: stream,
		state:  stateInit,
	}
	defer ctx.cleanup()

	buf := make([]byte, 8192)

	for {
		n, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF {
				logError("Stream read error: %v", err)
			}
			return
		}

		raw := buf[:n]

		typeName, msgID, perr := protocol.GetFrameType(raw)
		if perr != nil || msgID <= 0 {
			sendError(stream, "ERR_02", "Invalid frame header")
			continue
		}

		switch typeName {
		case protocol.FrameTypeHello:
			ctx.handleHello(raw)

		case protocol.FrameTypeAuth:
			ctx.handleAuth(raw)

		case protocol.FrameTypeRegister:
			ctx.handleRegister(raw)

		case protocol.FrameTypePing:
			ctx.handlePing(raw)

		case protocol.FrameTypeData:
			ctx.handleData(raw)

		case protocol.FrameTypeGetStatus:
			ctx.handleGetStatus(raw)

		case protocol.FrameTypeBye:
			ctx.handleBye(raw)
			return

		default:
			sendError(stream, "ERR_02", "Unknown frame type")
		}
	}
}
