package chat

import (
	"github.com/gabbla05/KittyProtocol/client/api"
)

// ChatEventBridge connects KittyClient event channels with ChatState.
// It is UI-agnostic; UI is notified via callbacks in App.
type ChatEventBridge struct {
	client    *api.KittyClient
	chatState *ChatState
}

// NewChatEventBridge constructs a new event bridge.
func NewChatEventBridge(client *api.KittyClient, state *ChatState) *ChatEventBridge {
	return &ChatEventBridge{client: client, chatState: state}
}

// Run starts a blocking loop that processes chat-related events.
func (b *ChatEventBridge) Run(onEvent func(msg string)) {
	for {
		select {
		case ev := <-b.client.ChatRequestEvents():
			b.chatState.SetPendingRequest(ev.From)
			onEvent("[CHAT] " + ev.From + " wants to chat with you. You can accept it or refuse it now")

		case ev := <-b.client.ChatAcceptEvents():
			b.chatState.SetActive(ev.From)
			onEvent("[CHAT] " + ev.From + " accepted the chat.")

		case ev := <-b.client.ChatRefuseEvents():
			b.chatState.ClearPendingRequest()
			onEvent("[CHAT] " + ev.From + " refused the chat: " + ev.Reason)

		case ev := <-b.client.ChatEndEvents():
			b.chatState.EndChat()
			onEvent("[CHAT] " + ev.From + " ended the chat: " + ev.Reason)

		case ev := <-b.client.ChatMessageEvents():
			onEvent("[" + ev.From + "] " + ev.Text)
		}
	}
}
