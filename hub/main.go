package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

var sm = protection.NewSessionManager()

func main() {
	cert, _ := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	tlsConf := &tls.Config{Certificates: []tls.Certificate{cert}, NextProtos: []string{"kitty-quic"}}
	listener, _ := quic.ListenAddr("127.0.0.1:9999", tlsConf, nil)

	fmt.Println("🐈 KittyProtocol Hub (Task 9 Active) na 127.0.0.1:9999")

	for {
		conn, _ := listener.Accept(context.Background())
		go handleClient(conn)
	}
}

func handleClient(conn *quic.Conn) {
	stream, _ := conn.AcceptStream(context.Background())
	defer stream.Close()

	// Task 9: Auth Timeout (20s) [cite: 662]
	authTimer := protection.StartAuthTimer(func() {
		conn.CloseWithError(0x03, "ERR_03: Auth Timeout")
	})

	var currentSess *protection.Session
	buf := make([]byte, 4096)

	for {
		n, err := stream.Read(buf)
		if err != nil {
			break
		}

		frame, _ := protocol.ParseFrame(buf[:n])
		if frame == nil {
			continue
		}

		if currentSess != nil {
			currentSess.LastActive = time.Now()
		}

		switch frame.Type {
		case "AUTH":
			if auth.CheckCredentials(frame.User, frame.Pass) { // Task 6 [cite: 659]
				authTimer.Stop()
				currentSess = &protection.Session{
					ID:         frame.User,
					LastActive: time.Now(),
					Limiter:    protection.NewRateLimiter(10), // Task 9 [cite: 663]
					CloseFunc:  func() { conn.CloseWithError(0x09, "Idle Timeout") },
				}
				sm.Add(frame.User, currentSess)
				stream.Write((&protocol.UniversalFrame{Type: "MEOW_OK", Status: "Logged in"}).ToJSON())
			}
		case "DATA":
			// Task 9: Rate Limit Check (10 msg/s) [cite: 207, 415]
			if currentSess == nil || !currentSess.Limiter.Allow() {
				stream.Write((&protocol.UniversalFrame{Type: "ERROR", Code: "ERR_07"}).ToJSON())
				continue
			}
			// Dalej routing MB...
		}
	}
}
