package protocol

import (
	"strings"
	"testing"
)

func TestParseErrorFrameValid(t *testing.T) {
	json := []byte(`{"type":"ERROR","msg_id":1,"code":"ERR_01","desc":"x"}`)
	_, err := ParseErrorFrame(json)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseErrorFrameInvalidJSON(t *testing.T) {
	_, err := ParseErrorFrame([]byte(`{invalid}`))
	if err == nil {
		t.Fatalf("expected JSON error")
	}
}

func TestParseErrorFrameWrongType(t *testing.T) {
	_, err := ParseErrorFrame([]byte(`{"type":"DATA","msg_id":1,"code":"ERR_01"}`))
	if err == nil || !strings.Contains(err.Error(), ErrCodeInvalidFrame) {
		t.Fatalf("expected wrong type error")
	}
}

func TestParseErrorFrameInvalidMsgID(t *testing.T) {
	_, err := ParseErrorFrame([]byte(`{"type":"ERROR","msg_id":0,"code":"ERR_01"}`))
	if err == nil {
		t.Fatalf("expected invalid msg_id error")
	}
}

func TestParseErrorFrameMissingCode(t *testing.T) {
	_, err := ParseErrorFrame([]byte(`{"type":"ERROR","msg_id":1}`))
	if err == nil {
		t.Fatalf("expected missing code error")
	}
}
