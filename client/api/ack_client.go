package api

// RegisterAckHandler registers a UI or application component to receive ACK events.
//
// Handlers are invoked synchronously in the caller goroutine.
// It is the caller's responsibility to offload heavy work to separate goroutines.
func (c *KittyClient) RegisterAckHandler(h AckEventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ackMgr != nil {
		c.ackMgr.RegisterHandler(h)
	}
}
