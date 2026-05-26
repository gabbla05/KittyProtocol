package hub

import (
	"encoding/json"
	"testing"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

func TestHandleBye(t *testing.T) {
	globalSessions = protection.NewSessionManager()

	sess := &protection.Session{ID: "alice"}
	globalSessions.Add("alice", sess)

	c := &clientContext{
		username: "alice",
		session:  sess,
		stream:   &fakeStream{},
		state:    stateAuthenticated, // BYE requires AUTH
	}

	frame := protocol.ByeFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeBye, MsgID: 1},
	}
	raw, _ := json.Marshal(frame)

	c.handleBye(raw)

	if _, ok := globalSessions.Get("alice"); ok {
		t.Fatalf("BYE should remove session")
	}
}
