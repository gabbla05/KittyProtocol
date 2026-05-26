package protocol

import (
	"strings"
	"testing"
)

func TestParsePingFrameValid(t *testing.T) {
	json := []byte(`{"type":"PING","msg_id":1}`)
	_, err := ParsePingFrame(json)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParsePingFrameInvalidJSON(t *testing.T) {
	_, err := ParsePingFrame([]byte(`{invalid}`))
	if err == nil {
		t.Fatalf("expected JSON error")
	}
}

func TestParsePingFrameWrongType(t *testing.T) {
	_, err := ParsePingFrame([]byte(`{"type":"HELLO","msg_id":1}`))
	if err == nil || !strings.Contains(err.Error(), ErrCodeInvalidFrame) {
		t.Fatalf("expected wrong type error")
	}
}

func TestParsePingFrameInvalidMsgID(t *testing.T) {
	_, err := ParsePingFrame([]byte(`{"type":"PING","msg_id":0}`))
	if err == nil {
		t.Fatalf("expected invalid msg_id error")
	}
}
