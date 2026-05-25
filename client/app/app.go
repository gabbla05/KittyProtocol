package app

import (
	"github.com/gabbla05/KittyProtocol/client/api"
)

type UI interface {
	ReadLine() string
	ReadSharedSecret() []byte
	Println(v ...any)
	Printf(format string, v ...any)
	Prompt()
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
	go a.handleChatEvents()

	return a
}

func (a *App) attachEventHandlers() {
	c := a.client

	// ERROR frame
	c.OnError(func(code, desc string) {
		// Jeśli w trakcie aktywnego czatu dostaniemy ERR_15,
		// potraktuj to jako „peer zniknął” i zamknij lokalnie czat.
		if code == "ERR_15" {
			if active, _ := a.chatState.IsActive(); active {
				a.chatState.EndChat()
				a.ui.Printf("\n[CHAT] Czat zakończony (peer unavailable: %s).\n> ", desc)
				return
			}
		}
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
		// Przy rozłączeniu zawsze czyścimy stan czatu lokalnie.
		a.chatState.EndChat()
		a.ui.Printf("\n[DISCONNECTED] %v\n> ", err)
	})
}

func (a *App) InitSecretStoreForUser(username string, masterKey []byte) {
	path := PathForUser(username)
	a.secrets = NewSecretStore(path, masterKey)

	for peer, secret := range a.secrets.All() {
		_ = a.client.SetSharedSecretForPeer(peer, secret)
	}
}

func (a *App) Client() *api.KittyClient      { return a.client }
func (a *App) Secrets() *SecretStore         { return a.secrets }
func (a *App) Disconnected() <-chan struct{} { return a.disconnected }

// Udostępniamy ChatState dla UI (np. do obsługi /quit).
func (a *App) ChatState() *ChatState {
	return a.chatState
}

func (a *App) handleChatEvents() {
	for {
		select {
		case ev := <-a.client.ChatRequestEvents():
			a.chatState.SetPendingRequest(ev.From)
			a.ui.Printf("\n[CHAT] %s chce z Tobą rozmawiać. Użyj /accept %s lub /refuse %s\n",
				ev.From, ev.From, ev.From)
			a.ui.Prompt()

		case ev := <-a.client.ChatAcceptEvents():
			a.chatState.SetActive(ev.From)
			a.ui.Printf("\n[CHAT] %s zaakceptował czat.\n", ev.From)
			a.ui.Prompt()

		case ev := <-a.client.ChatRefuseEvents():
			a.chatState.ClearPendingRequest()
			a.ui.Printf("\n[CHAT] %s odrzucił czat: %s\n", ev.From, ev.Reason)
			a.ui.Prompt()

		case ev := <-a.client.ChatEndEvents():
			a.chatState.EndChat()
			a.ui.Printf("\n[CHAT] %s zakończył czat: %s\n", ev.From, ev.Reason)
			a.ui.Prompt()

		case ev := <-a.client.ChatMessageEvents():
			a.ui.Printf("\n[%s] %s\n", ev.From, ev.Text)
			a.ui.Prompt()
		}
	}
}
