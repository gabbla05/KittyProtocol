package protocol

import (
	"strings"
	"testing"
)

// Test sprawdzający poprawne pobranie typu i ID ramki.
func TestGetFrameTypeValid(t *testing.T) {
	jsonInput := []byte(`{"type":"DATA","msg_id":123}`)
	typeName, msgID, err := GetFrameType(jsonInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if typeName != "DATA" || msgID != 123 {
		t.Fatalf("GetFrameType returned wrong values: got %s, %d", typeName, msgID)
	}
}

// Test sprawdzający brak wymaganych pól (Task 27 - Walidacja na obecność type/msg_id).
func TestGetFrameTypeMissingFields(t *testing.T) {
	jsonNoID := []byte(`{"type":"DATA"}`)
	_, _, err := GetFrameType(jsonNoID)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing msg_id, got: %v", err)
	}

	jsonNoType := []byte(`{"msg_id":123}`)
	_, _, err = GetFrameType(jsonNoType)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing type, got: %v", err)
	}
}

// Test dla niepoprawnego formatu JSON (uszkodzona struktura parsera).
func TestGetFrameTypeInvalidJSON(t *testing.T) {
	jsonInvalid := []byte(`{invalid json}`)
	_, _, err := GetFrameType(jsonInvalid)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for invalid JSON parsing")
	}
}

// Test dla nieznanego typu wiadomości
func TestGetFrameTypeUnknownType(t *testing.T) {
	jsonUnknown := []byte(`{"type":"HACK","msg_id":123}`)
	_, _, err := GetFrameType(jsonUnknown)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for unknown frame type")
	}
}

// Test sprawdzający, czy parser AUTH poprawnie odrzuca puste pola.
func TestParseAuthFrameValidation(t *testing.T) {
	jsonMissingPass := []byte(`{"type":"AUTH","msg_id":123,"user":"alice"}`)
	_, err := ParseAuthFrame(jsonMissingPass)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing pass in AUTH")
	}
}

// Test sprawdzający, czy specyficzne ramki DATA poprawnie się walidują.
func TestDataFrameValidation(t *testing.T) {
	importJSON := []byte(`{"type":"DATA","msg_id":123,"target":"bob","payload":"SGVsbG8=","mac":"hash"}`)
	f, err := ParseDataFrame(importJSON)
	if err != nil {
		t.Fatalf("failed to parse valid DataFrame: %v", err)
	}
	if f.Target != "bob" || f.Payload != "SGVsbG8=" {
		t.Errorf("DataFrame has wrong values after unmarshal")
	}

	missingMacJSON := []byte(`{"type":"DATA","msg_id":123,"target":"bob","payload":"SGVsbG8="}`)
	_, err = ParseDataFrame(missingMacJSON)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing MAC in DATA frame")
	}
}

// Test dla parsera StatusResFrame
func TestParseStatusResFrameValidation(t *testing.T) {
	jsonMissingStatus := []byte(`{"type":"STATUS_RES","msg_id":123,"target":"alice"}`)
	_, err := ParseStatusResFrame(jsonMissingStatus)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing status in STATUS_RES")
	}
}

// Test dla parsera ErrorFrame (sprawdza brak wymaganego kodu błędu)
func TestParseErrorFrameValidation(t *testing.T) {
	jsonMissingCode := []byte(`{"type":"ERROR","msg_id":123,"desc":"Something went wrong"}`)
	_, err := ParseErrorFrame(jsonMissingCode)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing code in ERROR frame")
	}
}

// Test dla parsera GetStatusFrame (sprawdza brak targetu)
func TestParseGetStatusFrameValidation(t *testing.T) {
	jsonMissingTarget := []byte(`{"type":"GET_STATUS","msg_id":123}`)
	_, err := ParseGetStatusFrame(jsonMissingTarget)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing target in GET_STATUS frame")
	}
}

// Test dla parsera HelloFrame (sprawdza uszkodzoną strukturę)
func TestParseHelloFrameValidation(t *testing.T) {
	jsonInvalid := []byte(`{"type":"HELLO"`) // brakująca klamra
	_, err := ParseHelloFrame(jsonInvalid)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for invalid JSON in HELLO frame")
	}
}
