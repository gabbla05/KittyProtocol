package api

import (
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendGetStatus sends a GET_STATUS frame for a given user.
func (c *KittyClient) SendGetStatus(target string) error {
	stream, err := c.ensureConnected()
	if err != nil {
		return err
	}

	if target == "" {
		return ErrTargetNotSet
	}

	if len(target) > maxUsernameLength {
		return ErrTargetNameTooLong
	}

	msgID := time.Now().UnixMilli()

	frame := protocol.GetStatusFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeGetStatus,
			MsgID: msgID,
		},
		Target: canonicalTarget(target),
	}

	return c.sendFrame(stream, frame)
}
