package api

import "encoding/json"

// handleChatPayload attempts to interpret the decrypted DATA payload as a
// chat control frame (CHAT_REQUEST / CHAT_ACCEPT / CHAT_REFUSE / CHAT_END / TEXT_MESSAGE).
// It returns true if the payload was recognized and dispatched as a chat event.
func (c *KittyClient) handleChatPayload(sender string, plaintext []byte) bool {
	var probe chatFrameProbe
	if err := json.Unmarshal(plaintext, &probe); err != nil || probe.Type == "" {
		return false
	}

	c.mu.Lock()
	chatReqCh := c.chatReqCh
	chatAcceptCh := c.chatAcceptCh
	chatRefuseCh := c.chatRefuseCh
	chatEndCh := c.chatEndCh
	chatMsgCh := c.chatMsgCh
	c.mu.Unlock()

	switch probe.Type {
	case "CHAT_REQUEST":
		if chatReqCh != nil {
			chatReqCh <- ChatRequestEvent{From: probe.From}
		}

	case "CHAT_ACCEPT":
		if chatAcceptCh != nil {
			chatAcceptCh <- ChatAcceptEvent{From: probe.From}
		}

	case "CHAT_REFUSE":
		var p chatRefusePayload
		_ = json.Unmarshal(probe.Payload, &p)
		if chatRefuseCh != nil {
			chatRefuseCh <- ChatRefuseEvent{From: probe.From, Reason: p.Reason}
		}

	case "CHAT_END":
		var p chatEndPayload
		_ = json.Unmarshal(probe.Payload, &p)
		if chatEndCh != nil {
			chatEndCh <- ChatEndEvent{From: probe.From, Reason: p.Reason}
		}

	case "TEXT_MESSAGE":
		var p textMessagePayload
		_ = json.Unmarshal(probe.Payload, &p)
		if chatMsgCh != nil {
			chatMsgCh <- ChatMessageEvent{From: probe.From, Text: p.Text}
		}

	default:
		log(LogWarn, "unknown chat frame type: %s", probe.Type)
	}

	return true
}
