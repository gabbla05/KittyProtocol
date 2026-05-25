package api

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// canonicalTarget normalizes the target username to a stable form.
// This must match the Hub's canonicalization logic.
func canonicalTarget(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}

// ensureConnected returns current stream or an error if the client
// is not in a usable state for sending frames.
func (c *KittyClient) ensureConnected() (StreamAdapter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stream == nil {
		return nil, ErrNoStream
	}
	if c.state == StateDisconnected {
		return nil, ErrNotConnected
	}
	return c.stream, nil
}

// SendAppFrameEncrypted sends an application-level frame (chat control, text, etc.)
// encrypted as a DATA frame. The Hub requires MAC for all DATA frames.
func (c *KittyClient) SendAppFrameEncrypted(target string, payload []byte) error {
	stream, err := c.ensureConnected()
	if err != nil {
		return err
	}

	if target == "" {
		return ErrTargetNotSet
	}

	c.mu.Lock()
	kEnc, kMac, ok := c.getKeysForPeer(target)
	c.mu.Unlock()

	if !ok {
		return ErrNoSharedSecret
	}

	msgID := time.Now().UnixMilli()

	c.mu.Lock()
	ackMgr := c.ackMgr
	if ackMgr != nil {
		ackMgr.AddPending(msgID)
	}
	c.mu.Unlock()

	canonTarget := canonicalTarget(target)

	payloadB64, macB64, err := cryptoee.EncryptAndMACWithKeys(
		msgID,
		canonTarget,
		string(payload),
		kEnc,
		kMac,
	)
	if err != nil {
		return err
	}

	frame := protocol.DataFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeData,
			MsgID: msgID,
		},
		Target:  canonTarget,
		Payload: payloadB64,
		MAC:     macB64,
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	// Store last raw frame for replay testing (dev-only helper).
	c.mu.Lock()
	c.lastFrame = b
	c.mu.Unlock()

	_, err = stream.Write(b)
	return err
}

// SendGetStatus sends a GET_STATUS frame for a given user.
func (c *KittyClient) SendGetStatus(target string) error {
	stream, err := c.ensureConnected()
	if err != nil {
		return err
	}

	if target == "" {
		return ErrTargetNotSet
	}

	msgID := time.Now().UnixMilli()

	frame := protocol.GetStatusFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeGetStatus,
			MsgID: msgID,
		},
		Target: canonicalTarget(target),
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
	stream, err := c.ensureConnected()
	if err != nil {
		// BYE jest „best-effort” — jeśli nie ma streama, nie robimy dramatu.
		if errors.Is(err, ErrNoStream) || errors.Is(err, ErrNotConnected) {
			return nil
		}
		return err
	}

	frame := protocol.BaseFrame{
		Type:  protocol.FrameTypeBye,
		MsgID: time.Now().UnixMilli(),
	}

	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	_, err = stream.Write(b)
	return err
}
