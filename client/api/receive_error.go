package api

import "github.com/gabbla05/KittyProtocol/protocol"

// handleErrorFrame processes an ERROR frame. For handshake-related states
// it forwards the error to the appropriate result channel. For established
// sessions it forwards the error to the registered ErrorHandler, if any.
func (c *KittyClient) handleErrorFrame(ef protocol.ErrorFrame) {
	c.mu.Lock()
	currentState := c.state
	helloCh := c.helloCh
	authCh := c.authCh
	registerCh := c.registerCh
	eh := c.errHandler
	c.mu.Unlock()

	switch currentState {
	case StateHandshaking:
		helloCh <- OpResult{OK: false, Code: ef.Code, Desc: ef.Desc}

	case StateAuthenticating:
		authCh <- OpResult{OK: false, Code: ef.Code, Desc: ef.Desc}

	case StateRegistering:
		registerCh <- OpResult{OK: false, Code: ef.Code, Desc: ef.Desc}

	default:
		if eh != nil {
			eh(ef.Code, ef.Desc)
		} else {
			log(LogError, "server error %s: %s", ef.Code, ef.Desc)
		}
	}
}
