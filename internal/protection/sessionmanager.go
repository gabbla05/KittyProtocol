package protection

import (
	"fmt"
	"sync"
	"time"
)

type Session struct {
	ID         string
	LastActive time.Time
	Limiter    *RateLimiter
	CloseFunc  func() // Funkcja do zamykania połączenia przez Hub
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
	}
	go sm.startCleaner() // Task 9: Demon czyszczący [cite: 581]
	return sm
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

// startCleaner co 10s sprawdza, kto przekroczył 60s bezczynności [cite: 662]
func (sm *SessionManager) startCleaner() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		sm.mu.Lock()
		for user, sess := range sm.sessions {
			if time.Since(sess.LastActive) > 60*time.Second {
				fmt.Printf("[Protection] Idle Timeout: %s. Usuwanie sesji.\n", user)
				sess.CloseFunc()
				delete(sm.sessions, user)
			}
		}
		sm.mu.Unlock()
	}
}
