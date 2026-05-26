package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gabbla05/KittyProtocol/client/api"
)

// ChatLogic contains all application-level chat operations.
// It is UI-agnostic and uses KittyClient only for encrypted transport.
type ChatLogic struct {
	client    *api.KittyClient
	chatState *ChatState
}

// NewChatLogic constructs a new chat logic layer.
func NewChatLogic(client *api.KittyClient, state *ChatState) *ChatLogic {
	return &ChatLogic{client: client, chatState: state}
}

// StartChatRequest initiates a chat session with a peer.
func (l *ChatLogic) StartChatRequest(target string) error {
	if target == "" {
		return errors.New("target cannot be empty")
	}
	if target == l.client.User() {
		return errors.New("cannot chat with yourself")
	}
	if active, peer := l.chatState.IsActive(); active {
		return fmt.Errorf("chat already active with %s", peer)
	}
	if pending, from := l.chatState.HasAnyPending(); pending {
		return fmt.Errorf("you have a pending request from %s — resolve it first", from)
	}
	if !l.client.HasSharedSecret(target) {
		return fmt.Errorf("no shared secret for %s", target)
	}

	frame := NewChatRequest(l.client.User(), target)
	return l.sendFrame(frame)
}

// AcceptChat accepts a pending chat request.
func (l *ChatLogic) AcceptChat(from string) error {
	if from == "" {
		return errors.New("from cannot be empty")
	}
	if from == l.client.User() {
		return errors.New("cannot chat with yourself")
	}
	if !l.chatState.HasPendingFrom(from) {
		return fmt.Errorf("no pending chat request from %s", from)
	}

	frame := NewChatAccept(l.client.User(), from)
	if err := l.sendFrame(frame); err != nil {
		return err
	}

	l.chatState.SetActive(from)
	return nil
}

// RefuseChat rejects a pending chat request.
func (l *ChatLogic) RefuseChat(from, reason string) error {
	if from == "" {
		return errors.New("from cannot be empty")
	}
	if from == l.client.User() {
		return errors.New("cannot chat with yourself")
	}
	if !l.chatState.HasPendingFrom(from) {
		return fmt.Errorf("no pending chat request from %s", from)
	}

	frame := NewChatRefuse(l.client.User(), from, reason)
	if err := l.sendFrame(frame); err != nil {
		return err
	}

	l.chatState.ClearPendingRequest()
	return nil
}

// EndChat terminates an active chat session.
func (l *ChatLogic) EndChat(reason string) error {
	active, peer := l.chatState.IsActive()
	if !active {
		return errors.New("no active chat")
	}
	if peer == "" {
		return errors.New("no active target")
	}
	if peer == l.client.User() {
		return errors.New("cannot chat with yourself")
	}

	frame := NewChatEnd(l.client.User(), peer, reason)
	if err := l.sendFrame(frame); err != nil {
		return err
	}

	l.chatState.EndChat()
	return nil
}

// SendTextMessage sends a text message inside an active chat session.
func (l *ChatLogic) SendTextMessage(text string) error {
	if text == "" {
		return errors.New("text cannot be empty")
	}

	active, peer := l.chatState.IsActive()
	if !active {
		return errors.New("chat not active")
	}
	if peer == "" {
		return errors.New("no active target")
	}
	if peer == l.client.User() {
		return errors.New("cannot chat with yourself")
	}

	frame := NewTextMessage(l.client.User(), peer, text)
	return l.sendFrame(frame)
}

// sendFrame serializes and sends a chat frame via encrypted API transport.
func (l *ChatLogic) sendFrame(frame ChatFrame) error {
	frame.To = strings.ToLower(strings.TrimSpace(frame.To))

	data, err := json.Marshal(frame)
	if err != nil {
		return fmt.Errorf("marshal chat frame: %w", err)
	}

	return l.client.SendAppFrameEncrypted(frame.To, data)
}
