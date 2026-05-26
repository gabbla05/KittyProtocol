package hub

import (
	"testing"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

func TestRouter_UnknownTarget(t *testing.T) {
	// Ensure globalSessions exists, but do NOT reassign it to avoid races
	// with Hub goroutines (cleanup, handleClient, etc.).
	if globalSessions == nil {
		globalSessions = protection.NewSessionManager()
	}
	// Clean per‑test state only.
	globalSessions.Remove("alice")
	globalSessions.Remove("bob")

	sender := &protection.Session{ID: "alice"}

	ok := routeData(
		protocol.DataFrame{
			BaseFrame: protocol.BaseFrame{MsgID: 1},
			Target:    "ghost",
			Payload:   "abc",
		},
		sender,
		&fakeStream{},
	)

	if ok {
		t.Fatalf("Routing to unknown target should fail")
	}
}

func TestRouter_OfflineTarget(t *testing.T) {
	if globalSessions == nil {
		globalSessions = protection.NewSessionManager()
	}
	globalSessions.Remove("alice")
	globalSessions.Remove("bob")

	receiver := &protection.Session{ID: "bob", Stream: nil}
	globalSessions.Add("bob", receiver)

	sender := &protection.Session{ID: "alice"}

	ok := routeData(
		protocol.DataFrame{
			BaseFrame: protocol.BaseFrame{MsgID: 1},
			Target:    "bob",
			Payload:   "abc",
		},
		sender,
		&fakeStream{},
	)

	if ok {
		t.Fatalf("Routing to offline target should fail")
	}
}

func TestRouter_Success(t *testing.T) {
	if globalSessions == nil {
		globalSessions = protection.NewSessionManager()
	}
	globalSessions.Remove("alice")
	globalSessions.Remove("bob")

	receiver := &protection.Session{ID: "bob", Stream: &fakeStream{}}
	globalSessions.Add("bob", receiver)

	sender := &protection.Session{ID: "alice"}

	ok := routeData(
		protocol.DataFrame{
			BaseFrame: protocol.BaseFrame{MsgID: 1},
			Target:    "bob",
			Payload:   "abc",
		},
		sender,
		&fakeStream{},
	)

	if !ok {
		t.Fatalf("Routing to online target should succeed")
	}
}
