package app

import (
	"encoding/json"
)

func (a *App) HandleIncomingChatFrame(frame ChatFrame) {
	switch frame.Type {

	case ChatRequest:
		// Jeśli jesteśmy w czacie → ignorujemy
		if a.chatState.Active {
			a.ui.Printf("\n[CHAT] Otrzymano CHAT_REQUEST od %s, ale czat jest już aktywny.\n> ", frame.From)
			return
		}

		a.chatState.SetPendingRequest(frame.From)
		a.ui.Printf("\n[CHAT REQUEST] %s chce z Tobą rozmawiać.\nUżyj: /accept %s lub /refuse %s\n> ",
			frame.From, frame.From, frame.From)

	case ChatAccept:
		a.chatState.SetActive(frame.From)
		a.ui.Printf("\n[CHAT ACCEPTED] %s zaakceptował czat.\n> ", frame.From)

	case ChatRefuse:
		var p ChatRefusePayload
		_ = json.Unmarshal(frame.Payload, &p)
		a.chatState.ClearPendingRequest()
		a.ui.Printf("\n[CHAT REFUSED] %s odrzucił czat: %s\n> ", frame.From, p.Reason)

	case ChatEnd:
		a.chatState.EndChat()
		a.ui.Printf("\n[CHAT ENDED] %s zakończył czat.\n> ", frame.From)

	case TextMessage:
		var p TextMessagePayload
		_ = json.Unmarshal(frame.Payload, &p)
		a.ui.Printf("\n[%s]: %s\n> ", frame.From, p.Text)

	default:
		a.ui.Printf("\n[CHAT] Nieznany typ ramki: %s\n> ", frame.Type)
	}
}
