package protocol

import (
	"strings"
	"testing"
)

func TestParseGetStatusFrameInvalidJSON(t *testing.T) {
	_, err := ParseGetStatusFrame([]byte(`{invalid}`))
	if err == nil {
		t.Fatalf("expected JSON error")
	}
}

func TestParseGetStatusFrameWrongType(t *testing.T) {
	_, err := ParseGetStatusFrame([]byte(`{"type":"DATA","msg_id":1,"target":"x"}`))
	if err == nil || !strings.Contains(err.Error(), ErrCodeInvalidFrame) {
		t.Fatalf("expected wrong type error")
	}
}

func TestParseGetStatusFrameInvalidMsgID(t *testing.T) {
	_, err := ParseGetStatusFrame([]byte(`{"type":"GET_STATUS","msg_id":0,"target":"x"}`))
	if err == nil {
		t.Fatalf("expected invalid msg_id error")
	}
}

func TestParseStatusResFrameInvalidJSON(t *testing.T) {
	_, err := ParseStatusResFrame([]byte(`{invalid}`))
	if err == nil {
		t.Fatalf("expected JSON error")
	}
}

func TestParseStatusResFrameWrongType(t *testing.T) {
	_, err := ParseStatusResFrame([]byte(`{"type":"DATA","msg_id":1,"target":"x","status":"y"}`))
	if err == nil {
		t.Fatalf("expected wrong type error")
	}
}

func TestParseStatusResFrameInvalidMsgID(t *testing.T) {
	_, err := ParseStatusResFrame([]byte(`{"type":"STATUS_RES","msg_id":0,"target":"x","status":"y"}`))
	if err == nil {
		t.Fatalf("expected invalid msg_id error")
	}
}
