package hub

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"testing"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// TestHappyPathE2E verifies the full end‑to‑end flow of the KittyProtocol Hub:
// 1. HELLO → MEOW_OK
// 2. AUTH → MEOW_OK
// 3. DATA routing from Alice → Bob
// 4. ACK confirmation back to Alice
//
// This test runs against an isolated Hub instance created by StartTestHub(),
// ensuring no interference with global state or other tests.
func TestHappyPathE2E(t *testing.T) {
	// Start isolated Hub instance
	addr, stop, err := StartTestHub()
	if err != nil {
		t.Fatalf("Failed to start test Hub: %v", err)
	}
	defer stop()

	// Client QUIC config
	clientTLS := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"kitty-quic-v1"},
	}

	// ============================================================
	// 1. ALICE CONNECTS AND AUTHENTICATES
	// ============================================================

	aliceConn, err := quic.DialAddr(context.Background(), addr, clientTLS, nil)
	if err != nil {
		t.Fatalf("Alice connection failed: %v", err)
	}
	defer aliceConn.CloseWithError(0, "")

	aliceStream, err := aliceConn.OpenStreamSync(context.Background())
	if err != nil {
		t.Fatalf("Alice stream failed: %v", err)
	}

	bufA := make([]byte, 4096)

	// HELLO
	aliceHello := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeHello, MsgID: 1},
		Version:   "1.0",
	}
	b, _ := json.Marshal(aliceHello)
	aliceStream.Write(b)

	n, _ := aliceStream.Read(bufA)
	var helloAck protocol.MeowOkFrame
	json.Unmarshal(bufA[:n], &helloAck)

	if helloAck.Type != protocol.FrameTypeMeowOK {
		t.Fatalf("Alice expected MEOW_OK after HELLO, got: %s", helloAck.Type)
	}

	// AUTH
	aliceAuth := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 2},
		User:      "alice",
		Pass:      "secret",
	}
	b, _ = json.Marshal(aliceAuth)
	aliceStream.Write(b)

	n, _ = aliceStream.Read(bufA)
	var authAck protocol.MeowOkFrame
	json.Unmarshal(bufA[:n], &authAck)

	if authAck.Status != "Logged in" {
		t.Fatalf("Alice AUTH failed, status: %s", authAck.Status)
	}

	// ============================================================
	// 2. BOB CONNECTS AND AUTHENTICATES
	// ============================================================

	bobConn, err := quic.DialAddr(context.Background(), addr, clientTLS, nil)
	if err != nil {
		t.Fatalf("Bob connection failed: %v", err)
	}
	defer bobConn.CloseWithError(0, "")

	bobStream, err := bobConn.OpenStreamSync(context.Background())
	if err != nil {
		t.Fatalf("Bob stream failed: %v", err)
	}

	bufB := make([]byte, 4096)

	// HELLO
	bobHello := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeHello, MsgID: 3},
		Version:   "1.0",
	}
	b, _ = json.Marshal(bobHello)
	bobStream.Write(b)
	bobStream.Read(bufB)

	// AUTH
	bobAuth := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 4},
		User:      "bob",
		Pass:      "password",
	}
	b, _ = json.Marshal(bobAuth)
	bobStream.Write(b)
	bobStream.Read(bufB)

	// ============================================================
	// 3. ALICE SENDS DATA → BOB
	// ============================================================

	msgID := time.Now().UnixMilli()
	dataFrame := protocol.DataFrame{
		BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeData, MsgID: msgID},
		Target:    "bob",
		Payload:   "SGVsbG8gQm9iIQ==",
		MAC:       "dummy_mac_123",
	}
	b, _ = json.Marshal(dataFrame)
	aliceStream.Write(b)

	// Alice receives ACK
	n, _ = aliceStream.Read(bufA)
	var aliceAck protocol.MeowOkFrame
	json.Unmarshal(bufA[:n], &aliceAck)

	if aliceAck.Status != "Delivered" {
		t.Fatalf("Alice expected Delivered ACK, got: %s", aliceAck.Status)
	}

	// ============================================================
	// 4. BOB RECEIVES FORWARDED DATA
	// ============================================================

	n, err = bobStream.Read(bufB)
	if err != nil {
		t.Fatalf("Bob failed to read DATA: %v", err)
	}

	var bobRecv protocol.DataFrame
	json.Unmarshal(bufB[:n], &bobRecv)

	if bobRecv.Sender != "alice" {
		t.Fatalf("Bob expected Sender=alice, got: %s", bobRecv.Sender)
	}
	if bobRecv.Payload != "SGVsbG8gQm9iIQ==" {
		t.Fatalf("Bob received wrong payload")
	}
}
