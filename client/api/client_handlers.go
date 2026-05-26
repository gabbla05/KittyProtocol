package api

// RegisterAppPayloadHandler registers a callback for application-level
// payloads that are not recognized as chat control frames.
//
// Thread-safety: safe to call at any time; handler replacement is atomic.
func (c *KittyClient) RegisterAppPayloadHandler(h AppPayloadHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.appHandler = h
}

// OnError registers a callback for protocol-level errors that are not
// directly tied to HELLO / AUTH / REGISTER operations.
func (c *KittyClient) OnError(h ErrorHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errHandler = h
}

// OnStatus registers a callback for presence/status updates.
func (c *KittyClient) OnStatus(h StatusHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statusHandler = h
}

// OnDisconnected registers a callback that is invoked when the underlying
// transport is broken or the receiver loop terminates due to a read error.
func (c *KittyClient) OnDisconnected(h DisconnectHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.disconnectHandler = h
}

// User returns the currently authenticated username, if any.
// If the client is not authenticated, an empty string is returned.
func (c *KittyClient) User() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.user
}

// getKeysForPeer returns the E2EE keys for a given peer, if present.
// It is intentionally unexported; key management is internal to KittyClient.
func (c *KittyClient) getKeysForPeer(peer string) (kEnc, kMac []byte, ok bool) {
	pk, exists := c.peerKeys[peer]
	if !exists {
		return nil, nil, false
	}
	return pk.kEnc, pk.kMac, true
}
