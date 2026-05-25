package client

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app"
	"github.com/gabbla05/KittyProtocol/client/ui_cli"
)

// Start is the main entrypoint for the KittyProtocol CLI client.
// It performs connection, HELLO handshake, AUTH flow, and launches the app.
func Start() {
	client := api.NewKittyClient()
	ui := ui_cli.NewCliUI(client)

	disconnected := make(chan struct{})

	hubAddr := os.Getenv("KITTY_HUB_ADDR")
	if hubAddr == "" {
		hubAddr = "127.0.0.1:9999"
	}

	ui.Println("[Client] Connecting to Hub:", hubAddr)

	if err := client.Connect(hubAddr); err != nil {
		ui.Println("[Client] Connection error:", err)
		return
	}

	if err := client.WaitForHelloOK(); err != nil {
		ui.Println("[Client] HELLO failed:", err)
		client.Close()
		return
	}

	// -------------------------------
	// AUTH FLOW (CLI-specific)
	// -------------------------------
	if err := ui.RunAuthFlow(client); err == ui_cli.ErrQuitRequested {
		client.Close()
		return
	}

	// -------------------------------
	// AUTH SUCCESS → start application
	// -------------------------------

	application := app.NewApp(client, ui, disconnected)

	application.InitSecretStoreForUser(client.User())

	client.RegisterAckHandler(ui)
	client.RegisterAppPayloadHandler(application.HandleIncomingPayload)

	setupSignalHandler(client)

	client.StartReceiverLoop(disconnected)
	client.StartPingLoop()

	ui.RunMainMenu(application)

	client.Close()
}

// setupSignalHandler installs OS signal handlers for graceful shutdown.
func setupSignalHandler(client *api.KittyClient) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigCh
		_ = client.SendBye()
		client.Close()
		os.Exit(0)
	}()
}
