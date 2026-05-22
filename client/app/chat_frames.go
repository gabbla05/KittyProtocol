package app

import "encoding/json"

// App frames types
type ChatFrameType string

const (
	ChatRequest ChatFrameType = "CHAT_REQUEST"
	ChatAccept  ChatFrameType = "CHAT_ACCEPT"
	ChatRefuse  ChatFrameType = "CHAT_REFUSE"
	ChatEnd     ChatFrameType = "CHAT_END"
	TextMessage ChatFrameType = "TEXT_MESSAGE"
)

// General app frame struccture
type ChatFrame struct {
	Type    ChatFrameType   `json:"type"`
	From    string          `json:"from"`
	To      string          `json:"to"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Payloady dla poszczególnych ramek

type ChatRequestPayload struct {
	// Na razie puste — w przyszłości można dodać np. "topic"
}

type ChatAcceptPayload struct {
	// Też puste — można dodać np. "sessionID"
}

type ChatRefusePayload struct {
	Reason string `json:"reason,omitempty"`
}

type ChatEndPayload struct {
	Reason string `json:"reason,omitempty"`
}

type TextMessagePayload struct {
	Text string `json:"text"`
}

// Helpery do tworzenia ramek

func NewChatRequest(from, to string) ChatFrame {
	payload, _ := json.Marshal(ChatRequestPayload{})
	return ChatFrame{
		Type:    ChatRequest,
		From:    from,
		To:      to,
		Payload: payload,
	}
}

func NewChatAccept(from, to string) ChatFrame {
	payload, _ := json.Marshal(ChatAcceptPayload{})
	return ChatFrame{
		Type:    ChatAccept,
		From:    from,
		To:      to,
		Payload: payload,
	}
}

func NewChatRefuse(from, to, reason string) ChatFrame {
	payload, _ := json.Marshal(ChatRefusePayload{Reason: reason})
	return ChatFrame{
		Type:    ChatRefuse,
		From:    from,
		To:      to,
		Payload: payload,
	}
}

func NewChatEnd(from, to, reason string) ChatFrame {
	payload, _ := json.Marshal(ChatEndPayload{Reason: reason})
	return ChatFrame{
		Type:    ChatEnd,
		From:    from,
		To:      to,
		Payload: payload,
	}
}

func NewTextMessage(from, to, text string) ChatFrame {
	payload, _ := json.Marshal(TextMessagePayload{Text: text})
	return ChatFrame{
		Type:    TextMessage,
		From:    from,
		To:      to,
		Payload: payload,
	}
}
