package hub

import (
	"encoding/json"
	"testing"

	"github.com/gabbla05/KittyProtocol/protocol"
)

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
