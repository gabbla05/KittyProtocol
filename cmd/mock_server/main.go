package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func main() {
	// Korzystamy z Twoich certyfikatów
	cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		fmt.Printf("Błąd certyfikatów: %v\n", err)
		return
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"kitty-quic-v1"},
	}

	// Słuchamy na 9999
	listener, err := quic.ListenAddr("127.0.0.1:9999", tlsConf, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("🤖 Kitty Mock Server (Task 13) działa na 127.0.0.1:9999")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			continue
		}
		// FIX: conn jest już wskaźnikiem (*quic.Conn), więc nie dajemy &
		go handleMockClient(conn)
	}
}

func handleMockClient(conn *quic.Conn) {
	// AcceptStream zgodnie z Twoim stylem z hub/handler.go
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return
	}
	defer stream.Close()

	buf := make([]byte, 4096)
	for {
		n, err := stream.Read(buf)
		if err != nil {
			return
		}

		// TASK 13: Zwracamy MEOW_OK z tym samym MsgID
		_, msgID, err := protocol.GetFrameType(buf[:n])
		if err != nil {
			msgID = 0
		}

		response := protocol.MeowOkFrame{
			BaseFrame: protocol.BaseFrame{
				Type:  "MEOW_OK",
				MsgID: msgID,
			},
			Status: "MOCK_OK_NO_ERRORS",
		}

		respBytes, _ := json.Marshal(response)
		stream.Write(respBytes)
		fmt.Printf("Replied MEOW_OK for MsgID: %d\n", msgID)
	}
}
