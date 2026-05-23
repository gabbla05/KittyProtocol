package app

import (
	"github.com/gabbla05/KittyProtocol/client/api"
)

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
	chatState    *ChatState
	secrets      *SecretStore
}

func NewApp(c *api.KittyClient, ui UI, disconnected <-chan struct{}) *App {
	return &App{
		client:       c,
		ui:           ui,
		disconnected: disconnected,
		chatState:    NewChatState(),
		secrets:      nil, // initialized after AUTH
	}
}

// InitSecretStoreForUser must be called AFTER successful AUTH.
func (a *App) InitSecretStoreForUser(username string) {
	path := PathForUser(username)
	a.secrets = NewSecretStore(path)

	// Auto-load all secrets into KittyClient
	for peer, secret := range a.secrets.All() {
		_ = a.client.SetSharedSecretForPeer(peer, secret)
	}
}
