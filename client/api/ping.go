package api

import (
	"encoding/json"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
)

// StartPingLoop starts a background goroutine that periodically sends
// PING frames to keep the KittyProtocol session active.
//
// Although QUIC has its own keep-alive mechanisms, the Hub expects
// application-level PING frames to detect idle clients.
//
// The loop terminates when:
//   - stopPing channel is closed,
//   - the underlying QUIC stream becomes nil,
//   - writing to the stream returns an error.
func (c *KittyClient) StartPingLoop() {
	c.mu.Lock()
	stream := c.stream
	stop := c.stopPing
	c.mu.Unlock()

	if stream == nil || stop == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(defaultPingInterval)
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
						Type:  protocol.FrameTypePing,
						MsgID: time.Now().UnixMilli(),
					},
				}

				b, err := json.Marshal(frame)
				if err != nil {
					// Serialization failure — stop ping loop.
					return
				}

				if _, err := s.Write(b); err != nil {
					// Stream closed or write error — stop ping loop.
					return
				}
			}
		}
	}()
}
