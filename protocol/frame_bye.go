package protocol

import (
	"encoding/json"
	"fmt"
)

type ByeFrame struct{ BaseFrame }

func ParseByeFrame(data []byte) (*ByeFrame, error) {
	var f ByeFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeBye {
		return nil, fmt.Errorf("%s: invalid type for BYE frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: invalid msg_id in BYE frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}
