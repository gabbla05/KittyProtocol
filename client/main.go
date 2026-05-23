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

	// Resolve Hub address.
	hubAddr := os.Getenv("KITTY_HUB_ADDR")
	if hubAddr == "" {
		hubAddr = "127.0.0.1:9999"
	}

	fmt.Println("[Client] Connecting to Hub:", hubAddr)

	// 1. QUIC connection + stream + HELLO (Connect() does all of this)
	if err := client.Connect(hubAddr); err != nil {
		fmt.Println("[Client] Connection error:", err)
		return
	}

	// 2. Wait for MEOW_OK after HELLO
	if err := client.WaitForHelloOK(); err != nil {
		fmt.Println("[Client] HELLO failed:", err)
		client.Close()
		return
	}

	// 3. AUTH
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

	// 4. Create App layer
	application := app.NewApp(client, ui, disconnected)

	// 5. Initialize per-user SecretStore
	application.InitSecretStoreForUser(client.User())

	// 6. Register handlers
	client.RegisterAckHandler(ui)
	client.RegisterAppPayloadHandler(application.HandleIncomingPayload)

	// 7. OS signal handling
	setupSignalHandler(client)

	// 8. Start background loops
	client.StartReceiverLoop(disconnected)
	client.StartPingLoop()

	// 9. Main workflow
	application.RunMainMenu()

	// 10. Cleanup
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
