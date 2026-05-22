package app

import "sync"

// ChatState holds the local chat session state for a single client.
// It is completely independent from the transport layer and protocol details.
type ChatState struct {
	mu sync.Mutex

	// Active indicates whether a chat session is currently active
	// (after a successful CHAT_ACCEPT exchange).
	Active bool

	// ActiveTarget is the username of the peer we are currently chatting with.
	ActiveTarget string

	// PendingRequestFrom holds the username of a peer who sent us a CHAT_REQUEST
	// that has not yet been accepted or refused.
	PendingRequestFrom string

	// SecretEstablished indicates whether an E2EE shared secret has been
	// configured for the current peer (according to the UI flow).
	SecretEstablished bool
}

// NewChatState creates an empty chat state instance.
func NewChatState() *ChatState {
	return &ChatState{}
}

// SetActive marks the chat as active with the given target and clears any
// pending incoming request.
func (s *ChatState) SetActive(target string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Active = true
	s.ActiveTarget = target
	s.PendingRequestFrom = ""
}

// SetPendingRequest records an incoming CHAT_REQUEST from the given user.
func (s *ChatState) SetPendingRequest(from string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.PendingRequestFrom = from
}

// ClearPendingRequest clears any pending CHAT_REQUEST information.
func (s *ChatState) ClearPendingRequest() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.PendingRequestFrom = ""
}

// EndChat resets the active chat state and clears any pending request.
func (s *ChatState) EndChat() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Active = false
	s.ActiveTarget = ""
	s.PendingRequestFrom = ""
}

// SetSecretEstablished updates the E2EE secret flag for the current context.
func (s *ChatState) SetSecretEstablished(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.SecretEstablished = v
}
