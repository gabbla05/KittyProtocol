package main

import (
	"encoding/json" // Dodano do obsługi json.Marshal
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// startPingLoop okresowo wysyła ramki PING, aby podtrzymać sesję (Keep-alive).
func startPingLoop(stream *quic.Stream) {
	go func() {
		for {
			// Zgodnie z dokumentacją wysyłamy PING co 30 sekund[cite: 2, 7].
			time.Sleep(30 * time.Second)

			// TASK 10: Użycie dedykowanej struktury PingFrame zamiast UniversalFrame.
			ping := protocol.PingFrame{
				BaseFrame: protocol.BaseFrame{
					Type:  "PING",
					MsgID: time.Now().UnixMilli(),
				},
			}

			// Serializacja przy użyciu standardowej biblioteki.
			b, _ := json.Marshal(ping)
			_, err := stream.Write(b)
			if err != nil {
				// Jeśli strumień jest zamknięty, kończymy pętlę.
				return
			}
		}
	}()
}
