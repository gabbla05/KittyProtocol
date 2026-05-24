// start.go
// Entry point for the Hub server. Initializes TLS, QUIC, authentication backend,
// session manager, and begins accepting client connections. This function blocks
// until the server is shut down.

package hub

import (
	"context"
	"database/sql"
	"os"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	_ "github.com/lib/pq"
	"github.com/quic-go/quic-go"
)

// globalSessions manages all active Hub sessions.
// It is initialized once at startup and shared across handlers.
var globalSessions = protection.NewSessionManager()

// globalAuth provides authentication backend for the Hub.
// It is replaced during startup depending on configuration (mock or DB).
var globalAuth auth.AuthProvider

// Start initializes the KittyProtocol Hub server, configures TLS, QUIC,
// authentication backend, and begins accepting incoming client connections.
func Start() {
	loadEnv()

	// Load TLS certificates for QUIC transport security.
	tlsConf, err := certmanager.SetupTLSConfig("certs/cert.pem", "certs/key.pem")
	if err != nil {
		logError("Failed to load TLS certificates: %v", err)
		return
	}

	// Initialize authentication backend (PostgreSQL).
	dsn := os.Getenv("KITTY_DB_DSN")
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logError("DB connection failed: %v", err)
		return
	}
	defer db.Close()

	globalAuth = auth.NewDBAuth(db)
	logInfo("Using authentication backend with DSN: %s", dsn)

	// Configure QUIC transport parameters.
	quicConf := &quic.Config{
		MaxIdleTimeout:          quicMaxIdleTimeout,
		KeepAlivePeriod:         quicKeepAlivePeriod,
		Allow0RTT:               quicAllow0RTT,
		DisablePathMTUDiscovery: quicDisablePMTU,
	}

	// Determine listening address.
	addr := os.Getenv("KITTY_INTERCEPT_ADDR")
	if addr == "" {
		addr = defaultHubAddress
	}

	// Start QUIC listener.
	listener, err := quic.ListenAddr(addr, tlsConf, quicConf)
	if err != nil {
		logError("Failed to start listener: %v", err)
		return
	}

	logInfo("🐈 KittyProtocol Hub listening on %s", addr)

	// Handle SIGINT/SIGTERM for graceful shutdown.
	setupSignalHandler(listener)

	// Accept incoming QUIC connections.
	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			logError("Accept error: %v", err)
			return
		}

		// Each client is handled in its own goroutine.
		go handleClient(conn)
	}
}
