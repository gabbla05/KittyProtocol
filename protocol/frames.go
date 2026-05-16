package protocol

import (
	"encoding/json"
	"fmt"
)

// Frame type constants – single source of truth for all frame type strings.
const (
	FrameTypeHello     = "HELLO"
	FrameTypeAuth      = "AUTH"
	FrameTypeData      = "DATA"
	FrameTypeMeowOK    = "MEOW_OK"
	FrameTypeError     = "ERROR"
	FrameTypeGetStatus = "GET_STATUS"
	FrameTypeStatusRes = "STATUS_RES"
	FrameTypePing      = "PING"
	FrameTypeBye       = "BYE"
)

// Common error code constants used in protocol-level validation.
const (
	ErrCodeInvalidFrame = "ERR_02" // generic "invalid frame" / "bad format" error
)

// BaseFrame contains fields common to every frame.
type BaseFrame struct {
	Type  string `json:"type"`   // e.g. "HELLO", "AUTH", "DATA"
	MsgID int64  `json:"msg_id"` // Timestamp used as a unique ID
}

// 1. HELLO – initial greeting frame.
type HelloFrame struct {
	BaseFrame
	Version string `json:"version"` // e.g. "1.0"
}

// 2. AUTH – authentication frame.
type AuthFrame struct {
	BaseFrame
	User string `json:"user"` // Username
	Pass string `json:"pass"` // Password
}

// 3. DATA – E2EE payload transfer frame.
type DataFrame struct {
	BaseFrame
	Target  string `json:"target,omitempty"` // Recipient (on sender side)
	Sender  string `json:"sender,omitempty"` // Sender (on receiver side – added by Hub)
	Payload string `json:"payload"`          // Encrypted Base64 payload
	MAC     string `json:"mac"`              // HMAC for E2EE integrity
}

// 4. MEOW_OK – application-level acknowledgment (ACK).
type MeowOkFrame struct {
	BaseFrame
	Status string `json:"status,omitempty"` // Optional status description
}

// 5. ERROR – error frame.
type ErrorFrame struct {
	BaseFrame
	Code string `json:"code"` // Error code (e.g. ERR_02)
	Desc string `json:"desc"` // Error description
}

// 6. GET_STATUS – query for user status.
type GetStatusFrame struct {
	BaseFrame
	Target string `json:"target"` // User whose status is being queried
}

// 7. STATUS_RES – response with user status.
type StatusResFrame struct {
	BaseFrame
	Target string `json:"target"` // Queried user identifier
	Status string `json:"status"` // "online" or "offline"
}

// 8. PING and 9. BYE – keep-alive and session termination.
type PingFrame struct{ BaseFrame }
type ByeFrame struct{ BaseFrame }

// GetFrameType performs a lightweight parse to extract frame type and msg_id,
// and strictly rejects malformed or incomplete frames.
func GetFrameType(data []byte) (string, int64, error) {
	var base BaseFrame
	if err := json.Unmarshal(data, &base); err != nil {
		return "", 0, fmt.Errorf("%s: JSON parsing error", ErrCodeInvalidFrame)
	}
	// Strict validation of required fields.
	if base.Type == "" || base.MsgID == 0 {
		return "", 0, fmt.Errorf("%s: missing required fields (type/msg_id)", ErrCodeInvalidFrame)
	}
	// Verify that the type is one of the supported protocol types.
	if !IsValidType(base.Type) {
		return "", 0, fmt.Errorf("%s: unknown or invalid frame type", ErrCodeInvalidFrame)
	}
	return base.Type, base.MsgID, nil
}

// IsValidType checks whether the given frame type is allowed by the protocol.
func IsValidType(t string) bool {
	switch t {
	case FrameTypeHello,
		FrameTypeAuth,
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

// ParseHelloFrame validates the initial HELLO frame.
func ParseHelloFrame(data []byte) (*HelloFrame, error) {
	var f HelloFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	return &f, nil
}

// ParseAuthFrame strictly validates the AUTH frame.
func ParseAuthFrame(data []byte) (*AuthFrame, error) {
	var f AuthFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.User == "" || f.Pass == "" {
		return nil, fmt.Errorf("%s: Missing user or pass in AUTH frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}

// ParseDataFrame validates the DATA frame and checks required E2EE fields.
func ParseDataFrame(data []byte) (*DataFrame, error) {
	var f DataFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Payload == "" || f.MAC == "" {
		return nil, fmt.Errorf("%s: Missing payload or MAC in DATA frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}

// ParseErrorFrame validates the ERROR frame.
func ParseErrorFrame(data []byte) (*ErrorFrame, error) {
	var f ErrorFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Code == "" {
		return nil, fmt.Errorf("%s: Missing error code in ERROR frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}

// ParseGetStatusFrame validates the GET_STATUS frame.
func ParseGetStatusFrame(data []byte) (*GetStatusFrame, error) {
	var f GetStatusFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Target == "" {
		return nil, fmt.Errorf("%s: Missing target in GET_STATUS frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}

// ParseStatusResFrame validates the STATUS_RES frame.
func ParseStatusResFrame(data []byte) (*StatusResFrame, error) {
	var f StatusResFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: Invalid JSON format", ErrCodeInvalidFrame)
	}
	if f.Target == "" || f.Status == "" {
		return nil, fmt.Errorf("%s: Missing target or status in STATUS_RES frame", ErrCodeInvalidFrame)
	}
	return &f, nil
}
