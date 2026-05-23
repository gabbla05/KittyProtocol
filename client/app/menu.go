package app

import (
	"bytes"
	"os"
	"strings"
)

// RunMainMenu displays the main command loop.
func (a *App) RunMainMenu() {
	for {
		select {
		case <-a.disconnected:
			a.ui.Println("[Client] Rozłączono z serwerem. Zamykanie aplikacji.")
			return
		default:
		}

		a.printMenu()
		line := strings.TrimSpace(a.ui.ReadLine())
		if line == "" {
			continue
		}

		switch {

		// Exit
		case line == "/quit":
			_ = a.client.SendBye()
			return

		// Presence status
		case strings.HasPrefix(line, "/status "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/status "))
			if user == "" {
				a.ui.Println("Usage: /status <user>")
				continue
			}
			_ = a.client.SendGetStatus(user)

			// Configure shared secret for a peer (E2EE)
		case strings.HasPrefix(line, "/secret "):
			args := strings.Fields(line)
			if len(args) < 2 {
				a.ui.Println("Usage: /secret <user> [file:<path>]")
				continue
			}

			user := args[1]

			var secret []byte
			if len(args) == 3 && strings.HasPrefix(args[2], "file:") {
				path := strings.TrimPrefix(args[2], "file:")
				data, err := os.ReadFile(path)
				if err != nil {
					a.ui.Printf("[E2EE] Failed to read secret file: %v\n", err)
					continue
				}
				secret = bytes.TrimSpace(data)
			} else {
				secret = a.ui.ReadSharedSecret()
			}

			if err := a.client.SetSharedSecretForPeer(user, secret); err != nil {
				a.ui.Println("[E2EE] Error deriving keys:", err)
				continue
			}
			if err := a.secrets.Set(user, secret); err != nil {
				a.ui.Println("[E2EE] Error saving secret:", err)
				continue
			}

			a.ui.Printf("[E2EE] Shared secret configured for %s.\n", user)

		// Start chat
		case strings.HasPrefix(line, "/chat "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/chat "))
			if user == "" {
				a.ui.Println("Usage: /chat <user>")
				continue
			}

			if a.chatState.Active {
				a.ui.Println("[CHAT] Masz już aktywny czat. Użyj /end aby zakończyć.")
				continue
			}
			if a.chatState.PendingRequestFrom != "" {
				a.ui.Printf("[CHAT] Masz oczekujący request od %s. Użyj /accept lub /refuse.\n",
					a.chatState.PendingRequestFrom)
				continue
			}

			// Ensure E2EE secret for this peer.
			if secret, ok := a.secrets.Get(user); ok {
				if err := a.client.SetSharedSecretForPeer(user, secret); err != nil {
					a.ui.Println("[E2EE] Error deriving keys from stored secret:", err)
					continue
				}
				a.ui.Printf("[E2EE] Loaded stored shared secret for %s.\n", user)
			} else {
				secret := a.ui.ReadSharedSecret()
				if err := a.client.SetSharedSecretForPeer(user, secret); err != nil {
					a.ui.Println("[E2EE] Error deriving keys:", err)
					continue
				}
				if err := a.secrets.Set(user, secret); err != nil {
					a.ui.Println("[E2EE] Error saving secret:", err)
					continue
				}
			}

			if err := a.StartChatRequest(user); err != nil {
				a.ui.Println("Błąd:", err)
			} else {
				a.ui.Printf("[CHAT] Wysłano CHAT_REQUEST do %s.\n", user)
			}

		// Accept chat request
		case strings.HasPrefix(line, "/accept "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/accept "))
			if user == "" {
				a.ui.Println("Usage: /accept <user>")
				continue
			}

			if a.chatState.Active {
				a.ui.Println("[CHAT] Już jesteś w czacie — nie możesz zaakceptować nowego.")
				continue
			}
			if a.chatState.PendingRequestFrom != user {
				a.ui.Printf("[CHAT] Nie masz oczekującego requestu od %s.\n", user)
				continue
			}

			// Ensure E2EE secret for this peer before accepting.
			if secret, ok := a.secrets.Get(user); ok {
				if err := a.client.SetSharedSecretForPeer(user, secret); err != nil {
					a.ui.Println("[E2EE] Error deriving keys from stored secret:", err)
					continue
				}
				a.ui.Printf("[E2EE] Loaded stored shared secret for %s.\n", user)
			} else {
				secret := a.ui.ReadSharedSecret()
				if err := a.client.SetSharedSecretForPeer(user, secret); err != nil {
					a.ui.Println("[E2EE] Error deriving keys:", err)
					continue
				}
				if err := a.secrets.Set(user, secret); err != nil {
					a.ui.Println("[E2EE] Error saving secret:", err)
					continue
				}
			}

			if err := a.AcceptChat(user); err != nil {
				a.ui.Println("Błąd:", err)
			} else {
				a.ui.Printf("[CHAT] Zaakceptowano czat z %s.\n", user)
			}

		// Refuse chat request
		case strings.HasPrefix(line, "/refuse "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/refuse "))
			if user == "" {
				a.ui.Println("Usage: /refuse <user>")
				continue
			}

			if a.chatState.Active {
				a.ui.Println("[CHAT] Jesteś w czacie — nie możesz odrzucać requestów.")
				continue
			}
			if a.chatState.PendingRequestFrom != user {
				a.ui.Printf("[CHAT] Nie masz oczekującego requestu od %s.\n", user)
				continue
			}

			if err := a.RefuseChat(user, "user refused"); err != nil {
				a.ui.Println("Błąd:", err)
			} else {
				a.ui.Printf("[CHAT] Odrzucono czat z %s.\n", user)
			}

		// Send message in active chat
		case strings.HasPrefix(line, "/msg "):
			if !a.chatState.Active {
				a.ui.Println("[CHAT] Nie jesteś w czacie. Użyj /chat <user>.")
				continue
			}

			text := strings.TrimSpace(strings.TrimPrefix(line, "/msg "))
			if text == "" {
				continue
			}

			if err := a.SendTextMessage(text); err != nil {
				a.ui.Println("Błąd:", err)
			}

		// End chat
		case line == "/end":
			if !a.chatState.Active {
				a.ui.Println("[CHAT] Nie jesteś w czacie.")
				continue
			}

			if err := a.EndChat("user ended chat"); err != nil {
				a.ui.Println("Błąd:", err)
			} else {
				a.ui.Println("[CHAT] Zakończono czat.")
			}

		default:
			a.ui.Println("Nieznana komenda.")
		}
	}
}

func (a *App) printMenu() {
	a.ui.Println("Dostępne komendy:")
	a.ui.Println("  /status <user>")
	a.ui.Println("  /secret <user>   # configure shared secret for peer")
	a.ui.Println("  /chat <user>")
	a.ui.Println("  /accept <user>")
	a.ui.Println("  /refuse <user>")
	a.ui.Println("  /msg <tekst>")
	a.ui.Println("  /end")
	a.ui.Println("  /quit")
}
