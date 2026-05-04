// client/ack.go
package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/clientutils"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// startReceiverLoop listens for incoming frames and handles MEOW_OK, ERROR and DATA.
func startReceiverLoop(stream *quic.Stream, disconnected chan struct{}) (map[int64]chan struct{}, *sync.Mutex) {

	pending := make(map[int64]chan struct{})
	var mu sync.Mutex

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stream.Read(buf)
			if err != nil {
				fmt.Println("\n[Client] Connection closed by server:", err)
				fmt.Println("[System] Returning to disconnected state.")
				// Signal to main that the session has ended.
				select {
				case <-disconnected:
				default:
					close(disconnected)
				}
				return
			}

			// TASK 8: Lightweight detection of frame type.
			typeName, msgID, err := protocol.GetFrameType(buf[:n])
			if err != nil {
				fmt.Println("[Client] Parse error:", err)
				continue
			}

			switch typeName {
			case "MEOW_OK":
				mu.Lock()
				if ch, ok := pending[msgID]; ok {
					close(ch)
					delete(pending, msgID)
					fmt.Printf("\n[Delivered] msg_id=%d\n> ", msgID)
				}
				mu.Unlock()

			case "ERROR":
				var errFrame protocol.ErrorFrame
				if json.Unmarshal(buf[:n], &errFrame) == nil {
					fmt.Printf("\n[Server ERROR] %s: %s\n", errFrame.Code, errFrame.Desc)

					switch errFrame.Code {
					case "ERR_15":
						// Receiver is offline – we do NOT close the session, only inform the user.
						fmt.Println("[System] Receiver is offline. Messages will not be delivered.")
					}

					fmt.Print("> ")
				} else {
					fmt.Println("\n[Client: ack] Failed to parse ERROR frame\n> ")
				}

			case "DATA":
				var df protocol.DataFrame
				if json.Unmarshal(buf[:n], &df) == nil {
					// Minimal chat output – sender + payload.
					fmt.Printf("\n[Message from %s]: %s\n> ", df.Sender, df.Payload)
				} else {
					fmt.Println("\n[Client] Failed to parse DATA frame\n> ")
				}
			}
		}
	}()

	return pending, &mu
}

// sendMessage sends a DATA frame and starts a 5-second ACK timer.
func sendMessage(stream *quic.Stream, target, text string, pending map[int64]chan struct{}, mu *sync.Mutex) {
	safe := clientutils.TruncateMessage(text)
	msgID := time.Now().UnixMilli() // msg_id as timestamp

	ch := make(chan struct{})
	mu.Lock()
	pending[msgID] = ch
	mu.Unlock()

	// Timer handling missing acknowledgments (Task 14).
	clientutils.StartAckTimer(msgID, ch, func() {
		mu.Lock()
		if _, ok := pending[msgID]; ok {
			delete(pending, msgID)
			fmt.Printf("\n[Timeout] msg_id=%d not delivered\n> ", msgID)
		}
		mu.Unlock()
	})

	// TASK 10: Create a strongly typed DataFrame.
	frame := protocol.DataFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "DATA",
			MsgID: msgID,
		},
		Target:  target,
		Payload: safe,
		MAC:     "placeholder",
	}

	b, _ := json.Marshal(frame)
	stream.Write(b)
}
