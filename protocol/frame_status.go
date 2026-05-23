package protocol

import (
	"encoding/json"
	"fmt"
)

type GetStatusFrame struct {
	BaseFrame
	Target string `json:"target"`
}

type StatusResFrame struct {
	BaseFrame
	Target string `json:"target"`
	Status string `json:"status"`
}

func ParseGetStatusFrame(data []byte) (*GetStatusFrame, error) {
	var f GetStatusFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeGetStatus {
		return nil, fmt.Errorf("%s: invalid type for GET_STATUS frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: invalid msg_id in GET_STATUS frame", ErrCodeInvalidFrame)
	}
	if f.Target == "" {
		return nil, fmt.Errorf("%s: missing target in GET_STATUS frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}

func ParseStatusResFrame(data []byte) (*StatusResFrame, error) {
	var f StatusResFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeStatusRes {
		return nil, fmt.Errorf("%s: invalid type for STATUS_RES frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: invalid msg_id in STATUS_RES frame", ErrCodeInvalidFrame)
	}
	if f.Target == "" || f.Status == "" {
		return nil, fmt.Errorf("%s: missing target or status in STATUS_RES frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}
