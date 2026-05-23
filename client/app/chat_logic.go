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
	if a.chatState.Active {
		return errors.New("chat already active")
	}

	frame := NewChatRequest(a.client.User(), target)
	return a.sendAppFrame(frame)
}

func (a *App) AcceptChat(from string) error {
	frame := NewChatAccept(a.client.User(), from)
	a.chatState.SetActive(from)
	return a.sendAppFrame(frame)
}

func (a *App) RefuseChat(from, reason string) error {
	frame := NewChatRefuse(a.client.User(), from, reason)
	a.chatState.ClearPendingRequest()
	return a.sendAppFrame(frame)
}

func (a *App) EndChat(reason string) error {
	if !a.chatState.Active {
		return errors.New("no active chat")
	}

	target := a.chatState.ActiveTarget
	frame := NewChatEnd(a.client.User(), target, reason)

	a.chatState.EndChat()
	return a.sendAppFrame(frame)
}

func (a *App) SendTextMessage(text string) error {
	if !a.chatState.Active {
		return errors.New("chat not active")
	}
	if text == "" {
		return errors.New("text cannot be empty")
	}

	target := a.chatState.ActiveTarget
	frame := NewTextMessage(a.client.User(), target, text)

	return a.sendAppFrame(frame)
}

func (a *App) sendAppFrame(frame ChatFrame) error {
	// Canonicalize target to match transport‑layer expectations
	frame.To = strings.ToLower(strings.TrimSpace(frame.To))

	data, err := json.Marshal(frame)
	if err != nil {
		return fmt.Errorf("marshal chat frame: %w", err)
	}

	// ALWAYS encrypted — Hub requires MAC
	return a.client.SendAppFrameEncrypted(frame.To, data)
}
