package api

// SetTarget sets the current chat target (peer username).
// This is a UI-level decision and does not involve protocol logic.
func (c *KittyClient) SetTarget(target string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(target) > maxUsernameLength {
		// UI-level error, nie protokołowy
		log(LogWarn, "target name too long")
	}

	c.target = target
}

// Target returns the current chat target (peer username).
func (c *KittyClient) Target() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.target
}
