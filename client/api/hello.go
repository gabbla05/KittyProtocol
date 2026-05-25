package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func (c *KittyClient) SendHello() error {
	c.mu.Lock()
	stream := c.stream
	c.mu.Unlock()

	if stream == nil {
		return fmt.Errorf("stream is nil")
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
