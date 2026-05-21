package app

import "strings"

// RunMainMenu displays the main command loop.
// It blocks until the user quits or the client disconnects.
func (a *App) RunMainMenu() {
	for {
		// Check for disconnection
		select {
		case <-a.disconnected:
			a.ui.Println("[Client] Rozłączono z serwerem. Zamykanie aplikacji.")
			return
		default:
		}

		a.printMenu()

		line := strings.TrimSpace(a.ui.ReadLine())

		switch {
		case line == "":
			continue

		case line == "/quit":
			_ = a.client.SendBye()
			return

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

		default:
			a.ui.Println("Nieznana komenda.")
		}
	}
}

// printMenu prints the list of available commands.
func (a *App) printMenu() {
	a.ui.Println("Dostępne komendy:")
	a.ui.Println("  /status <user>")
	a.ui.Println("  /chat <user>")
	a.ui.Println("  /quit")
}
