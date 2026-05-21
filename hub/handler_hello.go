package main

// handleHello processes the HELLO frame.
// It sends MEOW_OK("Ready for auth") and starts the AUTH timeout timer.
func (c *clientContext) handleHello() {
	c.authTimer = handleHELLO(c.stream, c.conn)
}
