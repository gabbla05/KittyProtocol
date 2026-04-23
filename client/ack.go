package main

import (
	"encoding/json" // Dodano do obsługi JSON
	"fmt"
	"sync"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/clientutils"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// startReceiverLoop słucha przychodzących ramek i obsługuje MEOW_OK oraz ERROR.
func startReceiverLoop(stream *quic.Stream) (map[int64]chan struct{}, *sync.Mutex) {
	pending := make(map[int64]chan struct{})
	var mu sync.Mutex

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stream.Read(buf)
			if err != nil {
				fmt.Println("\n[Client] Connection closed:", err)
				return
			}

			// TASK 8: Wstępne rozpoznanie typu ramki [cite: 2185]
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
				// TASK 10: Parsowanie do dedykowanej struktury ErrorFrame [cite: 2151, 2396]
				var errFrame protocol.ErrorFrame
				if json.Unmarshal(buf[:n], &errFrame) == nil {
					fmt.Printf("\n[Server ERROR] %s: %s\n> ", errFrame.Code, errFrame.Desc)
				}
			}
		}
	}()

	return pending, &mu
}

// sendMessage wysyła ramkę DATA i uruchamia 5-sekundowy timer ACK.
func sendMessage(stream *quic.Stream, target, text string, pending map[int64]chan struct{}, mu *sync.Mutex) {
	safe := clientutils.TruncateMessage(text)
	msgID := time.Now().UnixMilli() // msg_id jako timestamp [cite: 2146]

	ch := make(chan struct{})
	mu.Lock()
	pending[msgID] = ch
	mu.Unlock()

	// Timer obsługujący brak potwierdzenia (Task 14) [cite: 2267, 3283]
	clientutils.StartAckTimer(msgID, ch, func() {
		mu.Lock()
		if _, ok := pending[msgID]; ok {
			delete(pending, msgID)
			fmt.Printf("\n[Timeout] msg_id=%d not delivered\n> ", msgID)
		}
		mu.Unlock()
	})

	// TASK 10: Utworzenie twardej struktury DataFrame [cite: 2153, 2731]
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
