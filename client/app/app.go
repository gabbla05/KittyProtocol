package app

import "github.com/gabbla05/KittyProtocol/client/api"

// UI defines the minimal interface required by the App layer.
// It allows plugging in different frontends (CLI, GUI, tests).
type UI interface {
	ReadLine() string
	ReadSharedSecret() []byte
	Println(v ...any)
	Printf(format string, v ...any)
}

// App coordinates user interaction (UI) with the KittyClient API.
// It contains no networking or cryptography — only application logic.
type App struct {
	client       *api.KittyClient
	ui           UI
	disconnected <-chan struct{}
}

// NewApp creates a new application controller.
func NewApp(c *api.KittyClient, ui UI, disconnected <-chan struct{}) *App {
	return &App{
		client:       c,
		ui:           ui,
		disconnected: disconnected,
	}
}
