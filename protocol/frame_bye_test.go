package protocol

import (
	"strings"
	"testing"
)

func TestParseByeFrameValid(t *testing.T) {
	json := []byte(`{"type":"BYE","msg_id":1}`)
	_, err := ParseByeFrame(json)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseByeFrameInvalidJSON(t *testing.T) {
	_, err := ParseByeFrame([]byte(`{invalid}`))
	if err == nil {
		t.Fatalf("expected JSON error")
	}
}

func TestParseByeFrameWrongType(t *testing.T) {
	_, err := ParseByeFrame([]byte(`{"type":"DATA","msg_id":1}`))
	if err == nil || !strings.Contains(err.Error(), ErrCodeInvalidFrame) {
		t.Fatalf("expected wrong type error")
	}
}

func TestParseByeFrameInvalidMsgID(t *testing.T) {
	_, err := ParseByeFrame([]byte(`{"type":"BYE","msg_id":0}`))
	if err == nil {
		t.Fatalf("expected invalid msg_id error")
	}
}
