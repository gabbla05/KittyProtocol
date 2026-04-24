package protocol

import "testing"

// Test parsing a valid frame.
func TestParseFrameValid(t *testing.T) {
    json := []byte(`{"type":"DATA","msg_id":123}`)
    f, err := ParseFrame(json)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if f.Type != "DATA" || f.MsgID != 123 {
        t.Fatalf("parsed frame has wrong values")
    }
}

// Test missing required fields.
func TestParseFrameMissingFields(t *testing.T) {
    json := []byte(`{"type":"DATA"}`)
    _, err := ParseFrame(json)
    if err == nil {
        t.Fatalf("expected error for missing msg_id")
    }
}

// Test invalid JSON.
func TestParseFrameInvalidJSON(t *testing.T) {
    json := []byte(`{invalid json}`)
    _, err := ParseFrame(json)
    if err == nil {
        t.Fatalf("expected JSON parsing error")
    }
}
