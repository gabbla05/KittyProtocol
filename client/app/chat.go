package app

// RunChatSession enters an interactive chat loop with the selected target.
// The function blocks until the user exits the chat or the client disconnects.
func (a *App) RunChatSession(target string) {
	a.client.SetTarget(target)
	_ = a.client.SendGetStatus(target)
	a.ui.Printf("Wybrano rozmówcę: %s\n", target)

	secret := a.ui.ReadSharedSecret()
	if err := a.client.SetSharedSecret(secret); err != nil {
		a.ui.Println("Błąd ustawiania sekretu:", err)
		return
	}

	a.ui.Println("Sekret ustawiony. Możesz pisać.")
	a.ui.Println("Komendy: /quit (wyjście), /replay (wyślij ostatnią ramkę)")

	for {
		// Check for disconnection
		select {
		case <-a.disconnected:
			a.ui.Println("[Client] Rozłączono z serwerem. Zamykanie czatu.")
			return
		default:
		}

		line := a.ui.ReadLine()

		switch line {
		case "":
			continue

		case "/quit":
			// 1. Send to the hub information about the end of conversation
			_ = a.client.SendGetStatus("") // target = "" means ther e is no receiver

			// 2. Clean up target locally
			a.client.SetTarget("")

			// 3. Exit to menu
			a.printMenu()
			return

		case "/replay":
			if err := a.client.ReplayLastFrame(); err != nil {
				a.ui.Println("Replay error:", err)
			} else {
				a.ui.Println("Replay sent.")
			}

		default:
			if err := a.client.SendMessage(line); err != nil {
				a.ui.Println("Send error:", err)
			}
		}
	}
}
