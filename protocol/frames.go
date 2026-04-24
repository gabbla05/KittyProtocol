package protocol

import (
	"encoding/json"
	"fmt"
)

// BaseFrame zawiera pola wspólne dla każdej ramki
type BaseFrame struct {
	Type  string `json:"type"`   // np. "HELLO", "AUTH", "DATA"
	MsgID int64  `json:"msg_id"` // Timestamp jako unikalny ID
}

// 1. HELLO - Przywitanie
type HelloFrame struct {
	BaseFrame
	Version string `json:"version"` // np. "1.0"
}

// 2. AUTH - Uwierzytelnianie
type AuthFrame struct {
	BaseFrame
	User string `json:"user"` // Login użytkownika
	Pass string `json:"pass"` // Hasło użytkownika
}

// 3. DATA - Przesyłanie ładunku E2EE
type DataFrame struct {
	BaseFrame
	Target  string `json:"target,omitempty"` // Odbiorca (u nadawcy)
	Sender  string `json:"sender,omitempty"` // Nadawca (u odbiorcy - dodawane przez Hub)
	Payload string `json:"payload"`          // Zaszyfrowany Base64
	MAC     string `json:"mac"`              // HMAC dla integralności E2EE
}

// 4. MEOW_OK - Potwierdzenie aplikacyjne (ACK)
type MeowOkFrame struct {
	BaseFrame
	Status string `json:"status,omitempty"` // Opcjonalny opis statusu
}

// 5. ERROR - Ramka błędu
type ErrorFrame struct {
	BaseFrame
	Code string `json:"code"` // Kod błędu (np. ERR_02)
	Desc string `json:"desc"` // Opis błędu
}

// 6. GET_STATUS - Zapytanie o status użytkownika
type GetStatusFrame struct {
	BaseFrame
	Target string `json:"target"` // Kogo sprawdzamy
}

// 7. STATUS_RES - Odpowiedź o statusie
type StatusResFrame struct {
	BaseFrame
	Target string `json:"target"` // Identyfikator sprawdzanego
	Status string `json:"status"` // "online" lub "offline"
}

// 8. PING i 9. BYE - Keep-alive i zakończenie
type PingFrame struct{ BaseFrame }
type ByeFrame struct{ BaseFrame }

// GetFrameType pozwala sprawdzić typ ramki przed pełnym parsowaniem i rygorystycznie odrzuca błędy
func GetFrameType(data []byte) (string, int64, error) {
	var base BaseFrame
	if err := json.Unmarshal(data, &base); err != nil {
		return "", 0, fmt.Errorf("ERR_02: JSON parsing error")
	}
	// Rygorystyczna walidacja pól wymaganych
	if base.Type == "" || base.MsgID == 0 {
		return "", 0, fmt.Errorf("ERR_02: missing required fields (type/msg_id)")
	}
	// Weryfikacja, czy typ zgadza się z obsługiwanym enumem
	if !IsValidType(base.Type) {
		return "", 0, fmt.Errorf("ERR_02: unknown or invalid frame type")
	}
	return base.Type, base.MsgID, nil
}

// IsValidType sprawdza, czy typ wiadomości znajduje się na liście dozwolonych
func IsValidType(t string) bool {
	switch t {
	case "HELLO", "AUTH", "DATA", "MEOW_OK", "ERROR", "GET_STATUS", "STATUS_RES", "PING", "BYE":
		return true
	}
	return false
}

// ParseHelloFrame waliduje ramkę startową
func ParseHelloFrame(data []byte) (*HelloFrame, error) {
	var f HelloFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("ERR_02: Invalid JSON format")
	}
	return &f, nil
}

// ParseAuthFrame rygorystycznie waliduje ramkę logowania
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

// ParseDataFrame waliduje ramkę z wiadomością i weryfikuje niezbędne do E2EE pola
func ParseDataFrame(data []byte) (*DataFrame, error) {
	var f DataFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("ERR_02: Invalid JSON format")
	}
	if f.Payload == "" || f.MAC == "" {
		return nil, fmt.Errorf("ERR_02: Missing payload or MAC in DATA frame")
	}
	return &f, nil
}

// ParseErrorFrame waliduje formatkę błędu
func ParseErrorFrame(data []byte) (*ErrorFrame, error) {
	var f ErrorFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("ERR_02: Invalid JSON format")
	}
	if f.Code == "" {
		return nil, fmt.Errorf("ERR_02: Missing error code in ERROR frame")
	}
	return &f, nil
}

// ParseGetStatusFrame waliduje zapytanie o status
func ParseGetStatusFrame(data []byte) (*GetStatusFrame, error) {
	var f GetStatusFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("ERR_02: Invalid JSON format")
	}
	if f.Target == "" {
		return nil, fmt.Errorf("ERR_02: Missing target in GET_STATUS frame")
	}
	return &f, nil
}

// ParseStatusResFrame waliduje odpowiedź o statusie
func ParseStatusResFrame(data []byte) (*StatusResFrame, error) {
	var f StatusResFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("ERR_02: Invalid JSON format")
	}
	if f.Target == "" || f.Status == "" {
		return nil, fmt.Errorf("ERR_02: Missing target or status in STATUS_RES frame")
	}
	return &f, nil
}
