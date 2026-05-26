package secretstore

import "errors"

// Package‑level errors used across the secret store implementation.
var (
	ErrEmptyMasterKey = errors.New("master key is empty")
	ErrEmptyPeer      = errors.New("peer cannot be empty")
	ErrEmptySecret    = errors.New("secret cannot be empty")
	ErrShortCipher    = errors.New("ciphertext too short")
)
