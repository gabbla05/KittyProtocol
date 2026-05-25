package hub

import (
	"testing"

	"github.com/gabbla05/KittyProtocol/internal/protection"
)

func TestClientContextCleanup(t *testing.T) {
	globalSessions = protection.NewSessionManager()
	globalSessions.Add("alice", &protection.Session{ID: "alice"})

	sess := &protection.Session{ID: "alice"}
	globalSessions.Add("alice", sess)

	c := &clientContext{
		username: "alice",
		session:  sess,
	}

	c.cleanup()

	if _, ok := globalSessions.Get("alice"); ok {
		t.Fatalf("cleanup should remove session")
	}
}
