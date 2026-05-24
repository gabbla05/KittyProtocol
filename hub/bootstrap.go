// bootstrap.go
// Contains environment loading and graceful shutdown logic for the Hub.
// This file has no protocol logic — only process-level lifecycle management.

package hub

import (
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

// setupSignalHandler installs OS signal handlers for graceful shutdown.
// On SIGINT/SIGTERM, all sessions are terminated and the QUIC listener is closed.
func setupSignalHandler(listener *quic.Listener) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigCh
		logWarn("Caught signal: %v", sig)

		// Stop all active sessions
		globalSessions.Stop()

		// Close QUIC listener
		_ = listener.Close()

		logInfo("Graceful shutdown complete.")
	}()
}
