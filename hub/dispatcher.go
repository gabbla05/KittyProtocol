// dispatcher.go
// Central QUIC stream dispatcher. Continuously reads raw frames from the
// client stream, determines their type, and forwards them to the appropriate
// handler. This file contains no business logic — only frame routing and
// connection lifecycle management.

package hub

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func handleClient(conn *quic.Conn) {
	// Accept a bidirectional stream from the client.
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		// This is a real error — client failed to open a stream.
		logError("Stream accept error: %v", err)
		return
	}
	defer stream.Close()

	// Per-client context used by handlers.
	ctx := &clientContext{
		conn:   conn,
		stream: stream,
		state:  stateInit,
	}
	defer ctx.cleanup()

	buf := make([]byte, readBufferSize)
	formatErrCount := 0

	for {
		n, err := stream.Read(buf)
		if err != nil {

			// --- NORMAL STREAM TERMINATION ---
			// QUIC-go uses ApplicationError for remote close.
			// These are NOT real errors — they indicate the client closed the stream.
			var appErr *quic.ApplicationError
			if err == io.EOF ||
				errors.As(err, &appErr) ||
				strings.Contains(err.Error(), "client closed") ||
				strings.Contains(err.Error(), "canceled by remote") {

				logInfo("[Client] Stream closed by remote: %v", err)
				return
			}

			// --- REAL ERROR ---
			logError("Stream read error: %v", err)
			return
		}

		raw := buf[:n]

		// Determine frame type and validate header.
		typeName, msgID, perr := protocol.GetFrameType(raw)
		if perr != nil || msgID <= 0 {
			formatErrCount++
			sendError(stream, protocol.ErrFormatError, "Invalid frame header")

			if formatErrCount >= maxFormatErrors {
				logWarn("Too many malformed frames from client — closing connection")
				return
			}
			continue
		}

		formatErrCount = 0

		// Dispatch frame to the appropriate handler.
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
			sendError(stream, protocol.ErrFormatError, "Unknown frame type")
		}
	}
}
