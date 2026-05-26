package protection

import "time"

// -----------------------------------------------------------------------------
// Global protection-related constants used across the Hub.
// Centralizing these values eliminates magic numbers and makes the system
// easier to configure, test, and maintain.
// -----------------------------------------------------------------------------

// DefaultSessionRateLimit defines how many DATA frames per second a single
// authenticated session is allowed to send. This protects the Hub from
// message flooding by malicious or buggy clients.
const DefaultSessionRateLimit = 10

// DefaultIdleCloseErrorCode is the QUIC application error code used when
// the Hub closes an idle connection. This is intentionally distinct from
// protocol-level error codes.
const DefaultIdleCloseErrorCode = 0x09

// DefaultSessionIdleTimeout defines how long a session may remain inactive
// before being automatically closed by the SessionManager.
const DefaultSessionIdleTimeout = 60 * time.Second

// DefaultSessionCleanupInterval defines how often the SessionManager scans
// for idle sessions.
const DefaultSessionCleanupInterval = 10 * time.Second

// DefaultAuthTimeout defines how long a client has to complete the
// HELLO → AUTH/REGISTER handshake before the Hub closes the connection.
const DefaultAuthTimeout = 2 * time.Minute

// Replay protection parameters.
const (
	// Maximum number of tracked message IDs before forced cleanup.
	MaxReplayEntries = 10_000

	// TTL for replay entries.
	ReplayTTL = 2 * time.Minute

	// Sweep interval for replay cleanup.
	ReplaySweepInterval = 5 * time.Second
)
