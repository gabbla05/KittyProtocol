// hub/main.go
// Entry point for the KittyProtocol Hub.
// Responsible for:
//   - loading configuration
//   - initializing TLS + QUIC
//   - accepting incoming connections
//   - dispatching each connection to handleClient()

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/joho/godotenv"
	"github.com/quic-go/quic-go"
)

var (
	globalSessions                   = protection.NewSessionManager()
	globalAuth     auth.AuthProvider = auth.NewMockAuth()
)

func main() {
	loadEnv()
	// if sth goes wrong with reading env please try using _ = godotenv.Load() and

	tlsConf, err := certmanager.SetupTLSConfig("certs/cert.pem", "certs/key.pem")
	if err != nil {
		fmt.Println("[Hub] Failed to load TLS certificates:", err)
		return
	}

	quicConf := &quic.Config{
		MaxIdleTimeout:          60 * time.Second,
		KeepAlivePeriod:         30 * time.Second,
		Allow0RTT:               true,
		DisablePathMTUDiscovery: false,
	}

	addr := os.Getenv("KITTY_INTERCEPT_ADDR")
	if addr == "" {
		addr = "0.0.0.0:9999"
	}

	listener, err := quic.ListenAddr(addr, tlsConf, quicConf)
	if err != nil {
		fmt.Println("[Hub] Failed to start listener:", err)
		return
	}

	fmt.Println("[Hub] 🐈 KittyProtocol Hub listening on", addr)

	setupSignalHandler(listener)

	// Accept loop
	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Println("[Hub] Accept error:", err)
			return
		}

		go handleClient(conn)
	}
}

// loadEnv loads environment variables from .env if present.
func loadEnv() {
	_ = godotenv.Load()
}

// setupSignalHandler gracefully shuts down the listener on OS signals.
func setupSignalHandler(listener *quic.Listener) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigCh
		fmt.Println("\n[Hub] Caught signal:", sig)
		_ = listener.Close()
	}()
}
