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
		// Spawn a new Goroutine for each connected client
		go handleConnection(conn)
	}
}

func handleConnection(conn *quic.Conn) { // Note the *quic.Conn pointer!
	stream, err := (*conn).AcceptStream(context.Background())
	if err != nil {
		return
	}
	defer stream.Close()

	buf := make([]byte, 2048)
	for {
		n, err := stream.Read(buf)
		if err != nil {
			return // Client disconnected
		}

		frame, err := protocol.ParseFrame(buf[:n])
		if err != nil {
			errFrame := protocol.UniversalFrame{Type: "ERROR", Code: "ERR_02", Desc: err.Error()} // [cite: 677]
			stream.Write(errFrame.ToJSON())
			continue
		}

		fmt.Printf("[Hub] Received from %s: %s\n", (*conn).RemoteAddr(), frame.Type)

		switch frame.Type {
		case "HELLO":
			resp := protocol.UniversalFrame{Type: "MEOW_OK", MsgID: frame.MsgID, Status: "Ready for auth"}
			stream.Write(resp.ToJSON())

		case "AUTH":
			token := "tok_" + frame.User + "_123"
			onlineUsers[frame.User] = UserInfo{
				IP:    "127.0.0.1",
				Port:  8888,
				Token: token,
			}
			resp := protocol.UniversalFrame{Type: "MEOW_OK", MsgID: frame.MsgID, Token: token}
			stream.Write(resp.ToJSON())
			fmt.Printf("[Hub] User '%s' successfully authenticated!\n", frame.User)
		}
	}
}
