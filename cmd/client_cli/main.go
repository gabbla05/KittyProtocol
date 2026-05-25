package main

import (
	"github.com/gabbla05/KittyProtocol/client"
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/ui_cli"
)

// Entry point for the CLI client application.
func main() {
	api.SetLogger(ui_cli.CliLogger{})
	client.Start()
}
