package protocol

import (
	"encoding/json"
	"fmt"
)

// UniversalFrame to podstawowa struktura dla wszystkich typów komunikatów.
type UniversalFrame struct {
	Type    string `json:"type"`              // np. "HELLO", "AUTH", "DATA", "ERROR"
	MsgID   int64  `json:"msg_id"`            // Timestamp jako identyfikator [cite: 233]
	User    string `json:"user,omitempty"`    // Dla AUTH
	Pass    string `json:"pass,omitempty"`    // Dla AUTH
	Token   string `json:"token,omitempty"`   // Token sesyjny po AUTH
	Target  string `json:"target,omitempty"`  // Adresat ramki DATA lub GET_STATUS
	Sender  string `json:"sender,omitempty"`  // Wstrzykiwane przez Hub dla odbiorcy [cite: 590]
	Payload string `json:"payload,omitempty"` // Zaszyfrowana treść (Base64)
	HMAC    string `json:"hmac,omitempty"`    // Integralność E2E [cite: 248]
	Status  string `json:"status,omitempty"`  // np. "Logged in", "Ready for auth"
	Code    string `json:"code,omitempty"`    // Kod błędu (np. ERR_02)
	Desc    string `json:"desc,omitempty"`    // Opis błędu
}

func (f *UniversalFrame) ToJSON() []byte {
	bytes, _ := json.Marshal(f)
	return bytes
}

// ParseFrame rygorystycznie waliduje format JSON (Task 8 Gaby)[cite: 584].
func ParseFrame(data []byte) (*UniversalFrame, error) {
	var frame UniversalFrame
	if err := json.Unmarshal(data, &frame); err != nil {
		return nil, fmt.Errorf("ERR_02: JSON parsing error")
	}
	if frame.Type == "" || frame.MsgID == 0 {
		return nil, fmt.Errorf("ERR_02: missing required fields (type/msg_id)")
	}
	return &frame, nil
}
