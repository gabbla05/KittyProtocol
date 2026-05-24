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

// TestNegativeScenarios verifies that the Hub correctly returns protocol‑level
// error frames (ERR_XX) for invalid authentication and invalid DATA routing.
//
// This test runs against an isolated Hub instance created by StartTestHub(),
// ensuring no interference with global state or other tests.
func TestNegativeScenarios(t *testing.T) {
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
	// ERR_04 — Authentication Failed (wrong password)
	// ============================================================
	t.Run("ERR_04_BadPassword", func(t *testing.T) {
		conn, err := quic.DialAddr(context.Background(), addr, clientTLS, nil)
		if err != nil {
			t.Fatalf("Connection error (Dial): %v", err)
		}
		defer conn.CloseWithError(0, "")

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			t.Fatalf("Stream open error: %v", err)
		}

		// HELLO
		hello := protocol.HelloFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeHello, MsgID: 1},
			Version:   "1.0",
		}
		hb, _ := json.Marshal(hello)
		if _, err := stream.Write(hb); err != nil {
			t.Fatalf("HELLO write error: %v", err)
		}

		buf := make([]byte, 1024)
		if _, err := stream.Read(buf); err != nil && err != io.EOF {
			t.Fatalf("HELLO response read error: %v", err)
		}

		// AUTH (wrong password)
		authFrame := protocol.AuthFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 2},
			User:      "alice",
			Pass:      "wrongpassword",
		}
		ab, _ := json.Marshal(authFrame)
		if _, err := stream.Write(ab); err != nil {
			t.Fatalf("AUTH write error: %v", err)
		}

		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("AUTH response read error: %v", err)
		}
		if n == 0 {
			t.Fatalf("Received EOF without data")
		}

		var errResp protocol.ErrorFrame
		if err := json.Unmarshal(buf[:n], &errResp); err != nil {
			t.Fatalf("Failed to unmarshal ERROR frame: %v", err)
		}

		if errResp.Code != protocol.ErrAuthenticationFailed {
			t.Errorf("Expected ERR_04, got: %s", errResp.Code)
		}
	})

	// ============================================================
	// ERR_15 — Unknown Target (user does not exist)
	// ============================================================
	t.Run("ERR_15_UnknownTarget", func(t *testing.T) {
		conn, err := quic.DialAddr(context.Background(), addr, clientTLS, nil)
		if err != nil {
			t.Fatalf("Connection error (Dial): %v", err)
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

		// AUTH (correct)
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

		// DATA → ghostuser (does not exist)
		dataFrame := protocol.DataFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeData, MsgID: 3},
			Target:    "ghostuser",
			Payload:   "SGVsbG8=",
			MAC:       "dummyMAC",
		}
		db, _ := json.Marshal(dataFrame)
		if _, err := stream.Write(db); err != nil {
			t.Fatalf("DATA write error: %v", err)
		}

		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("DATA response read error: %v", err)
		}
		if n == 0 {
			t.Fatalf("Received EOF without data")
		}

		var errResp protocol.ErrorFrame
		if err := json.Unmarshal(buf[:n], &errResp); err != nil {
			t.Fatalf("Failed to unmarshal ERROR frame: %v", err)
		}

		if errResp.Code != protocol.ErrUnknownTarget {
			t.Errorf("Expected ERR_15, got: %s", errResp.Code)
		}
	})
}
