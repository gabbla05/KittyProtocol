package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendAuth sends an AUTH frame with username and password to the Hub.
//
// SECURITY:
//   - Credentials are transmitted inside a TLS 1.3 encrypted QUIC stream.
//   - The Hub validates credentials and responds with MEOW_OK or ERROR.
//
// PROTOCOL ORDER:
//  1. SendHello()
//  2. WaitForHelloOK()
//  3. SendAuth()
//  4. WaitForAuthOK()
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
//
// BLOCKING BEHAVIOR:
//   - This call blocks until the Hub responds or the stream errors.
//   - QUIC idle timeout ensures this does not block indefinitely.
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
