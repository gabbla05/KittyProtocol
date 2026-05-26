package api

import (
	"encoding/json"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// SendAppFrameEncrypted sends an application-level frame (chat control, text, etc.)
// encrypted as a DATA frame. The Hub requires MAC for all DATA frames.
//
// SECURITY:
//   - Uses per-peer E2EE keys derived from a shared secret.
//   - Payload is encrypted and authenticated (MAC) via internal/cryptoee.
func (c *KittyClient) SendAppFrameEncrypted(target string, payload []byte) error {
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

	c.mu.Lock()
	kEnc, kMac, ok := c.getKeysForPeer(target)
	ackMgr := c.ackMgr
	c.mu.Unlock()

	if !ok {
		return ErrNoSharedSecret
	}

	msgID := time.Now().UnixMilli()

	if ackMgr != nil {
		ackMgr.AddPending(msgID)
	}

	canonTarget := canonicalTarget(target)

	if len(payload) > maxPayloadSize {
		return ErrPayloadTooLarge
	}

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

	// Store last raw frame for replay testing (dev-only helper).
	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.lastFrame = b
	c.mu.Unlock()

	_, err = stream.Write(b)
	return err
}
