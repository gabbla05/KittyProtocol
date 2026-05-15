package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// sendBye wysyła poprawną ramkę BYE z niezerowym msg_id,
// żeby GetFrameType ją zaakceptował.
func sendBye(stream *quic.Stream) {
	frame := protocol.BaseFrame{
		Type:  "BYE",
		MsgID: time.Now().UnixMilli(), // <--- KLUCZOWA ZMIANA
	}

	b, err := json.Marshal(frame)
	if err != nil {
		fmt.Println("[Client] Failed to marshal BYE frame:", err)
		return
	}

	if _, err := stream.Write(b); err != nil {
		fmt.Println("[Client] Failed to send BYE frame:", err)
		return
	}

	fmt.Println("[Client] BYE RAW:", string(b))
	fmt.Println("[Client] BYE stream ID:", stream.StreamID())
	fmt.Println("[Client] sendBye method has ended correctly")
}
