package app

import "testing"

// Test that long messages are truncated.
func TestTruncateMessageLong(t *testing.T) {
	input := make([]byte, MaxPlaintextSize+100)
	for i := range input {
		input[i] = 'A'
	}

	out := TruncateMessage(string(input))
	if len(out) != MaxPlaintextSize {
		t.Fatalf("expected truncated length %d, got %d", MaxPlaintextSize, len(out))
	}
}

// Test that short messages are unchanged.
func TestTruncateMessageShort(t *testing.T) {
	input := "hello"
	out := TruncateMessage(input)
	if out != input {
		t.Fatalf("expected unchanged message")
	}
}
