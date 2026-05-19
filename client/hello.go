package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendHello sends the HELLO frame to the Hub.
func (c *KittyClient) SendHello() error {
	c.mu.Lock()
	stream := c.stream
	c.mu.Unlock()

	if stream == nil {
		return errors.New("stream is nil")
	}

	frame := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "HELLO",
			MsgID: time.Now().UnixMilli(),
		},
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	_, err = stream.Write(b)
	return err
}

// waitHelloOK waits for MEOW_OK or ERROR after HELLO.
// Returns (success, errorCode).
func (c *KittyClient) waitHelloOK() (bool, string) {
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
