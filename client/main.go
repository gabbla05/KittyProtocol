package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabbla05/KittyProtocol/client/api"
)

// Main entry point for the CLI version of the KittyProtocol client.
// All protocol logic is inside KittyClient.
// This file only orchestrates the flow.
func main() {
	client := api.NewKittyClient()
	ui := NewCliUI(client)

	// Register UI as ACK event handler
	client.RegisterAckHandler(ui)

	// Handle OS signals (Ctrl+C, SIGTERM)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigCh
		_ = client.SendBye()
		client.Close()
		fmt.Println("\n[Client] Session closed due to signal.")
		os.Exit(0)
	}()

	// Connect to Hub
	hubAddr := os.Getenv("KITTY_HUB_ADDR")
	if hubAddr == "" {
		hubAddr = "127.0.0.1:9999"
	}

	fmt.Println("[Client] Connecting to Hub:", hubAddr)

	if err := client.Connect(hubAddr); err != nil {
		fmt.Println("[Client] Connection error:", err)
		return
	}

	// HELLO handshake
	if err := client.WaitForHelloOK(); err != nil {
		fmt.Println("[Client] HELLO failed:", err)
		client.Close()
		return
	}

	// AUTH phase
	user, pass := ui.ReadCredentials()
	if err := client.SendAuth(user, pass); err != nil {
		fmt.Println("[Client] AUTH send error:", err)
		client.Close()
		return
	}

	if err := client.WaitForAuthOK(); err != nil {
		fmt.Println("[Client] AUTH failed:", err)
		client.Close()
		return
	}

	// Target selection
	target := ui.ReadTarget()
	client.SetTarget(target)

	secret := ui.ReadSharedSecret()
	if err := client.SetSharedSecret(secret); err != nil {
		fmt.Println("[Client] Failed to set shared secret:", err)
		client.Close()
		return
	}

	fmt.Println("[Client] Session established with target:", target)

	// Start background loops
	disconnected := make(chan struct{})
	client.StartReceiverLoop(disconnected)
	client.StartPingLoop()

	// Main send loop (CLI)
	ui.RunSendLoop(disconnected)

	// Cleanup
	client.Close()
	time.Sleep(200 * time.Millisecond)
}
