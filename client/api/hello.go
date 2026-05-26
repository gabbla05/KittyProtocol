package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendHello sends the initial HELLO frame on the control stream.
//
// BEHAVIOR:
//   - Requires an established QUIC connection (ensureConnected).
//   - Does not block on any response; the result is delivered via HelloResult().
//   - Intended to be called immediately after Connect().
func (c *KittyClient) SendHello() error {
	stream, err := c.ensureConnected()
	if err != nil {
		return fmt.Errorf("cannot send HELLO: %w", err)
	}

	frame := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeHello,
			MsgID: time.Now().UnixMilli(),
		},
		Version: protocol.CurrentProtocolVersion,
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return fmt.Errorf("failed to marshal HELLO: %w", err)
	}

	if _, err := stream.Write(b); err != nil {
		return fmt.Errorf("failed to send HELLO: %w", err)
	}

	return nil
}
