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
			ui.Println("[Client] Rozłączono z serwerem. Zamykanie aplikacji.")
			return
		default:
		}

		ui.printMenu()
		line := strings.TrimSpace(ui.ReadLine())
		if line == "" {
			continue
		}

		switch {

		// ----------------------------------------------------
		// QUIT
		// ----------------------------------------------------
		case line == "/quit":
			if active, peer := a.ChatState().IsActive(); active && peer != "" {
				if err := a.EndChat("user quit client"); err != nil {
					ui.Println("[CHAT ERROR]", err)
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
				ui.Println("Usage: /status <user>")
				continue
			}
			_ = a.Client().SendGetStatus(user)

		// ----------------------------------------------------
		// SECRET (local only)
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/secret "):
			args := strings.Fields(line)
			if len(args) < 2 {
				ui.Println("Usage: /secret <user> [file:<path>]")
				continue
			}

			user := strings.ToLower(args[1])

			var secret []byte
			if len(args) == 3 && strings.HasPrefix(args[2], "file:") {
				path := strings.TrimPrefix(args[2], "file:")
				data, err := os.ReadFile(path)
				if err != nil {
					ui.Printf("[E2EE] Failed to read secret file: %v\n", err)
					continue
				}
				secret = bytes.TrimSpace(data)
			} else {
				secret = ui.ReadSharedSecret()
			}

			if err := a.Client().SetSharedSecretForPeer(user, secret); err != nil {
				ui.Println("[E2EE] Error deriving keys:", err)
				continue
			}
			if err := a.Secrets().Set(user, secret); err != nil {
				ui.Println("[E2EE] Error saving secret:", err)
				continue
			}

			ui.Printf("[E2EE] Shared secret configured for %s.\n", user)

		// ----------------------------------------------------
		// CHAT REQUEST (async)
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/chat "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/chat "))
			user = strings.ToLower(user)
			if user == "" {
				ui.Println("Usage: /chat <user>")
				continue
			}

			if err := a.StartChatRequest(user); err != nil {
				// Ładny komunikat dla braku sekretu
				if errors.Is(err, api.ErrNoSharedSecret) || strings.Contains(err.Error(), "no shared secret") {
					ui.Printf("[CHAT] Brak wspólnego sekretu z %s. Użyj /secret %s.\n", user, user)
				} else {
					ui.Println("[CHAT ERROR]", err)
				}
			} else {
				ui.Printf("[CHAT] Wysłano CHAT_REQUEST do %s.\n", user)
			}

		// ----------------------------------------------------
		// ACCEPT CHAT
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/accept "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/accept "))
			user = strings.ToLower(user)
			if user == "" {
				ui.Println("Usage: /accept <user>")
				continue
			}

			if err := a.AcceptChat(user); err != nil {
				ui.Println("[CHAT ERROR]", err)
			}

		// ----------------------------------------------------
		// REFUSE CHAT
		// ----------------------------------------------------
		case strings.HasPrefix(line, "/refuse "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/refuse "))
			user = strings.ToLower(user)
			if user == "" {
				ui.Println("Usage: /refuse <user>")
				continue
			}

			if err := a.RefuseChat(user, "user refused"); err != nil {
				ui.Println("[CHAT ERROR]", err)
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
				ui.Println("[CHAT ERROR]", err)
			}

		// ----------------------------------------------------
		// END CHAT
		// ----------------------------------------------------
		case line == "/end":
			if err := a.EndChat("user ended chat"); err != nil {
				ui.Println("[CHAT ERROR]", err)
			}

		default:
			ui.Println("Nieznana komenda.")
		}
	}
}

func (ui *CliUI) printMenu() {
	ui.Println("Dostępne komendy:")
	ui.Println("  /status <user>")
	ui.Println("  /secret <user>")
	ui.Println("  /chat <user>")
	ui.Println("  /accept <user>")
	ui.Println("  /refuse <user>")
	ui.Println("  /msg <tekst>")
	ui.Println("  /end")
	ui.Println("  /quit")
}
