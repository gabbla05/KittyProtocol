package main

import "fmt"

// handleBye processes the BYE frame.
// It removes the session from SessionManager and triggers cleanup.
func (c *clientContext) handleBye() {
	if c.session != nil {
		fmt.Println("[Handler: Bye] Cleaning up session for:", c.username)
		globalSessions.Remove(c.username)
		if c.session.CloseFunc != nil {
			c.session.CloseFunc()
		}
		c.session = nil
	}
}
