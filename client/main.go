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

	// jeden wspólny kanał disconnected
	disconnected := make(chan struct{})

	application := app.NewApp(client, ui, disconnected)

	// ACK handler
	client.RegisterAckHandler(ui)

	// sygnały OS
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigCh
		_ = client.SendBye()
		client.Close()
		fmt.Println("\n[Client] Session closed due to signal.")
		os.Exit(0)
	}()

	// Connect
	hubAddr := os.Getenv("KITTY_HUB_ADDR")
	if hubAddr == "" {
		hubAddr = "127.0.0.1:9999"
	}

	fmt.Println("[Client] Connecting to Hub:", hubAddr)

	if err := client.Connect(hubAddr); err != nil {
		fmt.Println("[Client] Connection error:", err)
		return
	}

	// HELLO
	if err := client.WaitForHelloOK(); err != nil {
		fmt.Println("[Client] HELLO failed:", err)
		client.Close()
		return
	}

	// AUTH
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

	// background loops
	client.StartReceiverLoop(disconnected)
	client.StartPingLoop()

	// workflow
	application.RunMainMenu()

	// cleanup
	client.Close()
}
