package api

import (
	"encoding/json"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// handleStatusResFrame processes a STATUS_RES frame and forwards the
// status information either to a registered StatusHandler or to the log.
func (c *KittyClient) handleStatusResFrame(frameBytes []byte) {
	var sf protocol.StatusResFrame
	if json.Unmarshal(frameBytes, &sf) != nil {
		log(LogError, "failed to parse STATUS_RES frame")
		return
	}

	c.mu.Lock()
	sh := c.statusHandler
	c.mu.Unlock()

	if sh != nil {
		sh(sf.Target, sf.Status)
	} else {
		log(LogInfo, "status: %s is %s", sf.Target, sf.Status)
	}
}
