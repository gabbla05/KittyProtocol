package chat

import "encoding/json"

// ChatFrameType enumerates all application-level chat control frames.
// These frames are encrypted and transported inside DATA frames at the API layer.
type ChatFrameType string

const (
	ChatRequest ChatFrameType = "CHAT_REQUEST"
	ChatAccept  ChatFrameType = "CHAT_ACCEPT"
	ChatRefuse  ChatFrameType = "CHAT_REFUSE"
	ChatEnd     ChatFrameType = "CHAT_END"
	TextMessage ChatFrameType = "TEXT_MESSAGE"
)

// ChatFrame is the generic application-level chat frame.
// It is serialized to JSON and encrypted by KittyClient before sending.
type ChatFrame struct {
	Type    ChatFrameType   `json:"type"`
	From    string          `json:"from"`
	To      string          `json:"to"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Payload structures for each chat frame type.

type ChatRequestPayload struct{}

type ChatAcceptPayload struct{}

type ChatRefusePayload struct {
	Reason string `json:"reason,omitempty"`
}

type ChatEndPayload struct {
	Reason string `json:"reason,omitempty"`
}

type TextMessagePayload struct {
	Text string `json:"text"`
}

// Constructors for each chat frame type.

func NewChatRequest(from, to string) ChatFrame {
	payload, _ := json.Marshal(ChatRequestPayload{})
	return ChatFrame{Type: ChatRequest, From: from, To: to, Payload: payload}
}

func NewChatAccept(from, to string) ChatFrame {
	payload, _ := json.Marshal(ChatAcceptPayload{})
	return ChatFrame{Type: ChatAccept, From: from, To: to, Payload: payload}
}

func NewChatRefuse(from, to, reason string) ChatFrame {
	payload, _ := json.Marshal(ChatRefusePayload{Reason: reason})
	return ChatFrame{Type: ChatRefuse, From: from, To: to, Payload: payload}
}

func NewChatEnd(from, to, reason string) ChatFrame {
	payload, _ := json.Marshal(ChatEndPayload{Reason: reason})
	return ChatFrame{Type: ChatEnd, From: from, To: to, Payload: payload}
}

func NewTextMessage(from, to, text string) ChatFrame {
	payload, _ := json.Marshal(TextMessagePayload{Text: text})
	return ChatFrame{Type: TextMessage, From: from, To: to, Payload: payload}
}
