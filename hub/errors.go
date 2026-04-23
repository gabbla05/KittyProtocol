package main

import (
	"encoding/json" // Dodano dla json.Marshal
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// sendError wysyła ustandaryzowaną ramkę ERROR do klienta.
func sendError(stream *quic.Stream, code, desc string) {
	// TASK 10: Użycie ErrorFrame zamiast UniversalFrame
	errFrame := protocol.ErrorFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "ERROR",
			MsgID: time.Now().UnixMilli(),
		},
		Code: code,
		Desc: desc,
	}

	// Serializacja przy użyciu standardowej biblioteki
	b, _ := json.Marshal(errFrame)
	stream.Write(b)
}
