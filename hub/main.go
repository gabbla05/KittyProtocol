package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// UserInfo stores logged-in users' data in memory (stateless approach).
type UserInfo struct {
	IP    string
	Port  int
	Token string
}

var onlineUsers = make(map[string]UserInfo)

func main() {
	// Load TLS certificates required by QUIC
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
	fmt.Println("🐈 KittyProtocol Hub (QUIC) is listening on 127.0.0.1:9999...")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			continue
		}

		// Zaczynamy nową goroutine bez argumentów (Go 1.22+ bezpiecznie przechwytuje conn)
		go func() {
			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				return
			}
			defer stream.Close()

			buf := make([]byte, 2048)
			for {
				n, err := stream.Read(buf)
				if err != nil {
					return // Klient się rozłączył
				}

				frame, err := protocol.ParseFrame(buf[:n])
				if err != nil {
					errFrame := protocol.UniversalFrame{Type: "ERROR", Code: "ERR_02", Desc: err.Error()}
					stream.Write(errFrame.ToJSON())
					continue
				}

				fmt.Printf("[Hub] Received from %s: %s\n", conn.RemoteAddr(), frame.Type)

				switch frame.Type {
				case "HELLO":
					resp := protocol.UniversalFrame{Type: "MEOW_OK", MsgID: frame.MsgID, Status: "Ready for auth"}
					stream.Write(resp.ToJSON())

				case "AUTH":
					token := "tok_" + frame.User + "_123"

					// Dynamiczne przypisanie portu P2P zadeklarowanego przez klienta
					clientPort := frame.Port
					if clientPort == 0 {
						clientPort = 8888 // Fallback
					}

					onlineUsers[frame.User] = UserInfo{
						IP:    "127.0.0.1",
						Port:  clientPort,
						Token: token,
					}
					resp := protocol.UniversalFrame{Type: "MEOW_OK", MsgID: frame.MsgID, Token: token}
					stream.Write(resp.ToJSON())
					fmt.Printf("[Hub] User '%s' successfully authenticated on port %d!\n", frame.User, clientPort)

				case "LOOKUP":
					info, exists := onlineUsers[frame.Target]
					if exists {
						resp := protocol.UniversalFrame{Type: "MEOW_OK", MsgID: frame.MsgID, IP: info.IP, Port: info.Port}
						stream.Write(resp.ToJSON())
					} else {
						resp := protocol.UniversalFrame{Type: "ERROR", MsgID: frame.MsgID, Code: "ERR_09", Desc: "Recipient Offline"}
						stream.Write(resp.ToJSON())
					}
				}
			}
		}()
	}
}
