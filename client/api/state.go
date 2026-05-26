package api

// State returns the current client state.
//
// Thread-safety: safe to call from any goroutine.
func (c *KittyClient) State() ClientState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// setState updates the internal state.
//
// This method is intentionally unexported to prevent UI layers from
// mutating the protocol state machine directly.
func (c *KittyClient) setState(newState ClientState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = newState
}
