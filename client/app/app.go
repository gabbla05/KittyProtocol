package app

import "github.com/gabbla05/KittyProtocol/client/api"

type UI interface {
	ReadLine() string
	ReadSharedSecret() []byte
	Println(v ...any)
	Printf(format string, v ...any)
}

type App struct {
	client       *api.KittyClient
	ui           UI
	disconnected <-chan struct{}
}

func NewApp(c *api.KittyClient, ui UI, disconnected <-chan struct{}) *App {
	return &App{
		client:       c,
		ui:           ui,
		disconnected: disconnected,
	}
}
