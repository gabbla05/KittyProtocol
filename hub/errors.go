// errors.go
// Centralized helpers for sending protocol-level ERROR frames from the Hub.
// All error codes MUST come from protocol/error_codes.go.

package hub

import (
	"encoding/json"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// sendError sends a standardized ERROR frame to the client.
// Serialization or write failures are logged but do not panic.
// This function MUST be used by all handlers to ensure consistent error reporting.
func sendError(stream *quic.Stream, code, desc string) {
	errFrame := protocol.ErrorFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeError,
			MsgID: time.Now().UnixMilli(),
		},
		Code: code,
		Desc: desc,
	}

	b, err := json.Marshal(errFrame)
	if err != nil {
		logError("[Errors] Failed to marshal ERROR frame: %v", err)
		return
	}

	if _, err := stream.Write(b); err != nil {
		logError("[Errors] Failed to send ERROR frame: %v", err)
	}
}
