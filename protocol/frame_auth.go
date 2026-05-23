package protocol

import (
	"encoding/json"
	"fmt"
)

// AuthFrame is used for both AUTH and REGISTER operations.
type AuthFrame struct {
	BaseFrame
	User string `json:"user"`
	Pass string `json:"pass"`
}

// parseAuthLikeFrame is a shared validator for AUTH and REGISTER frames.
func parseAuthLikeFrame(data []byte, expectedType string) (*AuthFrame, error) {
	var f AuthFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != expectedType {
		return nil, fmt.Errorf("%s: invalid type for %s frame", ErrCodeInvalidFrame, expectedType)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: invalid msg_id in %s frame", ErrCodeInvalidFrame, expectedType)
	}
	if f.User == "" || f.Pass == "" {
		return nil, fmt.Errorf("%s: missing user or pass in %s frame", ErrCodeInvalidFrame, expectedType)
	}
	return &f, nil
}

func ParseAuthFrame(data []byte) (*AuthFrame, error) {
	return parseAuthLikeFrame(data, FrameTypeAuth)
}

func ParseRegisterFrame(data []byte) (*AuthFrame, error) {
	return parseAuthLikeFrame(data, FrameTypeRegister)
}
