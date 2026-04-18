package main

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

// globalSessions holds all active sessions on the Hub.
var globalSessions = protection.NewSessionManager()

// main starts the KittyProtocol Hub and listens for incoming QUIC connections.
func main() {
    cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
    if err != nil {
        panic(err)
    }
    tlsConf := &tls.Config{
        Certificates: []tls.Certificate{cert},
        NextProtos:   []string{"kitty-quic"},
    }

    listener, err := quic.ListenAddr("127.0.0.1:9999", tlsConf, nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("🐈 KittyProtocol Hub listening on 127.0.0.1:9999")

    for {
        conn, err := listener.Accept(context.Background())
        if err != nil {
            fmt.Println("Accept error:", err)
            continue
        }
        go handleClient(conn)
    }
}
