// hub/context.go
// Connection-scoped state for a single QUIC client.
// Each clientContext instance represents one connected client and tracks:
//   - QUIC connection and stream
//   - authentication state
//   - associated session (after AUTH)
//   - AUTH timeout timer
//   - last activity timestamps
//
// This file contains no protocol logic — only state management.

package hub

import (
	"time"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

// connectionState represents the current stage of the protocol handshake.
type connectionState int

const (
	stateInit          connectionState = iota // No HELLO yet
	stateHelloReceived                        // HELLO received, waiting for AUTH/REGISTER
	stateAuthenticated                        // AUTH successful, session active
)

// clientContext holds all per-connection state.
// It is created in dispatcher.go when a new QUIC stream is accepted.
type clientContext struct {
	conn      *quic.Conn            // Underlying QUIC connection
	stream    *quic.Stream          // Primary bidirectional stream
	session   *protection.Session   // Active session after AUTH
	username  string                // Set only after successful AUTH
	state     connectionState       // Protocol handshake state machine
	authTimer *protection.AuthTimer // Timer to enforce AUTH timeout
}

// cleanup releases all resources associated with this client.
// It is ALWAYS called via defer in dispatcher.go.
func (c *clientContext) cleanup() {
	// Remove session if present
	if c.session != nil {
		logInfo("[Context] Cleaning up session for user: %s", c.username)
		globalSessions.Remove(c.username)

		if c.session.CloseFunc != nil {
			c.session.CloseFunc()
		}

		c.session = nil
	}

	// Stop AUTH timeout timer
	if c.authTimer != nil {
		c.authTimer.Stop()
		c.authTimer = nil
	}
}

// touch updates the session's last activity timestamp.
// Called on every valid frame after AUTH.
func (c *clientContext) touch() {
	if c.session != nil {
		c.session.LastActive = time.Now()
	}
}
