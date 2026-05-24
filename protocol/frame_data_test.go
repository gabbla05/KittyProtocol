package protocol

import (
	"strings"
	"testing"
)

func TestParseDataFrameValid(t *testing.T) {
	json := []byte(`{"type":"DATA","msg_id":1,"target":"bob","payload":"x","mac":"y"}`)
	_, err := ParseDataFrame(json)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseDataFrameMissingFields(t *testing.T) {
	json := []byte(`{"type":"DATA","msg_id":1,"target":"bob","payload":"x"}`)
	_, err := ParseDataFrame(json)
	if err == nil || !strings.Contains(err.Error(), ErrCodeInvalidFrame) {
		t.Fatalf("expected missing MAC error")
	}
}

func TestParseDataFrameWrongType(t *testing.T) {
	json := []byte(`{"type":"HELLO","msg_id":1}`)
	_, err := ParseDataFrame(json)
	if err == nil {
		t.Fatalf("expected wrong type error")
	}
}
