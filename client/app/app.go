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
	a := &App{
		client:       c,
		ui:           ui,
		disconnected: disconnected,
		chatState:    NewChatState(),
		secrets:      nil,
	}

	a.attachEventHandlers()
	return a
}

func (a *App) attachEventHandlers() {
	c := a.client

	// Decrypted DATA → chat logic
	c.RegisterAppPayloadHandler(a.HandleIncomingPayload)

	// ERROR frame
	c.OnError(func(code, desc string) {
		a.ui.Printf("\n[ERROR] %s: %s\n> ", code, desc)
	})

	// STATUS_RES frame
	c.OnStatus(func(target, status string) {
		if target == "" && status == "no_target" {
			a.ui.Printf("\n[CHAT] Czat zakończony.\n> ")
			return
		}
		a.ui.Printf("\n[STATUS] %s is %s\n> ", target, status)
	})

	// Disconnect event
	c.OnDisconnected(func(err error) {
		a.ui.Printf("\n[DISCONNECTED] %v\n> ", err)
	})
}

func (a *App) InitSecretStoreForUser(username string) {
	path := PathForUser(username)
	a.secrets = NewSecretStore(path)

	for peer, secret := range a.secrets.All() {
		_ = a.client.SetSharedSecretForPeer(peer, secret)
	}
}

func (a *App) Client() *api.KittyClient      { return a.client }
func (a *App) Secrets() *SecretStore         { return a.secrets }
func (a *App) Disconnected() <-chan struct{} { return a.disconnected }
