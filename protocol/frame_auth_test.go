package protocol

import (
	"strings"
	"testing"
)

func TestParseAuthFrameInvalidJSON(t *testing.T) {
	_, err := ParseAuthFrame([]byte(`{invalid}`))
	if err == nil {
		t.Fatalf("expected JSON error")
	}
}

func TestParseAuthFrameWrongType(t *testing.T) {
	_, err := ParseAuthFrame([]byte(`{"type":"DATA","msg_id":1,"user":"a","pass":"b"}`))
	if err == nil {
		t.Fatalf("expected wrong type error")
	}
}

func TestParseAuthFrameInvalidMsgID(t *testing.T) {
	_, err := ParseAuthFrame([]byte(`{"type":"AUTH","msg_id":0,"user":"a","pass":"b"}`))
	if err == nil {
		t.Fatalf("expected invalid msg_id error")
	}
}

func TestParseRegisterFrameMissingFields(t *testing.T) {
	_, err := ParseRegisterFrame([]byte(`{"type":"REGISTER","msg_id":1,"user":"a"}`))
	if err == nil || !strings.Contains(err.Error(), ErrCodeInvalidFrame) {
		t.Fatalf("expected missing pass error")
	}
}

func TestParseRegisterFrameWrongType(t *testing.T) {
	_, err := ParseRegisterFrame([]byte(`{"type":"AUTH","msg_id":1,"user":"a","pass":"b"}`))
	if err == nil {
		t.Fatalf("expected wrong type error")
	}
}
