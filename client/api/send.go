package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendMessage encrypts the plaintext, registers ACK tracking,
// and sends a DATA frame to the Hub.
func (c *KittyClient) SendMessage(text string) error {
	c.mu.Lock()
	stream := c.stream
	target := c.target
	ackMgr := c.ackMgr
	kEnc := c.kEnc
	kMac := c.kMac
	c.mu.Unlock()

	if kEnc == nil || kMac == nil {
		return errors.New("shared secret not set")
	}

	if stream == nil {
		return errors.New("stream is nil")
	}
	if target == "" {
		return errors.New("target not set")
	}

	msgID := time.Now().UnixMilli()

	// Register pending ACK
	if ackMgr != nil {
		ackMgr.AddPending(msgID)
	}

	// E2EE encryption
	payloadB64, macB64, err := cryptoee.EncryptAndMACWithKeys(msgID, target, text, kEnc, kMac)
	if err != nil {
		return err
	}

	frame := protocol.DataFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "DATA",
			MsgID: msgID,
		},
		Target:  target,
		Payload: payloadB64,
		MAC:     macB64,
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	// --- message sending ---
	_, err = stream.Write(b)
	if err != nil {
		return err
	}

	// --- remember last frame (for /replay) ---
	c.mu.Lock()
	c.lastFrame = b
	c.mu.Unlock()

	return nil
}

// SendGetStatus sends a GET_STATUS frame for a given user.
func (c *KittyClient) SendGetStatus(target string) error {
	c.mu.Lock()
	stream := c.stream
	c.mu.Unlock()

	if stream == nil {
		return errors.New("stream is nil")
	}

	msgID := time.Now().UnixMilli()

	frame := protocol.GetStatusFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "GET_STATUS",
			MsgID: msgID,
		},
		Target: target,
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	_, err = stream.Write(b)
	return err
}

// SendBye sends a BYE frame and does NOT close the stream.
// Stream closing is handled by KittyClient.Close().
func (c *KittyClient) SendBye() error {
	c.mu.Lock()
	stream := c.stream
	c.mu.Unlock()

	if stream == nil {
		return errors.New("stream is nil")
	}

	frame := protocol.BaseFrame{
		Type:  "BYE",
		MsgID: time.Now().UnixMilli(),
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	_, err = stream.Write(b)
	return err
}
