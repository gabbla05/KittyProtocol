package api

// SetTarget sets the current chat target.
// This is a UI-level decision and does not involve protocol logic.
func (c *KittyClient) SetTarget(target string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.target = target
}

// Target returns the current chat target.
func (c *KittyClient) Target() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.target
}
