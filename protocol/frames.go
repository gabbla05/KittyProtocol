package protocol

import (
	"encoding/json"
	"fmt"
)

// UniversalFrame is a base structure for all KittyProtocol message types.
// "omitempty" ensures empty fields are not sent over the network.
type UniversalFrame struct {
	Type    string `json:"type"`   // e.g., HELLO, AUTH, DATA, ERROR
	MsgID   int    `json:"msg_id"` // Protects against Replay attacks
	Version string `json:"version,omitempty"`
	Status  string `json:"status,omitempty"`
	User    string `json:"user,omitempty"`
	Passw   string `json:"passw,omitempty"`
	Token   string `json:"token,omitempty"`
	Target  string `json:"target,omitempty"`
	IP      string `json:"ip,omitempty"`
	Port    int    `json:"port,omitempty"`
	Payload string `json:"payload,omitempty"`
	HMAC    string `json:"hmac,omitempty"`
	Code    string `json:"code,omitempty"`
	Desc    string `json:"desc,omitempty"`
}

// ToJSON serializes the struct into a JSON byte array.
func (f *UniversalFrame) ToJSON() []byte {
	bytes, _ := json.Marshal(f)
	return bytes
}

// ParseFrame strictly validates the incoming JSON (ERR_02 format error).
func ParseFrame(data []byte) (*UniversalFrame, error) {
	var frame UniversalFrame
	if err := json.Unmarshal(data, &frame); err != nil {
		return nil, fmt.Errorf("JSON parsing error")
	}
	if frame.Type == "" || frame.MsgID == 0 {
		return nil, fmt.Errorf("missing required fields: type or msg_id")
	}
	return &frame, nil
}
