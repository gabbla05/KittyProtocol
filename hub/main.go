package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

// globalSessions holds all active sessions on the Hub.
var globalSessions = protection.NewSessionManager()

// main starts the KittyProtocol Hub and listens for incoming QUIC connections.
func main() {

	// TLS 1.3 + ALPN kitty-quic-v1
	// ALPN = Application-Layer Protocol Negotiation
	// It is TLS mechanism that enables client and server to settle:
	// „Which application protocol will work inside TLS tunnel?”

	tlsConf, err := certmanager.SetupTLSConfig("certs/cert.pem", "certs/key.pem")
	if err != nil {
		panic(err)
	}

	// QUIC config consistent with documentation
	quicConf := &quic.Config{
		MaxIdleTimeout:          60 * time.Second,
		KeepAlivePeriod:         30 * time.Second,
		Allow0RTT:               true,
		DisablePathMTUDiscovery: false,
	}

	listener, err := quic.ListenAddr("127.0.0.1:9999", tlsConf, quicConf)
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
