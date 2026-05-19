package app

func (a *App) RunChatSession(target string) {
	a.client.SetTarget(target)
	a.ui.Printf("Wybrano rozmówcę: %s\n", target)

	secret := a.ui.ReadSharedSecret()
	if err := a.client.SetSharedSecret(secret); err != nil {
		a.ui.Println("Błąd ustawiania sekretu:", err)
		return
	}

	a.ui.Println("Sekret ustawiony. Możesz pisać. (/quit aby wrócić do menu, /replay aby wysłać ostatnią ramkę)")

	for {
		// sprawdzamy, czy klient nie został rozłączony
		select {
		case <-a.disconnected:
			a.ui.Println("[Client] Disconnected from server. Leaving chat.")
			return
		default:
		}

		line := a.ui.ReadLine()

		switch line {
		case "/quit":
			// kończymy tylko tryb rozmowy, wracamy do menu
			return

		case "/replay":
			if err := a.client.ReplayLastFrame(); err != nil {
				a.ui.Println("Replay error:", err)
			} else {
				a.ui.Println("Replay sent.")
			}

		case "":
			continue

		default:
			if err := a.client.SendMessage(line); err != nil {
				a.ui.Println("Send error:", err)
			}
		}
	}
}
