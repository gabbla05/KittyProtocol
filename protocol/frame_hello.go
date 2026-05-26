package protocol

import (
	"encoding/json"
	"fmt"
)

// HelloFrame is the first frame sent by the client.
// It announces the protocol version and initiates the handshake.
type HelloFrame struct {
	BaseFrame
	Version string `json:"version"`
}

// ParseHelloFrame validates and parses a HELLO frame.
func ParseHelloFrame(data []byte) (*HelloFrame, error) {
	var f HelloFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeHello {
		return nil, fmt.Errorf("%s: invalid type for HELLO frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: invalid msg_id in HELLO frame", ErrCodeInvalidFrame)
	}
	if f.Version == "" {
		return nil, fmt.Errorf("%s: missing version in HELLO frame", ErrCodeInvalidFrame)
	}
	if f.Version != CurrentProtocolVersion {
		return nil, fmt.Errorf("%s: unsupported protocol version %q", ErrCodeInvalidFrame, f.Version)
	}
	return &f, nil
}
