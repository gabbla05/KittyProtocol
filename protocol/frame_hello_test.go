package protocol

import "testing"

func TestParseHelloFrameValid(t *testing.T) {
	json := []byte(`{"type":"HELLO","msg_id":1,"version":"1.0"}`)
	_, err := ParseHelloFrame(json)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseHelloFrameInvalid(t *testing.T) {
	_, err := ParseHelloFrame([]byte(`{"type":"HELLO"}`))
	if err == nil {
		t.Fatalf("expected missing version error")
	}

	_, err = ParseHelloFrame([]byte(`{"type":"HELLO","msg_id":1,"version":"9.9"}`))
	if err == nil {
		t.Fatalf("expected version mismatch error")
	}
}
