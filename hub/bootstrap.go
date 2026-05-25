// bootstrap.go
// Contains environment loading and graceful shutdown logic for the Hub.
// This file has no protocol logic — only process-level lifecycle management.

package hub

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/quic-go/quic-go"
)

// loadEnv loads environment variables from a .env file if present.
// Missing .env is not treated as an error.
func loadEnv() {
	_ = godotenv.Load()

	// Disable colors if KITTY_LOG_COLOR=0
	if os.Getenv("KITTY_LOG_COLOR") == "0" {
		colorsEnabled = false
	}
}

// setupSignalHandler installs OS signal handlers and returns a context that is
// cancelled when the server should shut down.
//
// When SIGINT, SIGTERM, or SIGQUIT is received:
//   - all active sessions are terminated via globalSessions.Stop()
//   - the QUIC listener is closed (causing Accept() to unblock)
//   - the returned context is cancelled, allowing the accept loop to exit
//
// This enables a fully graceful shutdown of the Hub server.
func setupSignalHandler(listener *quic.Listener) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigCh
		logWarn("Caught signal: %v — initiating graceful shutdown", sig)

		// Stop all active sessions
		globalSessions.Stop()

		// Close QUIC listener (unblocks Accept)
		if err := listener.Close(); err != nil {
			logError("Error closing listener: %v", err)
		}

		// Cancel context to stop Accept loop
		cancel()

		logInfo("Graceful shutdown complete.")
	}()

	return ctx
}
