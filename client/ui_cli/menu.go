package ui_cli

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app"
)

func (ui *CliUI) RunMainMenu(a *app.App) {
	for {
		select {
		case <-a.Disconnected():
			ui.Println(ColorRed + "[Client] Rozłączono z serwerem. Zamykanie aplikacji." + ColorReset)
			return
		default:
		}

		ui.printMenu()
		ui.Prompt()

		line := strings.TrimSpace(ui.ReadLine())
		if line == "" {
			continue
		}

		switch {

		// ----------------------------------------------------
		// LOGOUT
		// ----------------------------------------------------
		case line == "/logout":
			if active, peer := a.ChatState().IsActive(); active && peer != "" {
				if err := a.EndChat("user logout client"); err != nil {
					ui.Println(ColorRed+"[CHAT ERROR]"+ColorReset, err)
				}
			}
			_ = a.Client().SendBye()
			return

		// ----------------------------------------------------
		// STATUS
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/status "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/status "))
			user = strings.ToLower(user)
			if user == "" {
				ui.Println(ColorYellow + "Usage: /status <user>" + ColorReset)
				continue
			}
			_ = a.Client().SendGetStatus(user)

		// ----------------------------------------------------
		// SECRET
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/secret "):
			args := strings.Fields(line)
			if len(args) < 2 {
				ui.Println(ColorYellow + "Usage: /secret <user> [file:<path>]" + ColorReset)
				continue
			}

			user := strings.ToLower(args[1])

			var secret []byte
			if len(args) == 3 && strings.HasPrefix(args[2], "file:") {
				path := strings.TrimPrefix(args[2], "file:")
				data, err := os.ReadFile(path)
				if err != nil {
					ui.Printf(ColorRed+"[E2EE] Failed to read secret file: %v\n"+ColorReset, err)
					continue
				}
				secret = bytes.TrimSpace(data)
			} else {
				secret = ui.ReadSharedSecret()
			}

			if err := a.Client().SetSharedSecretForPeer(user, secret); err != nil {
				ui.Println(ColorRed+"[E2EE] Error deriving keys:"+ColorReset, err)
				continue
			}
			if err := a.Secrets().Set(user, secret); err != nil {
				ui.Println(ColorRed+"[E2EE] Error saving secret:"+ColorReset, err)
				continue
			}

			ui.Printf(ColorGreen+"[E2EE] Shared secret configured for %s.\n"+ColorReset, user)

		// ----------------------------------------------------
		// CHAT REQUEST
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/chat "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/chat "))
			user = strings.ToLower(user)
			if user == "" {
				ui.Println(ColorYellow + "Usage: /chat <user>" + ColorReset)
				continue
			}

			if err := a.StartChatRequest(user); err != nil {
				if errors.Is(err, api.ErrNoSharedSecret) || strings.Contains(err.Error(), "no shared secret") {
					ui.Printf(ColorBlue+"[CHAT] Brak wspólnego sekretu z %s. Użyj /secret %s.\n"+ColorReset, user, user)
				} else {
					ui.Println(ColorRed+"[CHAT ERROR]"+ColorReset, err)
				}
			} else {
				ui.Printf(ColorBlue+"[CHAT] Wysłano CHAT_REQUEST do %s.\n"+ColorReset, user)
			}

		// ----------------------------------------------------
		// ACCEPT CHAT
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/accept "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/accept "))
			user = strings.ToLower(user)
			if user == "" {
				ui.Println(ColorYellow + "Usage: /accept <user>" + ColorReset)
				continue
			}

			if err := a.AcceptChat(user); err != nil {
				ui.Println(ColorRed+"[CHAT ERROR]"+ColorReset, err)
			}

		// ----------------------------------------------------
		// REFUSE CHAT
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/refuse "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/refuse "))
			user = strings.ToLower(user)
			if user == "" {
				ui.Println(ColorYellow + "Usage: /refuse <user>" + ColorReset)
				continue
			}

			if err := a.RefuseChat(user, "user refused"); err != nil {
				ui.Println(ColorRed+"[CHAT ERROR]"+ColorReset, err)
			}

		// ----------------------------------------------------
		// SEND MESSAGE
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/msg "):
			text := strings.TrimSpace(strings.TrimPrefix(line, "/msg "))
			if text == "" {
				continue
			}

			if err := a.SendTextMessage(text); err != nil {
				ui.Println(ColorRed+"[CHAT ERROR]"+ColorReset, err)
			}

		// ----------------------------------------------------
		// END CHAT
		// ----------------------------------------------------
		case line == "/end":
			if err := a.EndChat("user ended chat"); err != nil {
				ui.Println(ColorRed+"[CHAT ERROR]"+ColorReset, err)
			}

		default:
			ui.Println(ColorYellow + "Nieznana komenda." + ColorReset)
		}
	}
}

func (ui *CliUI) printMenu() {
	ui.Println(ColorBlue + "\n  ======================" + ColorReset)
	ui.Println(ColorBlue + " | " + ColorGreen + "Dostępne komendy:    " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /status <user>    " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /secret <user>    " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /chat <user>      " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /accept <user>    " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /refuse <user>    " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /msg <tekst>      " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /end              " + ColorBlue + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /logout           " + ColorBlue + "|")
	ui.Println(ColorBlue + "  ======================\n" + ColorReset)
}
