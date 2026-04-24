package protocol

import (
	"encoding/json"
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

// Test sprawdzający brak wymaganych pól (Task 8 - Walidacja).
func TestGetFrameTypeMissingFields(t *testing.T) {
	// Scenariusz: Brak msg_id
	jsonNoID := []byte(`{"type":"DATA"}`)
	_, _, err := GetFrameType(jsonNoID)
	if err == nil {
		t.Fatalf("expected error for missing msg_id")
	}

	// Scenariusz: Brak type
	jsonNoType := []byte(`{"msg_id":123}`)
	_, _, err = GetFrameType(jsonNoType)
	if err == nil {
		t.Fatalf("expected error for missing type")
	}
}

// Test dla niepoprawnego formatu JSON (ERR_02).
func TestGetFrameTypeInvalidJSON(t *testing.T) {
	jsonInvalid := []byte(`{invalid json}`)
	_, _, err := GetFrameType(jsonInvalid)
	if err == nil {
		t.Fatalf("expected JSON parsing error")
	}
}

// Test sprawdzający, czy specyficzne ramki poprawnie się unmarshalują.
func TestDataFrameUnmarshal(t *testing.T) {
	importJSON := []byte(`{"type":"DATA","msg_id":123,"target":"bob","payload":"SGVsbG8=","mac":"hash"}`)
	var f DataFrame
	if err := json.Unmarshal(importJSON, &f); err != nil {
		t.Fatalf("failed to unmarshal DataFrame: %v", err)
	}
	if f.Target != "bob" || f.Payload != "SGVsbG8=" {
		t.Errorf("DataFrame has wrong values after unmarshal")
	}
}
