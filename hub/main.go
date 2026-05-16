// hub/main.go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/joho/godotenv"
	"github.com/quic-go/quic-go"
)

// globalSessions holds all active sessions on the Hub.
var globalSessions = protection.NewSessionManager()

// main starts the KittyProtocol Hub and listens for incoming QUIC connections.
func main() {
	// Load environment variables from .env file (if present)
	_ = godotenv.Load()

	tlsConf, err := certmanager.SetupTLSConfig("certs/cert.pem", "certs/key.pem")
	if err != nil {
		panic(err)
	}

	quicConf := &quic.Config{
		MaxIdleTimeout:          60 * time.Second,
		KeepAlivePeriod:         30 * time.Second,
		Allow0RTT:               true,
		DisablePathMTUDiscovery: false,
	}

	// !! Read intercept address from env (allows flexibility in different environments). !!
	interceptAddr := os.Getenv("KITTY_INTERCEPT_ADDR")
	if interceptAddr == "" {
		interceptAddr = "0.0.0.0:9999" // sensowny default dla dev
	}

	listener, err := quic.ListenAddr(interceptAddr, tlsConf, quicConf)
	if err != nil {
		panic(err)
	}

	fmt.Println("🐈 KittyProtocol Hub listening on", interceptAddr)

	// Signal handling (SIGINT, SIGTERM, SIGQUIT) – graceful listener shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigCh
		fmt.Println("\n[Hub] Caught signal:", sig)
		// Close listener – Accept will start returning errors.
		listener.Close()
	}()

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Println("Accept error:", err)
			// After listener is closed by signal – exit main.
			return
		}
		go handleClient(conn)
	}
}
