package api

import (
	"encoding/json"
	"fmt"

	"github.com/gabbla05/KittyProtocol/internal/cryptoee"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// StartReceiverLoop starts a background goroutine that reads frames from the stream
// and handles MEOW_OK, ERROR, DATA and STATUS_RES.
// It stops when:
// - the stream is closed,
// - the disconnected channel is closed,
// - stopRecv is closed via client.Close().
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
				// Delivery acknowledgment
				if ackMgr != nil {
					ackMgr.NotifyDelivered(msgID)
				}

			case "ERROR":
				var errFrame protocol.ErrorFrame
				if json.Unmarshal(buf[:n], &errFrame) == nil {
					fmt.Printf("\n[Client: Receive] Server ERROR %s: %s\n> ", errFrame.Code, errFrame.Desc)
					if errFrame.Code == "ERR_15" {
						fmt.Println("[Client: Receive] Receiver is offline. Messages will not be delivered.")
					}
				} else {
					fmt.Println("\n[Client: Receive] Failed to parse ERROR frame\n> ")
				}

			case "DATA":
				c.mu.Lock()
				kEnc := c.kEnc
				kMac := c.kMac
				c.mu.Unlock()

				if kEnc == nil || kMac == nil {
					fmt.Println("\n[Client] No shared secret set — cannot decrypt.\n> ")
					continue
				}

				var df protocol.DataFrame
				if json.Unmarshal(buf[:n], &df) != nil {
					fmt.Println("\n[Client: Receive] Failed to parse DATA frame\n> ")
					continue
				}

				// Client-side replay protection (silent drop)
				if replay != nil && replay.MarkAndCheck(df.MsgID) {
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

				fmt.Printf("\n[Client: Receive] Message from %s: %s\n> ", df.Sender, plaintext)

			case "STATUS_RES":
				var sf protocol.StatusResFrame
				if json.Unmarshal(buf[:n], &sf) == nil {
					fmt.Printf("\n[Client: Receive] %s is %s\n> ", sf.Target, sf.Status)
				} else {
					fmt.Println("\n[Client: Receive] Failed to parse STATUS_RES frame\n> ")
				}
			}
		}
	}()
}
