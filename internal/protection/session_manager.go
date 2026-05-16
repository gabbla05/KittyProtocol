package protection

import (
	"fmt"
	"sync"
	"time"
)

// DefaultSessionIdleTimeout defines how long a session may stay inactive
// before it is considered idle and removed.
const DefaultSessionIdleTimeout = 60 * time.Second

// DefaultSessionCleanupInterval defines how often the SessionManager
// scans for idle sessions.
const DefaultSessionCleanupInterval = 10 * time.Second

// SessionManager manages all active sessions in memory.
// It periodically scans for idle sessions and closes them.
// This component is purely transport-level and does not contain
// any application-layer chat logic.
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewSessionManager creates a new SessionManager and starts the idle cleaner goroutine.
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
	}
	go sm.startCleaner(DefaultSessionCleanupInterval, DefaultSessionIdleTimeout)
	return sm
}

// Add registers a new session for the given user.
func (sm *SessionManager) Add(user string, sess *Session) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[user] = sess
}

// Get retrieves a session by username.
func (sm *SessionManager) Get(user string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.sessions[user]
	return s, ok
}

// Remove deletes a session from the manager.
func (sm *SessionManager) Remove(user string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	fmt.Println("[SessionManager] Removing session for:", user)
	delete(sm.sessions, user)
}

// startCleaner periodically checks for sessions idle for more than idleTimeout
// and closes them. This ensures resource cleanup and prevents stale sessions.
func (sm *SessionManager) startCleaner(interval, idleTimeout time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		sm.mu.Lock()
		for user, sess := range sm.sessions {
			if time.Since(sess.LastActive) > idleTimeout {
				fmt.Printf("[Protection] Idle Timeout: %s. Removing session.\n", user)
				if sess.CloseFunc != nil {
					sess.CloseFunc()
				}
				delete(sm.sessions, user)
			}
		}
		sm.mu.Unlock()
	}
}

// NewSessionManagerWithInterval is used only for tests.
func NewSessionManagerWithInterval(interval time.Duration, idle time.Duration) *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
	}
	go sm.startCleaner(interval, idle)
	return sm
}

// IsOnline returns true if there is an active session for the given user.
func (sm *SessionManager) IsOnline(user string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	_, ok := sm.sessions[user]
	return ok
}
