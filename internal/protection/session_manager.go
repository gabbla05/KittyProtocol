package protection

import (
	"fmt"
	"sync"
	"time"
)

// SessionManager manages all active sessions in memory.
// It periodically scans for idle sessions and closes them.
type SessionManager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
}

// NewSessionManager creates a new SessionManager and starts the idle cleaner goroutine.
func NewSessionManager() *SessionManager {
    sm := &SessionManager{
        sessions: make(map[string]*Session),
    }
    go sm.startCleaner()
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

// startCleaner periodically checks for sessions idle for more than 60 seconds
// and closes them.
func (sm *SessionManager) startCleaner() {
    ticker := time.NewTicker(10 * time.Second)
    for range ticker.C {
        sm.mu.Lock()
        for user, sess := range sm.sessions {
            if time.Since(sess.LastActive) > 60*time.Second {
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
    go func() {
        ticker := time.NewTicker(interval)
        for range ticker.C {
            sm.mu.Lock()
            for user, sess := range sm.sessions {
                if time.Since(sess.LastActive) > idle {
                    if sess.CloseFunc != nil {
                        sess.CloseFunc()
                    }
                    delete(sm.sessions, user)
                }
            }
            sm.mu.Unlock()
        }
    }()
    return sm
}
