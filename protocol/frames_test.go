package protocol

import (
	"strings"
	"testing"
)

// --- GetFrameType tests ---

// TestGetFrameTypeValid verifies that GetFrameType correctly extracts
// the frame type and message ID from a valid JSON frame.
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

// TestGetFrameTypeMissingFields ensures that missing required fields
// (type or msg_id) result in an ERR_02 validation error.
func TestGetFrameTypeMissingFields(t *testing.T) {
	jsonNoID := []byte(`{"type":"DATA"}`)
	_, _, err := GetFrameType(jsonNoID)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing msg_id")
	}

	jsonNoType := []byte(`{"msg_id":123}`)
	_, _, err = GetFrameType(jsonNoType)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing type")
	}
}

// TestGetFrameTypeInvalidJSON checks that malformed JSON is rejected.
func TestGetFrameTypeInvalidJSON(t *testing.T) {
	jsonInvalid := []byte(`{invalid json}`)
	_, _, err := GetFrameType(jsonInvalid)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for invalid JSON")
	}
}

// TestGetFrameTypeUnknownType verifies that unknown frame types
// are rejected with ERR_02.
func TestGetFrameTypeUnknownType(t *testing.T) {
	jsonUnknown := []byte(`{"type":"HACK","msg_id":123}`)
	_, _, err := GetFrameType(jsonUnknown)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for unknown frame type")
	}
}

// --- AUTH tests ---

// TestParseAuthFrameValidation ensures that missing required fields
// in AUTH frames (e.g., pass) produce an ERR_02 error.
func TestParseAuthFrameValidation(t *testing.T) {
	jsonMissingPass := []byte(`{"type":"AUTH","msg_id":123,"user":"alice"}`)
	_, err := ParseAuthFrame(jsonMissingPass)
	if err == nil || !strings.Contains(err.Error(), "ERR_02") {
		t.Fatalf("expected ERR_02 for missing pass")
	}
}

// --- DATA tests ---

// TestDataFrameValidation verifies correct parsing of DATA frames,
// including validation of required fields and rejection of invalid ones.
func TestDataFrameValidation(t *testing.T) {
	valid := []byte(`{"type":"DATA","msg_id":123,"target":"bob","payload":"SGVsbG8=","mac":"hash"}`)
	f, err := ParseDataFrame(valid)
	if err != nil {
		t.Fatalf("failed to parse valid DataFrame: %v", err)
	}
	if f.Target != "bob" {
		t.Errorf("wrong target")
	}

	// Sender is allowed (Hub sets it when forwarding)
	withSender := []byte(`{"type":"DATA","msg_id":123,"sender":"alice","target":"bob","payload":"x","mac":"y"}`)
	_, err = ParseDataFrame(withSender)
	if err != nil {
		t.Fatalf("sender should be allowed in forwarded DATA frames: %v", err)
	}

	// Missing MAC
	missingMac := []byte(`{"type":"DATA","msg_id":123,"target":"bob","payload":"x"}`)
	_, err = ParseDataFrame(missingMac)
	if err == nil {
		t.Fatalf("expected error for missing MAC")
	}
}

// --- STATUS_RES tests ---

// TestParseStatusResFrameValidation checks that STATUS_RES frames
// must include both target and status fields.
func TestParseStatusResFrameValidation(t *testing.T) {
	jsonMissingStatus := []byte(`{"type":"STATUS_RES","msg_id":123,"target":"alice"}`)
	_, err := ParseStatusResFrame(jsonMissingStatus)
	if err == nil {
		t.Fatalf("expected error for missing status")
	}
}

// --- ERROR tests ---

// TestParseErrorFrameValidation ensures that ERROR frames must include
// an error code.
func TestParseErrorFrameValidation(t *testing.T) {
	jsonMissingCode := []byte(`{"type":"ERROR","msg_id":123,"desc":"Something went wrong"}`)
	_, err := ParseErrorFrame(jsonMissingCode)
	if err == nil {
		t.Fatalf("expected error for missing code")
	}
}

// --- GET_STATUS tests ---

// TestParseGetStatusFrameValidation verifies that GET_STATUS frames
// must include a target field.
func TestParseGetStatusFrameValidation(t *testing.T) {
	jsonMissingTarget := []byte(`{"type":"GET_STATUS","msg_id":123}`)
	_, err := ParseGetStatusFrame(jsonMissingTarget)
	if err == nil {
		t.Fatalf("expected error for missing target")
	}
}

// --- HELLO tests ---

// TestParseHelloFrameValidation checks HELLO frame validation,
// including JSON format, version presence, and version compatibility.
func TestParseHelloFrameValidation(t *testing.T) {
	invalidJSON := []byte(`{"type":"HELLO"`)
	_, err := ParseHelloFrame(invalidJSON)
	if err == nil {
		t.Fatalf("expected error for invalid JSON")
	}

	missingVersion := []byte(`{"type":"HELLO","msg_id":1}`)
	_, err = ParseHelloFrame(missingVersion)
	if err == nil {
		t.Fatalf("expected error for missing version")
	}

	wrongVersion := []byte(`{"type":"HELLO","msg_id":1,"version":"9.9"}`)
	_, err = ParseHelloFrame(wrongVersion)
	if err == nil {
		t.Fatalf("expected error for unsupported version")
	}
}
