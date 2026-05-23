package app

import "sync"

// ChatState holds the current chat state on the client side.
// It is fully independent from the transport layer.
type ChatState struct {
	mu sync.Mutex

	// Whether a chat session is currently active (after CHAT_ACCEPT).
	Active bool

	// Logical username of the current chat peer.
	ActiveTarget string

	// If someone sent us a CHAT_REQUEST, this stores the sender username.
	PendingRequestFrom string
}

// NewChatState creates an empty chat state.
func NewChatState() *ChatState {
	return &ChatState{}
}

// SetActive marks a chat as active with the given target.
func (s *ChatState) SetActive(target string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Active = true
	s.ActiveTarget = target
	s.PendingRequestFrom = ""
}

// SetPendingRequest records an incoming chat request.
func (s *ChatState) SetPendingRequest(from string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.PendingRequestFrom = from
}

// ClearPendingRequest clears any pending chat request.
func (s *ChatState) ClearPendingRequest() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.PendingRequestFrom = ""
}

// EndChat ends the current chat and clears state.
func (s *ChatState) EndChat() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Active = false
	s.ActiveTarget = ""
	s.PendingRequestFrom = ""
}
