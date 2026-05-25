package api

import "errors"

var (
	// ErrNotConnected is returned when an operation requires an active
	// connection/stream but the client is disconnected.
	ErrNotConnected = errors.New("client not connected")

	// ErrNoStream is returned when the underlying stream is nil.
	ErrNoStream = errors.New("stream is nil")

	// ErrNoSharedSecret is returned when trying to send or decrypt
	// E2EE data without a derived secret for the peer.
	ErrNoSharedSecret = errors.New("no shared secret for target")

	// ErrTargetNotSet is returned when an operation requires a target
	// but none is configured.
	ErrTargetNotSet = errors.New("target not set")
)
