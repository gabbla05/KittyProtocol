package hub

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"testing"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// TestSecurityScenarios validates Hub‑level security mechanisms:
// 1. ERR_06 — replay attack detection
// 2. ERR_02 — malformed JSON injection.
//
// This test runs against an isolated Hub instance created by StartTestHub(),
// ensuring no interference with global state or other tests.
func TestSecurityScenarios(t *testing.T) {
	// Start isolated Hub instance
	addr, stop, err := StartTestHub()
	if err != nil {
		t.Fatalf("Failed to start test Hub: %v", err)
	}
	defer stop()

	clientTLS := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"kitty-quic-v1"},
	}

	// ============================================================
	// ERR_06 — Replay Attack Detection
	// ============================================================
	t.Run("ERR_06_ReplayAttack", func(t *testing.T) {
		conn, err := quic.DialAddr(context.Background(), addr, clientTLS, nil)
		if err != nil {
			t.Fatalf("Connection error: %v", err)
		}
		defer conn.CloseWithError(0, "")

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			t.Fatalf("Stream open error: %v", err)
		}

		buf := make([]byte, 1024)

		// HELLO
		hello := protocol.HelloFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeHello, MsgID: 1},
			Version:   "1.0",
		}
		hb, _ := json.Marshal(hello)
		if _, err := stream.Write(hb); err != nil {
			t.Fatalf("HELLO write error: %v", err)
		}
		if _, err := stream.Read(buf); err != nil && err != io.EOF {
			t.Fatalf("HELLO response read error: %v", err)
		}

		// AUTH
		authFrame := protocol.AuthFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 2},
			User:      "alice",
			Pass:      "secret",
		}
		ab, _ := json.Marshal(authFrame)
		if _, err := stream.Write(ab); err != nil {
			t.Fatalf("AUTH write error: %v", err)
		}
		if _, err := stream.Read(buf); err != nil && err != io.EOF {
			t.Fatalf("AUTH response read error: %v", err)
		}

		// DATA → Alice (self‑send to ensure routing succeeds)
		msgID := int64(100)
		dataFrame := protocol.DataFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeData, MsgID: msgID},
			Target:    "alice",
			Payload:   "SGVsbG8=",
			MAC:       "dummyMAC",
		}
		db, _ := json.Marshal(dataFrame)

		// First send — should be accepted
		if _, err := stream.Write(db); err != nil {
			t.Fatalf("First DATA write error: %v", err)
		}
		if _, err := stream.Read(buf); err != nil && err != io.EOF {
			t.Fatalf("First DATA response read error: %v", err)
		}

		// Replay — same MsgID
		if _, err := stream.Write(db); err != nil {
			t.Fatalf("Replay DATA write error: %v", err)
		}
		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Replay DATA response read error: %v", err)
		}
		if n == 0 {
			t.Fatalf("Received EOF without data on replay")
		}

		var errResp protocol.ErrorFrame
		if err := json.Unmarshal(buf[:n], &errResp); err != nil {
			t.Fatalf("Failed to unmarshal ERROR frame: %v", err)
		}

		if errResp.Code != protocol.ErrReplayDetected {
			t.Errorf("Expected ERR_06 (Replay detected), got: %s", errResp.Code)
		}
	})

	// ============================================================
	// ERR_02 — Malformed JSON Injection
	// ============================================================
	t.Run("ERR_02_Injection", func(t *testing.T) {
		conn, err := quic.DialAddr(context.Background(), addr, clientTLS, nil)
		if err != nil {
			t.Fatalf("Connection error: %v", err)
		}
		defer conn.CloseWithError(0, "")

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			t.Fatalf("Stream open error: %v", err)
		}

		// Malformed JSON payload (invalid syntax)
		badData := []byte(`{DROP TABLE users; HACK THE PLANET}`)
		if _, err := stream.Write(badData); err != nil {
			t.Fatalf("Malformed DATA write error: %v", err)
		}

		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Malformed DATA response read error: %v", err)
		}
		if n == 0 {
			t.Fatalf("Received EOF without data for malformed JSON")
		}

		var errResp protocol.ErrorFrame
		if err := json.Unmarshal(buf[:n], &errResp); err != nil {
			t.Fatalf("Failed to unmarshal ERROR frame: %v", err)
		}

		if errResp.Code != protocol.ErrFormatError {
			t.Errorf("Expected ERR_02 (Format error), got: %s", errResp.Code)
		}
	})
}
