package api

import (
	"encoding/json"
	"strings"
)

// canonicalTarget normalizes a username into a canonical form.
// This must match the Hub's canonicalization logic.
func canonicalTarget(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}

// ensureConnected returns the active QUIC stream or an error if the client
// is not in a valid state for sending frames.
func (c *KittyClient) ensureConnected() (StreamAdapter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stream == nil {
		return nil, ErrNoStream
	}
	if c.state == StateDisconnected {
		return nil, ErrNotConnected
	}
	return c.stream, nil
}

// sendFrame marshals the given frame to JSON and writes it to the stream.
// This helper centralizes JSON encoding and write logic to avoid duplication.
func (c *KittyClient) sendFrame(stream StreamAdapter, frame any) error {
	b, err := json.Marshal(frame)
	if err != nil {
		return err
	}
	_, err = stream.Write(b)
	return err
}
