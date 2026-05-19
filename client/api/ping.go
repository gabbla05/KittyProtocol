package api

import (
	"encoding/json"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// StartPingLoop launches a background goroutine that sends PING frames
// every 30 seconds. It stops when stopPing is closed or the stream ends.
func (c *KittyClient) StartPingLoop() {
	c.mu.Lock()
	stream := c.stream
	stop := c.stopPing
	c.mu.Unlock()

	if stream == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return

			case <-ticker.C:
				c.mu.Lock()
				s := c.stream
				c.mu.Unlock()

				if s == nil {
					return
				}

				frame := protocol.PingFrame{
					BaseFrame: protocol.BaseFrame{
						Type:  "PING",
						MsgID: time.Now().UnixMilli(),
					},
				}

				b, _ := json.Marshal(frame)
				_, err := s.Write(b)
				if err != nil {
					// Stream closed — stop loop
					return
				}
			}
		}
	}()
}
