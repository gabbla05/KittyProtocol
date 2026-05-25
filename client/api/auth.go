package api

import (
	"encoding/json"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// -----------------------------------------------------------------------------
// AUTH (asynchronous)
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
	c.state = StateAuthenticating
	c.mu.Unlock()

	_, err = stream.Write(b)
	return err
}

func (c *KittyClient) AuthResult() <-chan OpResult {
	return c.authCh
}

// -----------------------------------------------------------------------------
// REGISTER (asynchronous)
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

	c.mu.Lock()
	c.state = StateRegistering
	c.mu.Unlock()

	_, err = stream.Write(b)
	return err
}

func (c *KittyClient) RegisterResult() <-chan OpResult {
	return c.registerCh
}

// -----------------------------------------------------------------------------
// HELLO (asynchronous)
// -----------------------------------------------------------------------------

func (c *KittyClient) HelloResult() <-chan OpResult {
	return c.helloCh
}
