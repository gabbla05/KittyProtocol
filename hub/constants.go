package hub

import "time"

// Hub-level constants used across handlers, routing logic, and session management.
// These values complement protocol-level constants defined in protocol/error_codes.go.
// They MUST remain consistent with the protocol documentation and the protection package.

//===============================================================================
// NOTE:																		|
// These constants are part of the KittyProtocol specification.					|
// They may not be referenced directly inside the Hub because their enforcement	|
// happens in the protection package (rate limiting, replay, idle timeout)		|
// or on the client side (payload size).									    |
//																				|
// DO NOT REMOVE — they are required for protocol stability, documentation,		|
// and future compatibility with clients and tests.								|
// ==============================================================================

// Maximum allowed payload size for DATA frames.
// If exceeded, the Hub must return ERR_13 (Payload Too Large).
const maxPayloadBytes = 2048

// Size of the read buffer for incoming QUIC stream data.
// Must be large enough to hold the largest expected frame.
const readBufferSize = 8192

// Idle timeout for authenticated sessions.
// If exceeded, the Hub should return ERR_09 (Session Timeout / Idle).
// NOTE: This is separate from the QUIC idle timeout.
const sessionIdleTimeout = 60 * time.Second

// Maximum number of messages per second allowed by the Hub.
// Exceeding this should trigger ERR_07 (Rate Limit Exceeded).
const rateLimitPerSecond = 10

// Maximum number of messages per minute allowed by the Hub.
// Also part of ERR_07 enforcement.
const rateLimitPerMinute = 100

// Maximum number of consecutive format errors (ERR_02) before the Hub closes the connection.
// Helps prevent malformed spam or broken clients.
const maxFormatErrors = 5

// QUIC configuration defaults.
// These values define transport-level behavior and should remain aligned with
// performance and security expectations of the KittyProtocol Hub.
const (
	quicMaxIdleTimeout  = 60 * time.Second
	quicKeepAlivePeriod = 30 * time.Second
	quicAllow0RTT       = true
	quicDisablePMTU     = false
)

// Default listening address used when KITTY_INTERCEPT_ADDR is not set.
const defaultHubAddress = "0.0.0.0:9999"

// Default PostgreSQL DSN for Hub authentication backend.
// In production, this should be overridden via environment variables.
const defaultDSN = "postgres://kitty:kittypass@localhost:5432/kittyhub?sslmode=disable"
