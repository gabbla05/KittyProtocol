package hub

import (
	"encoding/json"
	"testing"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
)

// --- HELLO tests ---

func TestHandleHello(t *testing.T) {
	// We must provide a fake stream, otherwise handleHello will panic
	c := &clientContext{
		state:  stateInit,
		stream: &fakeStream{},
	}

	frame := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeHello, MsgID: 1},
		Version:   "1.0",
	}
	raw, _ := json.Marshal(frame)

	c.handleHello(raw)

	if c.state != stateHelloReceived {
		t.Fatalf("HELLO should transition to stateHelloReceived")
	}
}

// --- AUTH tests ---

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

// --- DATA tests ---

func TestHandleDataBeforeAuth(t *testing.T) {
	c := &clientContext{
		state:  stateHelloReceived,
		stream: &fakeStream{},
	}

	frame := protocol.DataFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeData, MsgID: 1},
		Target:    "bob",
		Payload:   "abc",
	}
	raw, _ := json.Marshal(frame)

	c.handleData(raw)

	if c.state != stateHelloReceived {
		t.Fatalf("DATA before AUTH should not change state")
	}
}
