package api

import (
	"errors"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendBye sends a BYE frame to the Hub.
//
// BYE is best-effort: if the client is already disconnected or has no stream,
// the method returns nil without treating it as an error.
func (c *KittyClient) SendBye() error {
	stream, err := c.ensureConnected()
	if err != nil {
		if errors.Is(err, ErrNoStream) || errors.Is(err, ErrNotConnected) {
			return nil
		}
		return err
	}

	frame := protocol.BaseFrame{
		Type:  protocol.FrameTypeBye,
		MsgID: time.Now().UnixMilli(),
	}

	return c.sendFrame(stream, frame)
}
