package hub

import (
	"encoding/json"
	"testing"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

func TestHandleAuthBeforeHello(t *testing.T) {
	globalAuth = auth.NewMockAuth()

	c := &clientContext{
		state:  stateInit,
		stream: &fakeStream{},
	}

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 1},
		User:      "alice",
		Pass:      "secret",
	}
	raw, _ := json.Marshal(frame)

	c.handleAuth(raw)

	if c.state != stateInit {
		t.Fatalf("AUTH before HELLO should not change state")
	}
}

func TestHandleAuthWrongPassword(t *testing.T) {
	c := &clientContext{
		state:  stateHelloReceived,
		stream: &fakeStream{},
	}

	globalAuth = auth.NewMockAuth()

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  protocol.FrameTypeAuth,
			MsgID: 1,
		},
		User: "alice",
		Pass: "wrong",
	}

	raw, _ := json.Marshal(frame)

	c.handleAuth(raw)

	if c.state != stateHelloReceived {
		t.Fatalf("AUTH with wrong password should not change state")
	}
}

func TestHandleAuthSuccess(t *testing.T) {
	globalAuth = auth.NewMockAuth()
	globalSessions = protection.NewSessionManager()

	c := &clientContext{
		state:  stateHelloReceived,
		stream: &fakeStream{},
	}

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 1},
		User:      "alice",
		Pass:      "secret",
	}
	raw, _ := json.Marshal(frame)

	c.handleAuth(raw)

	if c.state != stateAuthenticated {
		t.Fatalf("AUTH success should transition to stateAuthenticated")
	}
}

func TestHandleAuthUnknownUser(t *testing.T) {
	globalAuth = auth.NewMockAuth()
	globalSessions = protection.NewSessionManager()

	c := &clientContext{
		state:  stateHelloReceived,
		stream: &fakeStream{},
	}

	frame := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 1},
		User:      "ghost",
		Pass:      "whatever",
	}
	raw, _ := json.Marshal(frame)

	c.handleAuth(raw)

	if c.state != stateHelloReceived {
		t.Fatalf("AUTH with unknown user should not authenticate")
	}
}
