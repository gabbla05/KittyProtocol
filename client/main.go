// main.go
// Entry point for the CLI version of KittyClient.
// This file wires together the API, UI and App layers.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app"
	"github.com/gabbla05/KittyProtocol/client/ui_cli"
)

func main() {
	client := api.NewKittyClient()
	ui := ui_cli.NewCliUI(client)

	// Shared disconnection channel for App and background loops.
	disconnected := make(chan struct{})

	application := app.NewApp(client, ui, disconnected)

	// Register ACK event handler (UI implements AckEventHandler).
	client.RegisterAckHandler(ui)

	// OS signal handling (Ctrl+C, SIGTERM, SIGQUIT).
	setupSignalHandler(client)

	// Resolve Hub address.
	hubAddr := os.Getenv("KITTY_HUB_ADDR")
	if hubAddr == "" {
		hubAddr = "127.0.0.1:9999"
	}

	fmt.Println("[Client] Connecting to Hub:", hubAddr)

	// QUIC connection.
	if err := client.Connect(hubAddr); err != nil {
		fmt.Println("[Client] Connection error:", err)
		return
	}

	// HELLO handshake.
	if err := client.WaitForHelloOK(); err != nil {
		fmt.Println("[Client] HELLO failed:", err)
		client.Close()
		return
	}

	// AUTH.
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

	// Background loops.
	client.StartReceiverLoop(disconnected)
	client.StartPingLoop()

	// Main workflow.
	application.RunMainMenu()

	// Cleanup.
	client.Close()
}

// setupSignalHandler installs a handler for OS termination signals.
// On signal, the client sends BYE, closes the session and exits.
func setupSignalHandler(client *api.KittyClient) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigCh
		_ = client.SendBye()
		client.Close()
		fmt.Println("\n[Client] Session closed due to signal.")
		os.Exit(0)
	}()
}
