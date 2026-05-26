package protocol

import (
	"encoding/json"
	"testing"
)

func TestMeowOkFrameJSON(t *testing.T) {
	data := []byte(`{"type":"MEOW_OK","msg_id":1,"status":"ok"}`)
	var f MeowOkFrame
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
