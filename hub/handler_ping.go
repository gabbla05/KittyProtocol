package main

// handlePing updates session activity timestamp.
// This is used for idle timeout detection.
func (c *clientContext) handlePing() {
	c.touch()
}
