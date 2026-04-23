package protocol

import (
	"encoding/json"
	"fmt"
)

// BaseFrame zawiera pola wspólne dla każdej ramki [cite: 2143]
type BaseFrame struct {
	Type  string `json:"type"`   // np. "HELLO", "AUTH", "DATA" [cite: 2144]
	MsgID int64  `json:"msg_id"` // Timestamp jako unikalny ID [cite: 2146]
}

type HelloFrame struct {
	BaseFrame
	Version string `json:"version"` // np. "1.0"
}

// 2. AUTH - Uwierzytelnianie [cite: 2132]
type AuthFrame struct {
	BaseFrame
	User string `json:"user"` // Login użytkownika [cite: 2150]
	Pass string `json:"pass"` // Hasło użytkownika [cite: 2150]
}

// 3. DATA - Przesyłanie ładunku E2EE [cite: 2136]
type DataFrame struct {
	BaseFrame
	Target  string `json:"target,omitempty"` // Odbiorca (u nadawcy) [cite: 2150]
	Sender  string `json:"sender,omitempty"` // Nadawca (u odbiorcy - dodawane przez Hub) [cite: 2224]
	Payload string `json:"payload"`          // Zaszyfrowany Base64 [cite: 2160]
	MAC     string `json:"mac"`              // HMAC dla integralności E2EE [cite: 2161]
}

// 4. MEOW_OK - Potwierdzenie aplikacyjne (ACK) [cite: 2137]
type MeowOkFrame struct {
	BaseFrame
	Status string `json:"status,omitempty"` // Opcjonalny opis statusu [cite: 2196]
}

// 5. ERROR - Ramka błędu [cite: 2138]
type ErrorFrame struct {
	BaseFrame
	Code string `json:"code"` // Kod błędu (np. ERR_02) [cite: 2151]
	Desc string `json:"desc"` // Opis błędu [cite: 2151]
}

// 6. GET_STATUS - Zapytanie o status użytkownika [cite: 2133]
type GetStatusFrame struct {
	BaseFrame
	Target string `json:"target"` // Kogo sprawdzamy [cite: 2150]
}

// 7. STATUS_RES - Odpowiedź o statusie [cite: 2135]
type StatusResFrame struct {
	BaseFrame
	Target string `json:"target"` // Identyfikator sprawdzanego [cite: 2207]
	Status string `json:"status"` // "online" lub "offline" [cite: 2207]
}

// 8. PING i 9. BYE - Keep-alive i zakończenie [cite: 2138, 2139]
type PingFrame struct{ BaseFrame }
type ByeFrame struct{ BaseFrame }

// GetFrameType pozwala sprawdzić typ ramki przed pełnym parsowaniem (Task 8)
func GetFrameType(data []byte) (string, int64, error) {
	var base BaseFrame
	if err := json.Unmarshal(data, &base); err != nil {
		return "", 0, fmt.Errorf("ERR_02: JSON parsing error")
	}
	// Rygorystyczna walidacja pól wymaganych (Task 8)
	if base.Type == "" || base.MsgID == 0 {
		return "", 0, fmt.Errorf("ERR_02: missing required fields (type/msg_id)")
	}
	return base.Type, base.MsgID, nil
}

// IsValidType sprawdza, czy typ wiadomości znajduje się na liście dozwolonych [cite: 4254]
func IsValidType(t string) bool {
	switch t {
	case "HELLO", "AUTH", "DATA", "MEOW_OK", "ERROR", "GET_STATUS", "STATUS_RES", "PING", "BYE":
		return true
	}
	return false
}

// ParseAuthFrame rygorystycznie waliduje ramkę logowania (Task 8) [cite: 4254]
func ParseAuthFrame(data []byte) (*AuthFrame, error) {
	var f AuthFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("ERR_02: Invalid JSON format")
	}
	if f.User == "" || f.Pass == "" {
		return nil, fmt.Errorf("ERR_02: Missing user or pass in AUTH frame")
	}
	return &f, nil
}

// ParseDataFrame waliduje ramkę z wiadomością (Task 8) [cite: 4254]
func ParseDataFrame(data []byte) (*DataFrame, error) {
	var f DataFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("ERR_02: Invalid JSON format")
	}
	// Payload i MAC są krytyczne dla E2EE
	if f.Payload == "" || f.MAC == "" {
		return nil, fmt.Errorf("ERR_02: Missing payload or MAC in DATA frame")
	}
	return &f, nil
}
