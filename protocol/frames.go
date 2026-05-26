package protocol

import (
	"encoding/json"
	"fmt"
)

// BaseFrame contains fields common to every frame exchanged in the protocol.
// All frames must embed BaseFrame as their first field.
type BaseFrame struct {
	Type  string `json:"type"`
	MsgID int64  `json:"msg_id"`
}

// IsValidType returns true if the provided frame type is recognized by the protocol.
func IsValidType(t string) bool {
	switch t {
	case FrameTypeHello,
		FrameTypeAuth,
		FrameTypeRegister,
		FrameTypeData,
		FrameTypeMeowOK,
		FrameTypeError,
		FrameTypeGetStatus,
		FrameTypeStatusRes,
		FrameTypePing,
		FrameTypeBye:
		return true
	}
	return false
}

// GetFrameType extracts the "type" and "msg_id" fields from raw JSON.
// It performs minimal validation and is used by dispatchers before full parsing.
func GetFrameType(data []byte) (string, int64, error) {
	var base BaseFrame
	if err := json.Unmarshal(data, &base); err != nil {
		return "", 0, fmt.Errorf("%s: JSON parsing error", ErrCodeInvalidFrame)
	}
	if base.Type == "" || base.MsgID <= 0 {
		return "", 0, fmt.Errorf("%s: missing or invalid fields (type/msg_id)", ErrCodeInvalidFrame)
	}
	if !IsValidType(base.Type) {
		return "", 0, fmt.Errorf("%s: unknown or invalid frame type", ErrCodeInvalidFrame)
	}
	return base.Type, base.MsgID, nil
}
