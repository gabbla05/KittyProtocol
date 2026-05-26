package chat

import "sync"

// ChatState tracks the local chat session state.
// It is UI-facing and independent from the protocol state machine.
type ChatState struct {
	mu sync.Mutex

	Active       bool
	ActiveTarget string

	Pending     bool
	PendingFrom string
}

// NewChatState constructs a new empty chat state.
func NewChatState() *ChatState {
	return &ChatState{}
}

// SetPendingRequest marks that a peer has requested a chat session.
func (s *ChatState) SetPendingRequest(from string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Pending = true
	s.PendingFrom = from
}

// SetActive marks that a chat session with the given peer is active.
func (s *ChatState) SetActive(peer string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Active = true
	s.ActiveTarget = peer
	s.Pending = false
	s.PendingFrom = ""
}

// ClearPendingRequest removes any pending chat request.
func (s *ChatState) ClearPendingRequest() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Pending = false
	s.PendingFrom = ""
}

// EndChat resets the chat state to idle.
func (s *ChatState) EndChat() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Active = false
	s.ActiveTarget = ""
	s.Pending = false
	s.PendingFrom = ""
}

// IsActive returns whether a chat is active and with whom.
func (s *ChatState) IsActive() (bool, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Active, s.ActiveTarget
}

// HasPendingFrom returns true if there is a pending request from the given user.
func (s *ChatState) HasPendingFrom(user string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Pending && s.PendingFrom == user
}

// HasAnyPending returns whether any pending request exists.
func (s *ChatState) HasAnyPending() (bool, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Pending, s.PendingFrom
}
