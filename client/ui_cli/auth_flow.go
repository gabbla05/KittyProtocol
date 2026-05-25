package ui_cli

import (
	"errors"
	"strings"
	"time"

	"github.com/gabbla05/KittyProtocol/client/api"
)

var ErrQuitRequested = errors.New("quit requested")

// Async AUTH/REGISTER flow
func (ui *CliUI) RunAuthFlowAsync(client *api.KittyClient) (string, error) {
	for {
		ui.Println("Wybierz opcję:")
		ui.Println("  /login")
		ui.Println("  /register")
		ui.Println("  /quit")

		cmd := strings.TrimSpace(ui.ReadLine())

		switch cmd {

		case "/quit":
			client.Close()
			return "", ErrQuitRequested

		// ----------------------------------------------------
		// REGISTER (async)
		// ----------------------------------------------------
		case "/register":
			user, pass := ui.ReadCredentials()

			if err := client.SendRegister(user, pass); err != nil {
				ui.Println("[Client] REGISTER send error:", err)
				continue
			}

			select {
			case res := <-client.RegisterResult():
				if !res.OK {
					ui.Println("[Client] REGISTER failed:", res.Error())
					continue
				}
				ui.Println("[Client] REGISTER OK — możesz się teraz zalogować.")

			case <-time.After(5 * time.Second):
				ui.Println("[Client] REGISTER timeout")
				continue
			}

		// ----------------------------------------------------
		// LOGIN (async)
		// ----------------------------------------------------
		case "/login":
			user, pass := ui.ReadCredentials()

			if err := client.SendAuth(user, pass); err != nil {
				ui.Println("[Client] AUTH send error:", err)
				continue
			}

			select {
			case res := <-client.AuthResult():
				if !res.OK {
					ui.Println("[Client] AUTH failed:", res.Error())
					continue
				}
				ui.Println("[Client] AUTH OK — zalogowano.")
				return pass, nil

			case <-time.After(5 * time.Second):
				ui.Println("[Client] AUTH timeout")
				continue
			}

		default:
			ui.Println("Nieznana komenda.")
		}
	}
}
