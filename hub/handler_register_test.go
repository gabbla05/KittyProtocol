package hub

import (
	"encoding/json"
	"testing"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

func TestHandleRegisterSuccess(t *testing.T) {
	globalAuth = auth.NewMockAuth()
	globalSessions = protection.NewSessionManager()

	c := &clientContext{
		state:  stateHelloReceived,
		stream: &fakeStream{},
	}

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeRegister, MsgID: 1},
		User:      "alice",
		Pass:      "secret",
	}
	raw, _ := json.Marshal(frame)

	c.handleRegister(raw)

	// REGISTER should NOT authenticate user
	if c.state != stateHelloReceived {
		t.Fatalf("REGISTER should NOT authenticate user")
	}
}
