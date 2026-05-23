package protocol

import (
	"encoding/json"
	"fmt"
)

// CurrentProtocolVersion defines the version of the KittyProtocol
const CurrentProtocolVersion = "1.0"

// Frame type constants
const (
	FrameTypeHello     = "HELLO"
	FrameTypeAuth      = "AUTH"
	FrameTypeRegister  = "REGISTER"
	FrameTypeData      = "DATA"
	FrameTypeMeowOK    = "MEOW_OK"
	FrameTypeError     = "ERROR"
	FrameTypeGetStatus = "GET_STATUS"
	FrameTypeStatusRes = "STATUS_RES"
	FrameTypePing      = "PING"
	FrameTypeBye       = "BYE"
)

// Common error codes
const (
	ErrCodeInvalidFrame = "ERR_02"
)

// BaseFrame contains fields common to every frame.
type BaseFrame struct {
	Type  string `json:"type"`
	MsgID int64  `json:"msg_id"`
}

// -----------------------------------------------------------------------------
// Frame definitions
// -----------------------------------------------------------------------------

type HelloFrame struct {
	BaseFrame
	Version string `json:"version"`
}

type AuthFrame struct {
	BaseFrame
	User string `json:"user"`
	Pass string `json:"pass"`
}

type DataFrame struct {
	BaseFrame
	Target  string `json:"target,omitempty"`
	Sender  string `json:"sender,omitempty"`
	Payload string `json:"payload"`
	MAC     string `json:"mac"`
}

type MeowOkFrame struct {
	BaseFrame
	Status string `json:"status,omitempty"`
}

type ErrorFrame struct {
	BaseFrame
	Code string `json:"code"`
	Desc string `json:"desc"`
}

type GetStatusFrame struct {
	BaseFrame
	Target string `json:"target"`
}

type StatusResFrame struct {
	BaseFrame
	Target string `json:"target"`
	Status string `json:"status"`
}

type PingFrame struct{ BaseFrame }
type ByeFrame struct{ BaseFrame }

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
// Parsers
// -----------------------------------------------------------------------------

func ParseHelloFrame(data []byte) (*HelloFrame, error) {
	var f HelloFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeHello {
		return nil, fmt.Errorf("%s: Invalid type for HELLO frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in HELLO frame", ErrCodeInvalidFrame)
	}
	if f.Version == "" {
		return nil, fmt.Errorf("%s: Missing version in HELLO frame", ErrCodeInvalidFrame)
	}
	if f.Version != CurrentProtocolVersion {
		return nil, fmt.Errorf("%s: Unsupported protocol version %q", ErrCodeInvalidFrame, f.Version)
	}
	return &f, nil
}

// Shared parser for AUTH and REGISTER
func parseAuthLikeFrame(data []byte, expectedType string) (*AuthFrame, error) {
	var f AuthFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != expectedType {
		return nil, fmt.Errorf("%s: Invalid type for %s frame", ErrCodeInvalidFrame, expectedType)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in %s frame", ErrCodeInvalidFrame, expectedType)
	}
	if f.User == "" || f.Pass == "" {
		return nil, fmt.Errorf("%s: Missing user or pass in %s frame", ErrCodeInvalidFrame, expectedType)
	}
	return &f, nil
}

func ParseAuthFrame(data []byte) (*AuthFrame, error) {
	return parseAuthLikeFrame(data, FrameTypeAuth)
}

func ParseRegisterFrame(data []byte) (*AuthFrame, error) {
	return parseAuthLikeFrame(data, FrameTypeRegister)
}

func ParseDataFrame(data []byte) (*DataFrame, error) {
	var f DataFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeData {
		return nil, fmt.Errorf("%s: Invalid type for DATA frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in DATA frame", ErrCodeInvalidFrame)
	}
	if f.Sender != "" {
		return nil, fmt.Errorf("%s: Sender must be empty in client DATA frame", ErrCodeInvalidFrame)
	}
	if f.Target == "" {
		return nil, fmt.Errorf("%s: Missing target in DATA frame", ErrCodeInvalidFrame)
	}
	if f.Payload == "" || f.MAC == "" {
		return nil, fmt.Errorf("%s: Missing payload or MAC in DATA frame", ErrCodeInvalidFrame)
	}
	return &f, nil
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

func ParseGetStatusFrame(data []byte) (*GetStatusFrame, error) {
	var f GetStatusFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeGetStatus {
		return nil, fmt.Errorf("%s: Invalid type for GET_STATUS frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in GET_STATUS frame", ErrCodeInvalidFrame)
	}
	if f.Target == "" {
		return nil, fmt.Errorf("%s: Missing target in GET_STATUS frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}

func ParseStatusResFrame(data []byte) (*StatusResFrame, error) {
	var f StatusResFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeStatusRes {
		return nil, fmt.Errorf("%s: Invalid type for STATUS_RES frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in STATUS_RES frame", ErrCodeInvalidFrame)
	}
	if f.Target == "" || f.Status == "" {
		return nil, fmt.Errorf("%s: Missing target or status in STATUS_RES frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}

func ParsePingFrame(data []byte) (*PingFrame, error) {
	var f PingFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypePing {
		return nil, fmt.Errorf("%s: Invalid type for PING frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in PING frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}

func ParseByeFrame(data []byte) (*ByeFrame, error) {
	var f ByeFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Type != FrameTypeBye {
		return nil, fmt.Errorf("%s: Invalid type for BYE frame", ErrCodeInvalidFrame)
	}
	if f.MsgID <= 0 {
		return nil, fmt.Errorf("%s: Invalid msg_id in BYE frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}
