package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// -----------------------------------------------------------------------------
// AUTH
// -----------------------------------------------------------------------------

func (c *KittyClient) SendAuth(user, pass string) error {
	stream, err := c.ensureConnected()
	if err != nil {
		return err
	}

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeAuth,
			MsgID: time.Now().UnixMilli(),
		},
		User: user,
		Pass: pass,
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	// Save username
	c.mu.Lock()
	c.user = user
	c.mu.Unlock()

	_, err = stream.Write(b)
	return err
}

func (c *KittyClient) WaitForAuthOK() error {
	ok, code := c.waitOkOrError()
	if !ok {
		return errors.New(code)
	}
	return nil
}

// -----------------------------------------------------------------------------
// REGISTER
// -----------------------------------------------------------------------------

func (c *KittyClient) SendRegister(user, pass string) error {
	stream, err := c.ensureConnected()
	if err != nil {
		return err
	}

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeRegister,
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

func (c *KittyClient) WaitForRegisterOK() error {
	ok, code := c.waitOkOrError()
	if !ok {
		return errors.New(code)
	}
	return nil
}

// -----------------------------------------------------------------------------
// Shared helper
// -----------------------------------------------------------------------------

func (c *KittyClient) waitOkOrError() (bool, string) {
	c.mu.Lock()
	stream := c.stream
	c.mu.Unlock()

	if stream == nil {
		return false, "NO_STREAM"
	}

	buf := make([]byte, defaultRecvBufferSize)
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
			if errFrame.Desc != "" {
				return false, fmt.Sprintf("%s: %s", errFrame.Code, errFrame.Desc)
			}
			return false, errFrame.Code
		}
		return false, "PARSE_ERROR"

	case protocol.FrameTypeMeowOK:
		return true, ""
	}

	return false, "UNKNOWN_FRAME"
}
