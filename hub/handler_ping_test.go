package hub

import (
	"encoding/json"
	"testing"

	"github.com/gabbla05/KittyProtocol/protocol"
)

func TestHandlePing(t *testing.T) {
	c := &clientContext{
		state:  stateAuthenticated,
		stream: &fakeStream{},
	}

	frame := protocol.PingFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypePing, MsgID: 1},
	}
	raw, _ := json.Marshal(frame)

	c.handlePing(raw)

	if c.state != stateAuthenticated {
		t.Fatalf("PING should not change state")
	}
}
