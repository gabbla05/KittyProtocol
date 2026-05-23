package protection

import (
	"fmt"
	"sync"
	"time"
)

const (
	DefaultSessionIdleTimeout     = 60 * time.Second
	DefaultSessionCleanupInterval = 10 * time.Second
)

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	stopChan chan struct{}
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
		stopChan: make(chan struct{}),
	}
	go sm.startCleaner(DefaultSessionCleanupInterval, DefaultSessionIdleTimeout)
	return sm
}

func (sm *SessionManager) Stop() {
	close(sm.stopChan)
}

func (sm *SessionManager) Add(user string, sess *Session) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[user] = sess
}

func (sm *SessionManager) Get(user string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.sessions[user]
	return s, ok
}

func (sm *SessionManager) Remove(user string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	fmt.Println("[SessionManager] Removing session for:", user)
	delete(sm.sessions, user)
}

func (sm *SessionManager) startCleaner(interval, idleTimeout time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-sm.stopChan:
			return

		case <-ticker.C:
			sm.mu.Lock()
			for user, sess := range sm.sessions {
				if time.Since(sess.LastActive) > idleTimeout {
					fmt.Printf("[SessionManager: Protection] Idle Timeout: %s. Removing session.\n", user)
					if sess.CloseFunc != nil {
						sess.CloseFunc()
					}
					delete(sm.sessions, user)
				}
			}
			sm.mu.Unlock()
		}
	}
}

func (sm *SessionManager) IsOnline(user string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.sessions[user]
	return ok
}
