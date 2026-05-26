package api

import "time"

// Centralized constants for KittyClient behavior.
// This keeps magic numbers out of the logic files and makes tuning easier.
const (
	// defaultRecvBufferSize is the buffer size used for reading frames
	// from the Hub in the receiver loop.
	defaultRecvBufferSize = 4096

	// defaultAckTimeout is the time after which a pending message is
	// considered undelivered if no MEOW_OK arrives.
	defaultAckTimeout = 5 * time.Second

	// defaultPingInterval controls how often PING frames are sent to the Hub.
	defaultPingInterval = 30 * time.Second

	// minSharedSecretLength defines the minimum number of bytes required
	// for a valid E2EE shared secret.
	minSharedSecretLength = 16

	// maxPayloadSize defines the maximum allowed plaintext payload size
	// before encryption. This protects memory usage and prevents abuse.
	maxPayloadSize = 16 * 1024 // 16 KB

	// maxUsernameLength defines the maximum allowed username length.
	maxUsernameLength = 64
)
