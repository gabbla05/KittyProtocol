package main

import (
	"context"
	"fmt"
	"io"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// handleClient is the entry point for handling a single QUIC connection.
//
// Responsibilities:
//   - accept the bidirectional stream
//   - read incoming frames
//   - perform initial validation (type + msg_id)
//   - dispatch to dedicated handlers
func handleClient(conn *quic.Conn) {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		fmt.Println("[Hub: HandlerDispatcher] Stream accept error:", err)
		return
	}
	defer stream.Close()

	ctx := &clientContext{
		conn:   conn,
		stream: stream,
	}
	defer ctx.cleanup()

	buf := make([]byte, 4096)

	for {
		n, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("[Hub: HandlerDispatcher] Stream read error:", err)
			}
			return
		}

		raw := buf[:n]
		fmt.Println("[Hub: HandlerDispatcher] STREAM ID:", stream.StreamID())
		fmt.Println("[Hub: HandlerDispatcher] RAW:", string(raw))

		// Initial validation: extract type and msg_id.
		typeName, _, perr := protocol.GetFrameType(raw)
		if perr != nil {
			sendError(stream, "ERR_02", perr.Error())
			continue
		}

		fmt.Println("[Hub: HandlerDispatcher] Received frame type:", typeName)

		switch typeName {
		case "HELLO":
			ctx.handleHello()

		case "AUTH":
			ctx.handleAuth(raw)

		case "PING":
			ctx.handlePing()

		case "DATA":
			ctx.handleData(raw)

		case "GET_STATUS":
			ctx.handleGetStatus(raw)

		case "STATUS_RES":
			// Hub does not expect STATUS_RES from clients.
			fmt.Println("[Hub: HandlerDispatcher] Unexpected STATUS_RES from client – ignoring")

		case "BYE":
			ctx.handleBye()
			return

		default:
			sendError(stream, "ERR_02", "Unknown frame type: "+typeName)
		}
	}
}
