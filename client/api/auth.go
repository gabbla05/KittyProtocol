package api

import (
	"encoding/json"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// -----------------------------------------------------------------------------
// AUTH / REGISTER (asynchronous)
// -----------------------------------------------------------------------------

// SendAuth sends an AUTH frame with the provided credentials.
//
// BEHAVIOR:
//   - Requires an established QUIC connection (ensureConnected).
//   - Updates internal client state to StateAuthenticating.
//   - Stores the username on the client instance.
//   - The result is delivered asynchronously via AuthResult().
func (c *KittyClient) SendAuth(user, pass string) error {
	return c.sendAuthLikeFrame(protocol.FrameTypeAuth, user, pass, StateAuthenticating, func() {
		c.user = user
	})
}

// AuthResult returns a read-only channel that delivers the result
// of the last AUTH operation.
func (c *KittyClient) AuthResult() <-chan OpResult {
	return c.authCh
}

// SendRegister sends a REGISTER frame with the provided credentials.
//
// BEHAVIOR:
//   - Requires an established QUIC connection (ensureConnected).
//   - Updates internal client state to StateRegistering.
//   - The result is delivered asynchronously via RegisterResult().
func (c *KittyClient) SendRegister(user, pass string) error {
	return c.sendAuthLikeFrame(protocol.FrameTypeRegister, user, pass, StateRegistering, nil)
}

// RegisterResult returns a read-only channel that delivers the result
// of the last REGISTER operation.
func (c *KittyClient) RegisterResult() <-chan OpResult {
	return c.registerCh
}

// sendAuthLikeFrame is a shared helper for AUTH and REGISTER flows.
//
// PARAMETERS:
//   - frameType: protocol.FrameTypeAuth or protocol.FrameTypeRegister.
//   - user, pass: credentials to send.
//   - nextState: client state to set before sending.
//   - beforeUnlock: optional hook executed under lock before releasing it
//     (e.g. to store username).
func (c *KittyClient) sendAuthLikeFrame(
	frameType string,
	user, pass string,
	nextState ClientState,
	beforeUnlock func(),
) error {
	stream, err := c.ensureConnected()
	if err != nil {
		return err
	}

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  frameType,
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
	if beforeUnlock != nil {
		beforeUnlock()
	}
	c.state = nextState
	c.mu.Unlock()

	_, err = stream.Write(b)
	return err
}

// -----------------------------------------------------------------------------
// HELLO (asynchronous)
// -----------------------------------------------------------------------------

// HelloResult returns a read-only channel that delivers the result
// of the initial HELLO handshake.
func (c *KittyClient) HelloResult() <-chan OpResult {
	return c.helloCh
}
