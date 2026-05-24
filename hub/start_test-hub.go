package hub

import (
	"context"
	"os"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/quic-go/quic-go"
)

// StartTestHub starts an isolated Hub instance for tests and returns:
//   - QUIC address to dial
//   - stop() function that shuts down listener and accept loop
//
// It intentionally does NOT reassign globalSessions to avoid data races.
func StartTestHub() (addr string, stop func(), err error) {
	// Ensure certs directory exists
	if err := os.MkdirAll("../certs", 0755); err != nil {
		return "", nil, err
	}

	tlsConf, err := certmanager.SetupTLSConfig("../certs/cert.pem", "../certs/key.pem")
	if err != nil {
		return "", nil, err
	}

	// Ensure globalSessions exists, but do NOT reassign it later.
	if globalSessions == nil {
		globalSessions = protection.NewSessionManager()
	}
	// Clean per‑test users only (safe, SessionManager is synchronized).
	globalSessions.Remove("alice")
	globalSessions.Remove("bob")

	// Fresh mock auth backend for this test run.
	globalAuth = auth.NewMockAuth()

	// Start listener on random port.
	listener, err := quic.ListenAddr("127.0.0.1:0", tlsConf, &quic.Config{
		KeepAlivePeriod: 1 * time.Second,
	})
	if err != nil {
		return "", nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Accept loop.
	go func() {
		for {
			conn, err := listener.Accept(ctx)
			if err != nil {
				return
			}
			go handleClient(conn)
		}
	}()

	stop = func() {
		// Caution: Do not touch globalSessions here.
		cancel()
		_ = listener.Close()
	}

	return listener.Addr().String(), stop, nil
}
