package state

import (
	"fyne.io/fyne/v2"
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app"
)

type UIState struct {
	Window       fyne.Window
	Client       *api.KittyClient
	App          *app.App
	UI           *GuiUI
	SwitchToChat func(target string)
}

func (s *UIState) SwitchView(content fyne.CanvasObject) {
	s.Window.SetContent(content)
}

// GuiUI — UI-agnostic adapter dla App.
// NIE parsuje logów. Reaguje tylko na zdarzenia z ChatEventBridge.
type GuiUI struct {
	State *UIState

	// Callbacki zdarzeń czatu
	OnMessage     func(from string, text string)
	OnChatRequest func(from string)
	OnChatAccept  func(from string)
	OnChatRefuse  func(from string, reason string)
	OnChatEnd     func(from string, reason string)
}

// Wymagane przez App.UI
func (g *GuiUI) ReadLine() string               { return "" }
func (g *GuiUI) ReadSharedSecret() []byte       { return nil }
func (g *GuiUI) Prompt()                        {}
func (g *GuiUI) Println(v ...any)               {}
func (g *GuiUI) Printf(format string, v ...any) {}
