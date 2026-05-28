package state

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app"
)

type UIState struct {
	Window fyne.Window
	Client *api.KittyClient
	App    *app.App
	UI     *GuiUI // NOWOŚĆ: dostęp do naszego interfejsu logów
}

func (s *UIState) SwitchView(content fyne.CanvasObject) {
	s.Window.SetContent(content)
}

type GuiUI struct {
	OnMessage func(msg string) // Callback, który wywoła się przy nowej wiadomości
}

func (g *GuiUI) ReadLine() string         { return "" }
func (g *GuiUI) ReadSharedSecret() []byte { return nil }
func (g *GuiUI) Prompt()                  {}

func (g *GuiUI) Println(v ...any) {
	msg := fmt.Sprintln(v...)
	fmt.Print(msg) // Zostawiamy log w konsoli
	if g.OnMessage != nil {
		g.OnMessage(msg) // Przekazujemy log do okienka GUI!
	}
}

func (g *GuiUI) Printf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Print(msg) // Zostawiamy log w konsoli
	if g.OnMessage != nil {
		g.OnMessage(msg) // Przekazujemy log do okienka GUI!
	}
}