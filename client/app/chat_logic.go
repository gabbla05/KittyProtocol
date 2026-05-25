package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (a *App) StartChatRequest(target string) error {
	if target == "" {
		return errors.New("target cannot be empty")
	}

	if active, peer := a.chatState.IsActive(); active {
		return fmt.Errorf("chat already active with %s", peer)
	}

	if pending, from := a.chatState.HasAnyPending(); pending {
		return fmt.Errorf("you have a pending request from %s — resolve it first", from)
	}

	if !a.client.HasSharedSecret(target) {
		return fmt.Errorf("no shared secret for %s", target)
	}

	frame := NewChatRequest(a.client.User(), target)
	return a.sendAppFrame(frame)
}

func (a *App) AcceptChat(from string) error {
	if from == "" {
		return errors.New("from cannot be empty")
	}

	if !a.chatState.HasPendingFrom(from) {
		return fmt.Errorf("no pending chat request from %s", from)
	}

	frame := NewChatAccept(a.client.User(), from)

	if err := a.sendAppFrame(frame); err != nil {
		return err
	}

	// Responder wchodzi w stan Active lokalnie
	a.chatState.SetActive(from)
	return nil
}

func (a *App) RefuseChat(from, reason string) error {
	if from == "" {
		return errors.New("from cannot be empty")
	}

	if !a.chatState.HasPendingFrom(from) {
		return fmt.Errorf("no pending chat request from %s", from)
	}

	frame := NewChatRefuse(a.client.User(), from, reason)

	if err := a.sendAppFrame(frame); err != nil {
		return err
	}

	a.chatState.ClearPendingRequest()
	return nil
}

func (a *App) EndChat(reason string) error {
	active, peer := a.chatState.IsActive()
	if !active {
		return errors.New("no active chat")
	}
	if peer == "" {
		return errors.New("no active target")
	}

	frame := NewChatEnd(a.client.User(), peer, reason)

	if err := a.sendAppFrame(frame); err != nil {
		return err
	}

	a.chatState.EndChat()
	return nil
}

func (a *App) SendTextMessage(text string) error {
	if text == "" {
		return errors.New("text cannot be empty")
	}

	active, peer := a.chatState.IsActive()
	if !active {
		return errors.New("chat not active")
	}
	if peer == "" {
		return errors.New("no active target")
	}

	frame := NewTextMessage(a.client.User(), peer, text)
	return a.sendAppFrame(frame)
}

func (a *App) sendAppFrame(frame ChatFrame) error {
	frame.To = strings.ToLower(strings.TrimSpace(frame.To))

	data, err := json.Marshal(frame)
	if err != nil {
		return fmt.Errorf("marshal chat frame: %w", err)
	}

	return a.client.SendAppFrameEncrypted(frame.To, data)
}
