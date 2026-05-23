package api

import (
	"encoding/json"
	"fmt"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// StartReceiverLoop launches a background goroutine responsible for reading
// all incoming frames from the QUIC stream. This is the only reader for the
// stream; all other components must communicate through higher‑level APIs.
//
// The loop terminates when:
//   - stopRecv is closed,
//   - the QUIC stream returns an error,
//   - the disconnected channel is closed (exactly once).
func (c *KittyClient) StartReceiverLoop(disconnected chan struct{}) {
	c.mu.Lock()
	stream := c.stream
	replay := c.replay
	ackMgr := c.ackMgr
	stopRecv := c.stopRecv
	c.mu.Unlock()

	if stream == nil {
		return
	}

	go func() {
		buf := make([]byte, 4096)

		for {
			select {
			case <-stopRecv:
				return
			default:
			}

			n, err := stream.Read(buf)
			if err != nil {
				fmt.Println("\n[Client: Receive] Connection closed by server:", err)
				fmt.Println("[Client: Receive] Returning to disconnected state.")

				select {
				case <-disconnected:
				default:
					close(disconnected)
				}
				return
			}

			typeName, msgID, err := protocol.GetFrameType(buf[:n])
			if err != nil {
				fmt.Println("[Client: Receive] Parse error:", err)
				continue
			}

			switch typeName {

			case "MEOW_OK":
				if ackMgr != nil {
					ackMgr.NotifyDelivered(msgID)
				}

			case "ERROR":
				var errFrame protocol.ErrorFrame
				if json.Unmarshal(buf[:n], &errFrame) == nil {
					fmt.Printf("\n[Client: Receive] Server ERROR %s: %s\n> ",
						errFrame.Code, errFrame.Desc)
				} else {
					fmt.Println("\n[Client: Receive] Failed to parse ERROR frame\n> ")
				}

			case "DATA":
				var df protocol.DataFrame
				if json.Unmarshal(buf[:n], &df) != nil {
					fmt.Println("\n[Client: Receive] Failed to parse DATA frame\n> ")
					continue
				}

				// Replay protection
				if replay != nil && replay.MarkAndCheck(df.MsgID) {
					continue
				}

				// Retrieve handler + keys
				c.mu.Lock()
				handler := c.appHandler

				var (
					kEnc []byte
					kMac []byte
					ok   bool
				)
				if c.peerKeys != nil {
					pk, exists := c.peerKeys[df.Sender]
					if exists {
						kEnc = pk.kEnc
						kMac = pk.kMac
						ok = true
					}
				}
				c.mu.Unlock()

				if !ok {
					fmt.Printf("\n[Client] No shared secret for sender %s — cannot decrypt.\n> ",
						df.Sender)
					continue
				}

				plaintext, err := cryptoee.DecryptAndVerifyWithKeys(
					df.MsgID,
					df.Target,
					df.Payload,
					df.MAC,
					kEnc,
					kMac,
				)
				if err != nil {
					fmt.Printf("\n[Client: Receive] E2EE error: %v\n> ", err)
					continue
				}

				if handler != nil {
					handler(df.Sender, []byte(plaintext))
				} else {
					fmt.Printf("\n[Client: Receive] Message from %s: %s\n> ",
						df.Sender, string(plaintext))
				}

			case "STATUS_RES":
				var sf protocol.StatusResFrame
				if json.Unmarshal(buf[:n], &sf) != nil {
					fmt.Println("\n[Client: Receive] Failed to parse STATUS_RES frame\n> ")
					continue
				}

				if sf.Target == "" && sf.Status == "no_target" {
					fmt.Printf("\n[Client: Receive] Chat ended. No active target.\n> ")
					continue
				}

				fmt.Printf("\n[Client: Receive] %s is %s\n> ", sf.Target, sf.Status)
			}
		}
	}()
}
