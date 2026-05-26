package api

import "errors"

// Transport / connection errors
var (
	// ErrNotConnected is returned when an operation requires an active
	// connection/stream but the client is disconnected.
	ErrNotConnected = errors.New("client not connected")

	// ErrNoStream is returned when the underlying stream is nil.
	ErrNoStream = errors.New("stream is nil")
)

// Target / peer errors
var (
	// ErrTargetNotSet is returned when an operation requires a target
	// but none is configured.
	ErrTargetNotSet = errors.New("target not set")

	// ErrTargetNameTooLong is returned when the configured target name
	// exceeds the maximum allowed length.
	ErrTargetNameTooLong = errors.New("target name too long")

	// ErrPeerNameTooLong is returned when the identifier of a peer
	// exceeds the maximum allowed length.
	ErrPeerNameTooLong = errors.New("peer name too long")

	// ErrSharedSecretTooShort is returned when the provided shared secret
	// does not meet the minimum length requirement.
	ErrSharedSecretTooShort = errors.New("shared secret too short")

	// ErrNoSharedSecret is returned when trying to send or decrypt
	// E2EE data without a derived secret for the peer.
	ErrNoSharedSecret = errors.New("no shared secret for target")

	// ErrEmptyPeer is returned when SetSharedSecretForPeer is called
	// with an empty peer identifier.
	ErrEmptyPeer = errors.New("peer cannot be empty")

	// ErrEmptySecret is returned when SetSharedSecretForPeer is called
	// with an empty secret.
	ErrEmptySecret = errors.New("secret cannot be empty")
)

// Frame / protocol errors
var (
	// ErrPayloadTooLarge is returned when the frame payload size
	// exceeds the protocol's defined maximum limit.
	ErrPayloadTooLarge = errors.New("payload too large")

	// ErrUnknownFrameType indicates that the received frame type is not recognized.
	ErrUnknownFrameType = errors.New("unknown frame type")

	// ErrFrameParseFailed indicates that a frame could not be unmarshaled.
	ErrFrameParseFailed = errors.New("failed to parse frame")
)
