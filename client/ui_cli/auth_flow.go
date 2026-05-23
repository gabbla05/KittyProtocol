package ui_cli

import (
	"errors"
	"strings"

	"github.com/gabbla05/KittyProtocol/client/api"
)

var ErrQuitRequested = errors.New("quit requested")

// RunAuthFlow handles /login and /register before entering main menu.
// This is CLI-specific and will be replaced by GUI in the future.
func (ui *CliUI) RunAuthFlow(client *api.KittyClient) error {
	for {
		ui.Println("Wybierz opcję:")
		ui.Println("  /login")
		ui.Println("  /register")
		ui.Println("  /quit")

		cmd := strings.TrimSpace(ui.ReadLine())

		switch cmd {

		case "/quit":
			client.Close()
			return ErrQuitRequested

		case "/register":
			user, pass := ui.ReadCredentials()

			if err := client.SendRegister(user, pass); err != nil {
				ui.Println("[Client] REGISTER send error:", err)
				continue
			}

			if err := client.WaitForRegisterOK(); err != nil {
				ui.Println("[Client] REGISTER failed:", err)
				continue
			}

			ui.Println("[Client] REGISTER OK — możesz się teraz zalogować.")

		case "/login":
			user, pass := ui.ReadCredentials()

			if err := client.SendAuth(user, pass); err != nil {
				ui.Println("[Client] AUTH send error:", err)
				continue
			}

			if err := client.WaitForAuthOK(); err != nil {
				ui.Println("[Client] AUTH failed:", err)
				continue
			}

			ui.Println("[Client] AUTH OK — zalogowano.")
			return nil

		default:
			ui.Println("Nieznana komenda.")
		}
	}
}
