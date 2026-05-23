package protocol

import (
	"encoding/json"
	"fmt"
)

// DataFrame carries encrypted application payloads between clients.
type DataFrame struct {
	BaseFrame
	Target  string `json:"target"`
	Sender  string `json:"sender,omitempty"`
	Payload string `json:"payload"`
	MAC     string `json:"mac"`
}

func ParseDataFrame(data []byte) (*DataFrame, error) {
	var f DataFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeData {
		return nil, fmt.Errorf("%s: invalid type for DATA frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: invalid msg_id in DATA frame", ErrCodeInvalidFrame)
	}
	if f.Target == "" {
		return nil, fmt.Errorf("%s: missing target in DATA frame", ErrCodeInvalidFrame)
	}
	if f.Payload == "" || f.MAC == "" {
		return nil, fmt.Errorf("%s: missing payload or MAC in DATA frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}
