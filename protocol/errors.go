package protocol

// Error codes defined by the KittyProtocol specification.
// These MUST remain stable across versions, as both Hub and Clients rely on them.
//
// Each error code corresponds to a transport‑level or protocol‑level failure.
// Application‑level errors SHOULD reuse these codes and place human‑readable
// details in the Desc field of ErrorFrame.

const (
	// ERR_01 — Protocol Violation
	// Example: sending DATA before AUTH.
	ErrProtocolViolation = "ERR_01"

	// ERR_02 — Format Error
	// Invalid JSON, missing type/msg_id, wrong types.
	ErrFormatError = "ERR_02"

	// ERR_03 — Authorization Timeout
	// No AUTH received within 20 seconds after HELLO.
	ErrAuthorizationTimeout = "ERR_03"

	// ERR_04 — Authentication Failed
	// Wrong username or password.
	ErrAuthenticationFailed = "ERR_04"

	// ERR_05 — Session Error
	// Undefined session error; client should re‑authenticate.
	ErrSessionError = "ERR_05"

	// ERR_06 — Replay Detected
	// Reuse of msg_id.
	ErrReplayDetected = "ERR_06"

	// ERR_07 — Rate Limit Exceeded
	// >10 messages/s or >100/min.
	ErrRateLimitExceeded = "ERR_07"

	// ERR_08 — Delivery Failed – Recipient Offline
	// Target user is not online.
	ErrDeliveryFailedOffline = "ERR_08"

	// ERR_09 — Session Timeout (Idle)
	// No activity for 60 seconds.
	ErrSessionTimeoutIdle = "ERR_09"

	// ERR_10 — Resource Exhaustion
	// Hub overloaded; retry with jitter.
	ErrResourceExhaustion = "ERR_10"

	// ERR_11 — Internal Server Error
	// Hub internal failure (e.g., DB offline).
	ErrInternalServerError = "ERR_11"

	// ERR_12 — Version Mismatch
	// Unsupported protocol version.
	ErrVersionMismatch = "ERR_12"

	// ERR_13 — Payload Too Large
	// Payload > 2048 bytes.
	ErrPayloadTooLarge = "ERR_13"

	// ERR_14 — Not Authorized
	// User lacks permission for the action.
	ErrNotAuthorized = "ERR_14"

	// ERR_15 — Unknown Target
	// Example: GET_KEY for non‑existent user.
	ErrUnknownTarget = "ERR_15"
)
