package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/quic-go/quic-go"
)

var (
	globalSessions                   = protection.NewSessionManager()
	globalAuth     auth.AuthProvider = auth.NewMockAuth()
)

func main() {
	loadEnv()

	tlsConf, err := certmanager.SetupTLSConfig("certs/cert.pem", "certs/key.pem")
	if err != nil {
		logError("Failed to load TLS certificates: %v", err)
		return
	}

	dsn := "postgres://kitty:kittypass@localhost:5432/kittyhub?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logError("DB connection failed: %v", err)
		return
	}

	globalAuth = auth.NewDBAuth(db)

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
		logError("Failed to start listener: %v", err)
		return
	}

	logInfo("🐈 KittyProtocol Hub listening on %s", addr)

	setupSignalHandler(listener)

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			logError("Accept error: %v", err)
			return
		}

		go handleClient(conn)
	}
}

func loadEnv() {
	_ = godotenv.Load()
}

func setupSignalHandler(listener *quic.Listener) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigCh
		logWarn("Caught signal: %v", sig)

		globalSessions.Stop()
		_ = listener.Close()

		logInfo("Graceful shutdown complete.")
	}()
}
