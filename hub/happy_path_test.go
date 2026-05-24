package hub

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// TestHappyPathE2E verifies the full end‑to‑end flow of the KittyProtocol Hub:
// 1. HELLO → MEOW_OK
// 2. AUTH → MEOW_OK
// 3. DATA routing between authenticated users
// 4. ACK delivery confirmation
//
// This test uses a mock authentication backend and an in‑memory QUIC listener.
// It does NOT test TLS correctness — only QUIC transport and protocol logic.
func TestHappyPathE2E(t *testing.T) {

	// Reset global state (Hub uses package‑level singletons)
	globalSessions = protection.NewSessionManager()
	globalAuth = auth.NewMockAuth() // mock DB with alice/secret, bob/password

	// Ensure certs directory exists (relative to project root)
	err := os.MkdirAll("../certs", 0755)
	if err != nil {
		t.Fatalf("Failed to create certs directory: %v", err)
	}

	// Load TLS certificates for QUIC
	tlsConf, err := certmanager.SetupTLSConfig("../certs/cert.pem", "../certs/key.pem")
	if err != nil {
		t.Fatalf("TLS setup failed: %v", err)
	}

	// Start Hub QUIC listener on random port
	listener, err := quic.ListenAddr("127.0.0.1:0", tlsConf, &quic.Config{
		KeepAlivePeriod: 1 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	// Hub accept loop (simulates Start(), but without TLS/DB/env)
	go func() {
		for {
			conn, err := listener.Accept(context.Background())
			if err != nil {
				return
			}
			go handleClient(conn)
		}
	}()

	// Client QUIC config
	clientTLS := &tls.Config{
		InsecureSkipVerify: true, // acceptable for local integration test
		NextProtos:         []string{"kitty-quic-v1"},
	}

	// ============================================================
	// 1. ALICE CONNECTS AND AUTHENTICATES
	// ============================================================

	aliceConn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
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

	bobConn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
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
