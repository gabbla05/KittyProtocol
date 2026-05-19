package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendAuth sends the AUTH frame with username and password.
func (c *KittyClient) SendAuth(user, pass string) error {
	c.mu.Lock()
	stream := c.stream
	c.mu.Unlock()

	if stream == nil {
		return errors.New("stream is nil")
	}

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "AUTH",
			MsgID: time.Now().UnixMilli(),
		},
		User: user,
		Pass: pass,
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	_, err = stream.Write(b)
	return err
}

// waitAuthOK waits for MEOW_OK or ERROR after AUTH.
// Returns (success, errorCode).
func (c *KittyClient) waitAuthOK() (bool, string) {
	c.mu.Lock()
	stream := c.stream
	c.mu.Unlock()

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
	case "ERROR":
		var errFrame protocol.ErrorFrame
		if json.Unmarshal(buf[:n], &errFrame) == nil {
			return false, errFrame.Code
		}
		return false, "PARSE_ERROR"

	case "MEOW_OK":
		var okFrame protocol.MeowOkFrame
		if json.Unmarshal(buf[:n], &okFrame) != nil {
			return false, "PARSE_ERROR"
		}
		return true, ""
	}

	return false, "UNKNOWN_FRAME"
}
