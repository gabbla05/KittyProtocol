package api

// handleMeowOK processes a MEOW_OK frame. Depending on the current client
// state it completes the HELLO / AUTH / REGISTER handshake or forwards
// the ACK to the AckManager for application-level messages.
func (c *KittyClient) handleMeowOK(msgID int64) {
	c.mu.Lock()
	currentState := c.state
	helloCh := c.helloCh
	authCh := c.authCh
	registerCh := c.registerCh
	ackMgr := c.ackMgr
	c.mu.Unlock()

	switch currentState {
	case StateHandshaking:
		helloCh <- OpResult{OK: true}
		c.setState(StateAuthenticating)

	case StateAuthenticating:
		authCh <- OpResult{OK: true}
		c.setState(StateSelectingTarget)

	case StateRegistering:
		registerCh <- OpResult{OK: true}
		c.setState(StateAuthenticating)

	default:
		if ackMgr != nil {
			ackMgr.NotifyDelivered(msgID)
		}
	}
}
