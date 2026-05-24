package protocol

import "testing"

func TestIsValidType(t *testing.T) {
	valid := []string{
		FrameTypeHello, FrameTypeAuth, FrameTypeRegister,
		FrameTypeData, FrameTypeMeowOK, FrameTypeError,
		FrameTypeGetStatus, FrameTypeStatusRes,
		FrameTypePing, FrameTypeBye,
	}

	for _, v := range valid {
		if !IsValidType(v) {
			t.Fatalf("expected valid type: %s", v)
		}
	}

	if IsValidType("HACK") {
		t.Fatalf("expected invalid type")
	}
}

func TestIsValidTypeEdgeCases(t *testing.T) {
	if IsValidType("") {
		t.Fatalf("empty string should be invalid")
	}
	if IsValidType("data") {
		t.Fatalf("lowercase should be invalid")
	}
	if IsValidType(" ") {
		t.Fatalf("whitespace should be invalid")
	}
}
