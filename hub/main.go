package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// Mapa aktywnych sesji: User -> Strumień danych
var activeSessions = make(map[string]*quic.Stream)

func main() {
	// Ścieżki do certyfikatów muszą być poprawne! [cite: 577]
	cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		log.Fatal("TLS Cert loading error:", err)
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"kitty-quic"},
	}

	listener, err := quic.ListenAddr("127.0.0.1:9999", tlsConf, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("🐈 KittyProtocol Hub (Router) is listening on 127.0.0.1:9999...")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn *quic.Conn) {

	for {
		// Akceptujemy strumień wewnątrz połączenia
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return
		}

		go func(s *quic.Stream) {
			defer s.Close()
			buf := make([]byte, 4096)
			for {
				n, err := s.Read(buf)
				if err != nil {
					return
				}

				frame, err := protocol.ParseFrame(buf[:n])
				if err != nil {
					// ERR_02 Format Error [cite: 584, 433]
					resp := protocol.UniversalFrame{Type: "ERROR", Code: "ERR_02", Desc: err.Error()}
					s.Write(resp.ToJSON())
					continue
				}

				switch frame.Type {
				case "HELLO":
					resp := protocol.UniversalFrame{Type: "MEOW_OK", MsgID: frame.MsgID, Status: "Ready for auth"}
					s.Write(resp.ToJSON())

				case "AUTH":
					// Task 6: Mock DB [cite: 659]
					activeSessions[frame.User] = s
					resp := protocol.UniversalFrame{Type: "MEOW_OK", MsgID: frame.MsgID, Status: "Logged in"}
					s.Write(resp.ToJSON())
					fmt.Printf("[Hub] User %s logged in\n", frame.User)

				case "DATA":
					// Message Broker Logic [cite: 586-591]
					targetStream, exists := activeSessions[frame.Target]
					if exists {
						targetStream.Write(buf[:n]) // Forward zaszyfrowanej ramki
					} else {
						resp := protocol.UniversalFrame{Type: "ERROR", Code: "ERR_15", Desc: "Target Offline"}
						s.Write(resp.ToJSON())
					}
				}
			}
		}(stream)
	}
}
