package client

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app"
	"github.com/gabbla05/KittyProtocol/client/ui_cli"
)

func Start() {
	LoadEnv()
	
	client := api.NewKittyClient()
	ui := ui_cli.NewCliUI(client)

	// Logger CLI
	api.SetLogger(ui_cli.CliLogger{})

	disconnected := make(chan struct{})

	hubAddr := os.Getenv("KITTY_HUB_ADDR")
	if hubAddr == "" {
		hubAddr = "127.0.0.1:9999"
	}

	ui_cli.PrintBanner()
	ui.Println("\n[Client] Connecting to Hub:", hubAddr)

	// CONNECT
	if err := client.Connect(hubAddr); err != nil {
		ui.Println("[Client] Connection error:", err)
		return
	}

	// ----------------------------------------------------
	// START RECEIVER LOOP BEFORE WAITING FOR HELLO
	// ----------------------------------------------------
	client.StartReceiverLoop(disconnected)

	// ----------------------------------------------------
	// ASYNC HELLO
	// ----------------------------------------------------
	select {
	case res := <-client.HelloResult():
		if !res.OK {
			ui.Println("[Client] HELLO failed:", res.Error())
			client.Close()
			return
		}
	case <-time.After(5 * time.Second):
		ui.Println("[Client] HELLO timeout")
		client.Close()
		return
	}

	// ----------------------------------------------------
	// AUTH FLOW (CLI-specific)
	// ----------------------------------------------------
	pass, err := ui.RunAuthFlowAsync(client)
	if err == ui_cli.ErrQuitRequested {
		client.Close()
		return
	}
	if err != nil {
		ui.Println("[Client] AUTH error:", err)
		client.Close()
		return
	}

	// ----------------------------------------------------
	// AUTH SUCCESS → start application
	// ----------------------------------------------------
	application := app.NewApp(client, ui, disconnected)

	// ⬇️ NOWOŚĆ: przekazujemy masterKey = hasło użytkownika
	application.InitSecretStoreForUser(client.User(), []byte(pass))

	client.RegisterAckHandler(ui)

	setupSignalHandler(client)

	// Ping loop dopiero po AUTH
	client.StartPingLoop()

	ui.RunMainMenu(application)

	client.Close()
}

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
