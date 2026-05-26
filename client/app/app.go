package app

import (
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app/chat"
)

// UI defines the minimal interface required by App.
// It allows App to remain UI-agnostic (CLI, GUI, Wails, etc.).
type UI interface {
	ReadLine() string
	ReadSharedSecret() []byte
	Println(v ...any)
	Printf(format string, v ...any)
	Prompt()
}

// App is the high-level application layer.
// It wires together:
//   - KittyClient (transport + E2EE)
//   - ChatLogic (chat operations)
//   - ChatState (local chat state)
//   - ChatEventBridge (incoming events)
//   - UI (presentation layer)
type App struct {
	client       *api.KittyClient
	ui           UI
	disconnected <-chan struct{}

	chatState  *chat.ChatState
	chatLogic  *chat.ChatLogic
	chatBridge *chat.ChatEventBridge

	secrets *SecretStore
}

// NewApp constructs a new application layer instance.
func NewApp(c *api.KittyClient, ui UI, disconnected <-chan struct{}) *App {
	state := chat.NewChatState()

	a := &App{
		client:       c,
		ui:           ui,
		disconnected: disconnected,
		chatState:    state,
		chatLogic:    chat.NewChatLogic(c, state),
		chatBridge:   chat.NewChatEventBridge(c, state),
	}

	a.attachCoreEventHandlers()

	// Start chat event loop
	go a.chatBridge.Run(func(msg string) {
		a.ui.Printf("\n%s\n", msg)
		a.ui.Prompt()
	})

	return a
}

func (a *App) attachCoreEventHandlers() {
	c := a.client

	// ERROR frame
	c.OnError(func(code, desc string) {
		if code == "ERR_15" {
			if active, _ := a.chatState.IsActive(); active {
				a.chatState.EndChat()
				a.ui.Printf("\n[CHAT] Chat ended (peer unavailable: %s).\n", desc)
				a.ui.Prompt()
				return
			}
		}
		a.ui.Printf("\n[ERROR] %s: %s\n", code, desc)
		a.ui.Prompt()
	})

	// STATUS_RES frame
	c.OnStatus(func(target, status string) {
		if target == "" && status == "no_target" {
			a.ui.Printf("\n[CHAT] Chat ended.\n")
			a.ui.Prompt()
			return
		}
		a.ui.Printf("\n[STATUS] %s is %s\n", target, status)
		a.ui.Prompt()
	})

	// Disconnect event
	c.OnDisconnected(func(err error) {
		a.chatState.EndChat()
		a.ui.Printf("\n[DISCONNECTED] %v\n", err)
		a.ui.Prompt()
	})
}

// Disconnected returns a channel that is closed when the client disconnects.
func (a *App) Disconnected() <-chan struct{} {
	return a.disconnected
}

// Client exposes the underlying KittyClient for low-level operations
// (e.g. SendBye, SendGetStatus, SetSharedSecretForPeer).
func (a *App) Client() *api.KittyClient {
	return a.client
}

// ChatState exposes the chat state for UI (e.g. to check active chat on /logout).
func (a *App) ChatState() *chat.ChatState {
	return a.chatState
}

// Secrets returns the secret store used for persisting shared secrets.
func (a *App) Secrets() *SecretStore {
	return a.secrets
}

// InitSecretStoreForUser initializes the secret store for a given user and
// loads all stored shared secrets into KittyClient.
func (a *App) InitSecretStoreForUser(username string, masterKey []byte) {
	path := PathForUser(username)
	a.secrets = NewSecretStore(path, masterKey)

	for peer, secret := range a.secrets.All() {
		_ = a.client.SetSharedSecretForPeer(peer, secret)
	}
}

// High-level chat operations — thin wrappers delegating to ChatLogic.

func (a *App) StartChatRequest(target string) error {
	return a.chatLogic.StartChatRequest(target)
}

func (a *App) AcceptChat(from string) error {
	return a.chatLogic.AcceptChat(from)
}

func (a *App) RefuseChat(from, reason string) error {
	return a.chatLogic.RefuseChat(from, reason)
}

func (a *App) EndChat(reason string) error {
	return a.chatLogic.EndChat(reason)
}

func (a *App) SendTextMessage(text string) error {
	return a.chatLogic.SendTextMessage(text)
}
