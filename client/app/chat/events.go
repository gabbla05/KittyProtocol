package chat

import (
	"github.com/gabbla05/KittyProtocol/client/api"
)

// ChatEvent is a UI-agnostic event struct delivered to GUI.
type ChatEvent struct {
	Type   string // "request", "accept", "refuse", "end", "message"
	From   string
	Text   string
	Reason string
}

// ChatEventBridge connects KittyClient event channels with ChatState.
// It remains UI-agnostic. GUI listens on Events channel.
type ChatEventBridge struct {
	client    *api.KittyClient
	chatState *ChatState

	// GUI listens here. CLI ignores it.
	Events chan ChatEvent
}

// NewChatEventBridge constructs a new event bridge.
func NewChatEventBridge(client *api.KittyClient, state *ChatState) *ChatEventBridge {
	return &ChatEventBridge{
		client:    client,
		chatState: state,
		Events:    make(chan ChatEvent, 32),
	}
}

// Run starts a blocking loop that processes chat-related events.
// CLI receives log strings via onEvent.
// GUI receives structured events via Events channel.
func (b *ChatEventBridge) Run(onEvent func(msg string)) {
	for {
		select {

		// -------------------------
		// CHAT_REQUEST
		// -------------------------
		case ev := <-b.client.ChatRequestEvents():
			b.chatState.SetPendingRequest(ev.From)

			// CLI log
			onEvent("[CHAT] " + ev.From + " wants to chat with you. You can accept it or refuse it now")

			// GUI event
			b.Events <- ChatEvent{
				Type: "request",
				From: ev.From,
			}

		// -------------------------
		// CHAT_ACCEPT
		// -------------------------
		case ev := <-b.client.ChatAcceptEvents():
			b.chatState.SetActive(ev.From)

			onEvent("[CHAT] " + ev.From + " accepted the chat.")

			b.Events <- ChatEvent{
				Type: "accept",
				From: ev.From,
			}

		// -------------------------
		// CHAT_REFUSE
		// -------------------------
		case ev := <-b.client.ChatRefuseEvents():
			b.chatState.ClearPendingRequest()

			onEvent("[CHAT] " + ev.From + " refused the chat: " + ev.Reason)

			b.Events <- ChatEvent{
				Type:   "refuse",
				From:   ev.From,
				Reason: ev.Reason,
			}

		// -------------------------
		// CHAT_END
		// -------------------------
		case ev := <-b.client.ChatEndEvents():
			b.chatState.EndChat()

			onEvent("[CHAT] " + ev.From + " ended the chat: " + ev.Reason)

			b.Events <- ChatEvent{
				Type:   "end",
				From:   ev.From,
				Reason: ev.Reason,
			}

		// -------------------------
		// TEXT_MESSAGE
		// -------------------------
		case ev := <-b.client.ChatMessageEvents():
			onEvent("[" + ev.From + "] " + ev.Text)

			b.Events <- ChatEvent{
				Type: "message",
				From: ev.From,
				Text: ev.Text,
			}
		}
	}
}
