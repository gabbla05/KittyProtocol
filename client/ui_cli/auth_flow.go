package ui_cli

import (
	"errors"
	"strings"
	"time"

	"github.com/gabbla05/KittyProtocol/client/api"
)

// ErrQuitRequested is returned when the user chooses /quit during auth flow.
var ErrQuitRequested = errors.New("quit requested")

// authTimeout defines how long the UI waits for AUTH/REGISTER results.
const authTimeout = 5 * time.Second

// RunAuthFlowAsync drives the interactive LOGIN/REGISTER flow.
// It is UI-only logic: no protocol state is mutated here.
// Returns the user's password (for secret store) or ErrQuitRequested.
func (ui *CliUI) RunAuthFlowAsync(client *api.KittyClient) (string, error) {
	for {
		ui.printAuthMenu()
		ui.Prompt()

		cmd := strings.TrimSpace(ui.ReadLine())

		switch cmd {

		case "/quit":
			client.Close()
			return "", ErrQuitRequested

		case "/register":
			if err := ui.handleRegister(client); err != nil {
				ui.Println("[Client] REGISTER error:", err)
			}

		case "/login":
			pass, err := ui.handleLogin(client)
			if err != nil {
				ui.Println("[Client] AUTH error:", err)
				continue
			}
			return pass, nil

		default:
			ui.Println("Nieznana komenda.")
		}
	}
}

// printAuthMenu prints the main AUTH/REGISTER menu.
func (ui *CliUI) printAuthMenu() {
	ui.Println(ColorBlue + "\n  ==================" + ColorReset)
	ui.Println(ColorBlue + " | Wybierz opcję:   |")
	ui.Println(" |                  |")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + "   /login      " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + "   /register   " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + "   /quit       " + ColorBlue + "|")
	ui.Println("  ==================\n" + ColorReset)
}

// handleRegister performs the REGISTER flow.
func (ui *CliUI) handleRegister(client *api.KittyClient) error {
	user, pass := ui.ReadCredentials()

	if err := client.SendRegister(user, pass); err != nil {
		return err
	}

	select {
	case res := <-client.RegisterResult():
		if !res.OK {
			return res
		}
		ui.Println("[Client] REGISTER OK — możesz się teraz zalogować.")
		return nil

	case <-time.After(authTimeout):
		return errors.New("REGISTER timeout")
	}
}

// handleLogin performs the AUTH flow and returns the password on success.
func (ui *CliUI) handleLogin(client *api.KittyClient) (string, error) {
	user, pass := ui.ReadCredentials()

	if err := client.SendAuth(user, pass); err != nil {
		return "", err
	}

	select {
	case res := <-client.AuthResult():
		if !res.OK {
			return "", res
		}
		ui.Println("[Client] AUTH OK — zalogowano.")
		return pass, nil

	case <-time.After(authTimeout):
		return "", errors.New("AUTH timeout")
	}
}
