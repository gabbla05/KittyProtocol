package hub

import (
	"encoding/json"
	"testing"

	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

func TestHandleGetStatus(t *testing.T) {
	globalSessions = protection.NewSessionManager()
	globalSessions.Add("bob", &protection.Session{ID: "bob"})

	c := &clientContext{
		state:  stateAuthenticated,
		stream: &fakeStream{},
	}

	frame := protocol.GetStatusFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeGetStatus, MsgID: 1},
		Target:    "bob",
	}
	raw, _ := json.Marshal(frame)

	c.handleGetStatus(raw)
}
