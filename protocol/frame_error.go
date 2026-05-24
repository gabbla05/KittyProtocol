package protocol

import (
	"encoding/json"
	"fmt"
)

// Standardized error codes used in ERROR frames.
// These codes are intentionally short and stable to avoid breaking clients.
const (
	ErrCodeInvalidFrame = "ERR_02"
)

// ErrorFrame represents a transport-level error returned by the Hub.
// It is used for protocol violations, malformed frames, or invalid state transitions.
type ErrorFrame struct {
	BaseFrame
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func ParseErrorFrame(data []byte) (*ErrorFrame, error) {
	var f ErrorFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeError {
		return nil, fmt.Errorf("%s: Invalid type for ERROR frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in ERROR frame", ErrCodeInvalidFrame)
	}
	if f.Code == "" {
		return nil, fmt.Errorf("%s: Missing error code in ERROR frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}
