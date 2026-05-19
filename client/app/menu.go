package app

import "strings"

func (a *App) RunMainMenu() {
	for {
		// sprawdzamy, czy klient nie został rozłączony
		select {
		case <-a.disconnected:
			a.ui.Println("[Client] Disconnected from server. Exiting menu.")
			return
		default:
		}

		a.ui.Println("Dostępne komendy:")
		a.ui.Println("  /status <user>")
		a.ui.Println("  /chat <user>")
		a.ui.Println("  /quit")

		line := a.ui.ReadLine()
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "/status "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/status "))
			if user == "" {
				a.ui.Println("Usage: /status <user>")
				continue
			}
			_ = a.client.SendGetStatus(user)

		case strings.HasPrefix(line, "/chat "):
			user := strings.TrimSpace(strings.TrimPrefix(line, "/chat "))
			if user == "" {
				a.ui.Println("Usage: /chat <user>")
				continue
			}
			a.RunChatSession(user)

		case line == "/quit":
			_ = a.client.SendBye()
			return

		case line == "":
			// pusta linia – po prostu pokaż menu jeszcze raz
			continue

		default:
			a.ui.Println("Nieznana komenda.")
		}
	}
}
