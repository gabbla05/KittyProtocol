// hub/main.go
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

var globalSessions = protection.NewSessionManager()
var globalAuth auth.AuthProvider = auth.NewMockAuth()

func main() {
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

	interceptAddr := os.Getenv("KITTY_INTERCEPT_ADDR")
	if interceptAddr == "" {
		interceptAddr = "0.0.0.0:9999"
	}

	listener, err := quic.ListenAddr(interceptAddr, tlsConf, quicConf)
	if err != nil {
		panic(err)
	}

	fmt.Println("[Hub] 🐈 KittyProtocol Hub listening on", interceptAddr)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigCh
		fmt.Println("\n[Hub] Caught signal:", sig)
		listener.Close()
	}()

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Println("[Hub] Accept error:", err)
			return
		}
		go handleClient(conn)
	}
}
