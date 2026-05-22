package app

import (
	"github.com/gabbla05/KittyProtocol/client/api"
)

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
	chatState    *ChatState
	secrets      *SecretStore
}

// NewApp creates a new application controller and initializes the secret store.
//
// The secret store is responsible for persisting per-peer shared secrets
// to a local file (e.g. ~/.kitty/secrets.json).
func NewApp(c *api.KittyClient, ui UI, disconnected <-chan struct{}) *App {
	storePath := defaultSecretStorePath()
	secretStore := NewSecretStore(storePath)

	return &App{
		client:       c,
		ui:           ui,
		disconnected: disconnected,
		chatState:    NewChatState(),
		secrets:      secretStore,
	}
}
