package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendAppFrameEncrypted sends an application-level frame (chat control, text, etc.)
// encrypted as a DATA frame. The Hub requires MAC for all DATA frames.
func (c *KittyClient) SendAppFrameEncrypted(target string, payload []byte) error {
	c.mu.Lock()
	stream := c.stream
	ackMgr := c.ackMgr

	var (
		kEnc []byte
		kMac []byte
		ok   bool
	)
	if c.peerKeys != nil {
		pk, exists := c.peerKeys[target]
		if exists {
			kEnc = pk.kEnc
			kMac = pk.kMac
			ok = true
		}
	}
	c.mu.Unlock()

	if !ok {
		return errors.New("no shared secret for target")
	}
	if stream == nil {
		return errors.New("stream is nil")
	}
	if target == "" {
		return errors.New("target not set")
	}

	msgID := time.Now().UnixMilli()

	if ackMgr != nil {
		ackMgr.AddPending(msgID)
	}

	payloadB64, macB64, err := cryptoee.EncryptAndMACWithKeys(
		msgID,
		target,
		string(payload),
		kEnc,
		kMac,
	)
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

	_, err = stream.Write(b)
	return err
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

// SendBye sends a BYE frame to the Hub.
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
