package api

// State returns the current client state (thread-safe).
func (c *KittyClient) State() ClientState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// setState updates the internal state.
// This method is intentionally unexported to prevent misuse by UI layers.
func (c *KittyClient) setState(newState ClientState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = newState
}
