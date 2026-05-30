package state

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app"
)

type UIState struct {
	Window       fyne.Window
	Client       *api.KittyClient
	App          *app.App
	UI           *GuiUI
	SwitchToChat func(target string) // NOWOŚĆ: Funkcja przekazywana z zewnątrz do zmiany okna na widok czatu
}

func (s *UIState) SwitchView(content fyne.CanvasObject) {
	s.Window.SetContent(content)
}

type GuiUI struct {
	OnMessage func(msg string)
	State     *UIState
}

func (g *GuiUI) ReadLine() string         { return "" }
func (g *GuiUI) ReadSharedSecret() []byte { return nil }
func (g *GuiUI) Prompt()                  {}

func (g *GuiUI) handleLog(msg string) {
	cleanMsg := strings.TrimSpace(msg)

	// 1. Wykrywanie przychodzącego zapytania o czat
	if strings.HasPrefix(cleanMsg, "[CHAT] ") && strings.Contains(cleanMsg, " wants to chat with you") {
		parts := strings.Split(cleanMsg, " ")
		if len(parts) >= 2 {
			user := parts[1]

			// Jeśli nie mamy klucza dla tego użytkownika
			if g.State != nil && g.State.Client != nil && !g.State.Client.HasSharedSecret(user) {
				if g.State.App != nil {
					go g.State.App.RefuseChat(user, "User has no shared secret configured")
				}
				dialog.ShowInformation("Chat Request Rejected", "User "+user+" wants to chat, but you have no shared secret configured for them.", g.State.Window)
				return
			}

			// Pokaż okienko typu Confirm z dwoma opcjami
			if g.State != nil && g.State.Window != nil {
				cnf := dialog.NewConfirm("Incoming Chat Request", "User "+user+" wants to chat with you.", func(accept bool) {
					if accept {
						// Użytkownik kliknął ACCEPT
						if g.State.App != nil {
							err := g.State.App.AcceptChat(user)
							if err != nil {
								dialog.ShowError(err, g.State.Window)
							} else if g.State.SwitchToChat != nil {
								// Sukces - przełączamy ekran na czat z tym użytkownikiem!
								g.State.SwitchToChat(user)
							}
						}
					} else {
						// Użytkownik kliknął REFUSE
						if g.State.App != nil {
							g.State.App.RefuseChat(user, "User refused via GUI popup")
						}
					}
				}, g.State.Window)

				// Zmiana domyślnego tekstu przycisków ("Yes/No") na "Accept/Refuse"
				cnf.SetConfirmText("Accept")
				cnf.SetDismissText("Refuse")
				cnf.Show()
				
				return // Przerywamy dalsze wykonywanie, żeby nie pokazać pod spodem zwykłego okienka info
			}
		}
	}

	// 2. Obsługa innych powiadomień (np. druga strona zaakceptowała/odrzuciła czat)
	if strings.HasPrefix(cleanMsg, "[CHAT] ") && g.State != nil && g.State.Window != nil {
		dialog.ShowInformation("Chat Event", cleanMsg, g.State.Window)
	}

	// 3. Zapis do historii czatu
	if g.OnMessage != nil {
		g.OnMessage(msg)
	}
}

func (g *GuiUI) Println(v ...any) {
	msg := fmt.Sprintln(v...)
	fmt.Print(msg)
	g.handleLog(msg)
}

func (g *GuiUI) Printf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Print(msg)
	g.handleLog(msg)
}