package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *KittyClient) SendHello() error {
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

	_, err = c.stream.Write(b)
	if err != nil {
		return fmt.Errorf("failed to send HELLO: %w", err)
	}

	return nil
}

// waitHelloOK waits for MEOW_OK or ERROR after HELLO.
// Returns (success, errorCode).
//
// BLOCKING BEHAVIOR:
//   - This call blocks until the Hub responds or the stream errors.
//   - It does not enforce a timeout; QUIC idle timeout applies.
//
// PROTOCOL:
//   - Expected responses: MEOW_OK or ERROR.
//   - Any other frame type is treated as UNKNOWN_FRAME.
func (c *KittyClient) waitHelloOK() (bool, string) {
	c.mu.Lock()
	stream := c.stream
	c.mu.Unlock()

	if stream == nil {
		return false, "NO_STREAM"
	}

	buf := make([]byte, 4096)
	n, err := stream.Read(buf)
	if err != nil {
		return false, "READ_ERROR"
	}

	typeName, _, err := protocol.GetFrameType(buf[:n])
	if err != nil {
		return false, "PARSE_ERROR"
	}

	switch typeName {
	case protocol.FrameTypeError:
		var errFrame protocol.ErrorFrame
		if json.Unmarshal(buf[:n], &errFrame) == nil {
			return false, errFrame.Code
		}
		return false, "PARSE_ERROR"

	case protocol.FrameTypeMeowOK:
		var okFrame protocol.MeowOkFrame
		if json.Unmarshal(buf[:n], &okFrame) != nil {
			return false, "PARSE_ERROR"
		}
		return true, ""
	}

	return false, "UNKNOWN_FRAME"
}
