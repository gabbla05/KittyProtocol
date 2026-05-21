package api

// ====================================================
// DELETE THIS FILE AND ITS USAGE FOR PRODCUTION BUILDS
// ====================================================

import "fmt"

// ReplayLastFrame resends the last raw frame written to the stream.
// This is used exclusively for replay protection testing and should not be
// exposed in production builds.
func (c *KittyClient) ReplayLastFrame() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lastFrame == nil {
		return fmt.Errorf("no frame to replay")
	}
	if c.stream == nil {
		return fmt.Errorf("stream is nil")
	}

	_, err := c.stream.Write(c.lastFrame)
	return err
}
