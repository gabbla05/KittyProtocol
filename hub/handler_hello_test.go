package hub

import (
	"encoding/json"
	"testing"

	"github.com/gabbla05/KittyProtocol/protocol"
)

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

func TestHandleHelloWrongVersion(t *testing.T) {
	c := &clientContext{
		state:  stateInit,
		stream: &fakeStream{},
	}

	frame := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeHello, MsgID: 1},
		Version:   "9.9",
	}
	raw, _ := json.Marshal(frame)

	c.handleHello(raw)

	if c.state != stateInit {
		t.Fatalf("HELLO with wrong version should not change state")
	}
}

func TestHandleHelloTwice(t *testing.T) {
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
	c.handleHello(raw) // drugi raz

	if c.state != stateHelloReceived {
		t.Fatalf("Second HELLO should not break state machine")
	}
}
