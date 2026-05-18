package main

import (
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

// clientContext stores all state related to a single client connection.
// It keeps track of:
// - QUIC connection and stream
// - authenticated session (after AUTH)
// - username
// - AUTH timeout timer
type clientContext struct {
	conn      *quic.Conn
	stream    *quic.Stream
	session   *protection.Session
	username  string
	authTimer *protection.AuthTimer
}

// cleanup is executed when the handler finishes.
// It removes the session from the SessionManager and stops the AUTH timer.
func (c *clientContext) cleanup() {
	if c.session != nil {
		fmt.Println("[Handler] Cleaning up session for:", c.username)
		globalSessions.Remove(c.username)
		if c.session.CloseFunc != nil {
			c.session.CloseFunc()
		}
	}
	if c.authTimer != nil {
		c.authTimer.Stop()
	}
}

// touch updates the session's LastActive timestamp.
// This is used for idle timeout detection.
func (c *clientContext) touch() {
	if c.session != nil {
		c.session.LastActive = time.Now()
	}
}
