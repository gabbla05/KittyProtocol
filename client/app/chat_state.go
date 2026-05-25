package app

import "sync"

type ChatState struct {
	mu sync.Mutex

	Active       bool
	ActiveTarget string

	Pending     bool
	PendingFrom string
}

func NewChatState() *ChatState {
	return &ChatState{}
}

// Incoming CHAT_REQUEST
func (s *ChatState) SetPendingRequest(from string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Pending = true
	s.PendingFrom = from
}

// User accepted chat
func (s *ChatState) SetActive(peer string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Active = true
	s.ActiveTarget = peer

	s.Pending = false
	s.PendingFrom = ""
}

// User refused chat
func (s *ChatState) ClearPendingRequest() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Pending = false
	s.PendingFrom = ""
}

// Chat ended
func (s *ChatState) EndChat() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Active = false
	s.ActiveTarget = ""

	s.Pending = false
	s.PendingFrom = ""
}

// --- bezpieczne gettery / checki ---

func (s *ChatState) IsActive() (bool, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Active, s.ActiveTarget
}

func (s *ChatState) HasPendingFrom(user string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Pending && s.PendingFrom == user
}

func (s *ChatState) HasAnyPending() (bool, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Pending, s.PendingFrom
}
