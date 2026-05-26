package protocol

import (
	"strings"
	"testing"
)

func TestGetFrameTypeValid(t *testing.T) {
	jsonInput := []byte(`{"type":"DATA","msg_id":123}`)
	typ, id, err := GetFrameType(jsonInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if typ != FrameTypeData || id != 123 {
		t.Fatalf("wrong values: %s %d", typ, id)
	}
}

func TestGetFrameTypeMissingFields(t *testing.T) {
	_, _, err := GetFrameType([]byte(`{"type":"DATA"}`))
	if err == nil || !strings.Contains(err.Error(), ErrCodeInvalidFrame) {
		t.Fatalf("expected missing msg_id error")
	}

	_, _, err = GetFrameType([]byte(`{"msg_id":1}`))
	if err == nil {
		t.Fatalf("expected missing type error")
	}
}

func TestGetFrameTypeInvalidJSON(t *testing.T) {
	_, _, err := GetFrameType([]byte(`{invalid json}`))
	if err == nil {
		t.Fatalf("expected JSON error")
	}
}

func TestGetFrameTypeUnknownType(t *testing.T) {
	_, _, err := GetFrameType([]byte(`{"type":"HACK","msg_id":1}`))
	if err == nil {
		t.Fatalf("expected unknown type error")
	}
}
