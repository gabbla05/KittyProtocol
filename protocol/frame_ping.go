package protocol

import (
	"encoding/json"
	"fmt"
)

type PingFrame struct{ BaseFrame }

func ParsePingFrame(data []byte) (*PingFrame, error) {
	var f PingFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypePing {
		return nil, fmt.Errorf("%s: invalid type for PING frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: invalid msg_id in PING frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}
