package protocol

import (
	"encoding/json"
	"fmt"
)

// UniversalFrame is a generic JSON frame used by KittyProtocol.
// It is intentionally flexible for the first implementation stage.
// Later it can be split into dedicated types per message kind.
type UniversalFrame struct {
    Type    string `json:"type"`              // e.g. "HELLO", "AUTH", "DATA", "ERROR"
    MsgID   int64  `json:"msg_id"`            // Timestamp-based message identifier
    User    string `json:"user,omitempty"`    // For AUTH
    Pass    string `json:"pass,omitempty"`    // For AUTH
    Token   string `json:"token,omitempty"`   // Reserved for future session tokens
    Target  string `json:"target,omitempty"`  // For DATA / GET_STATUS
    Sender  string `json:"sender,omitempty"`  // Filled by Hub for delivered DATA
    Payload string `json:"payload,omitempty"` // Encrypted payload (Base64) or plaintext for now
    HMAC    string `json:"hmac,omitempty"`    // E2E integrity (future)
    Status  string `json:"status,omitempty"`  // Human-readable status (e.g. "Logged in")
    Code    string `json:"code,omitempty"`    // Error code (e.g. "ERR_02")
    Desc    string `json:"desc,omitempty"`    // Error description
}

// ToJSON serializes the frame into JSON bytes.
// Errors are ignored here for simplicity; in production they should be handled.
func (f *UniversalFrame) ToJSON() []byte {
    b, _ := json.Marshal(f)
    return b
}

// ParseFrame parses raw JSON bytes into a UniversalFrame and performs basic validation.
// It enforces presence of "type" and "msg_id" as required by the protocol.
func ParseFrame(data []byte) (*UniversalFrame, error) {
    var frame UniversalFrame
    if err := json.Unmarshal(data, &frame); err != nil {
        return nil, fmt.Errorf("ERR_02: JSON parsing error: %w", err)
    }
    if frame.Type == "" || frame.MsgID == 0 {
        return nil, fmt.Errorf("ERR_02: missing required fields (type/msg_id)")
    }
    return &frame, nil
}
