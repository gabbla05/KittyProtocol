package ui_cli

import (
	"strings"

	"github.com/gabbla05/KittyProtocol/client/app"
)

// RunMainMenu drives the main interactive CLI loop.
// It blocks until the user logs out or the client disconnects.
func (ui *CliUI) RunMainMenu(a *app.App) {

	ui.printMenu()

	for {
		select {
		case <-a.Disconnected():
			ui.Println(ColorRed + "[Client] Disconnected from server. Exiting." + ColorReset)
			return
		default:
		}

		ui.Prompt()

		line := strings.TrimSpace(ui.ReadLine())
		if line == "" {
			continue
		}

		// handleCommand returns true when menu should exit (logout)
		if ui.handleCommand(line, a) {
			return
		}
	}
}
