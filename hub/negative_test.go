package hub

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// TestNegativeScenarios verifies that the Hub correctly returns protocol‑level
// error frames (ERR_XX) for invalid authentication and invalid DATA routing.
// This test uses a mock authentication backend and an in‑memory QUIC listener.
func TestNegativeScenarios(t *testing.T) {

	// Reset global state
	globalSessions = protection.NewSessionManager()
	globalAuth = auth.NewMockAuth() // mock DB: alice/secret, bob/password

	// Prepare TLS certs (relative to project root)
	err := os.MkdirAll("../certs", 0755)
	if err != nil {
		t.Fatalf("Nie udało się utworzyć folderu certs: %v", err)
	}

	tlsConf, err := certmanager.SetupTLSConfig("../certs/cert.pem", "../certs/key.pem")
	if err != nil {
		t.Fatalf("Błąd konfiguracji TLS: %v", err)
	}

	// Start Hub listener
	listener, err := quic.ListenAddr("127.0.0.1:0", tlsConf, nil)
	if err != nil {
		t.Fatalf("Błąd uruchamiania listenera: %v", err)
	}
	defer listener.Close()

	// Accept loop
	go func() {
		for {
			conn, err := listener.Accept(context.Background())
			if err != nil {
				return
			}
			go handleClient(conn)
		}
	}()

	clientTLS := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"kitty-quic-v1"},
	}

	// ============================================================
	// ERR_04 — Authentication Failed (wrong password)
	// ============================================================
	t.Run("ERR_04_BadPassword", func(t *testing.T) {

		conn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
		if err != nil {
			t.Fatalf("Błąd połączenia (Dial): %v", err)
		}
		defer conn.CloseWithError(0, "")

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			t.Fatalf("Błąd otwarcia strumienia: %v", err)
		}

		// HELLO
		hello := protocol.HelloFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeHello, MsgID: 1},
			Version:   "1.0",
		}
		hb, _ := json.Marshal(hello)
		stream.Write(hb)

		buf := make([]byte, 1024)
		stream.Read(buf) // MEOW_OK

		// AUTH (wrong password)
		authFrame := protocol.AuthFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 2},
			User:      "alice",
			Pass:      "wrongpassword",
		}
		ab, _ := json.Marshal(authFrame)
		stream.Write(ab)

		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Błąd odczytu odpowiedzi: %v", err)
		}
		if n == 0 {
			t.Fatalf("Otrzymano EOF bez danych")
		}

		var errResp protocol.ErrorFrame
		json.Unmarshal(buf[:n], &errResp)

		if errResp.Code != protocol.ErrAuthenticationFailed {
			t.Errorf("Oczekiwano ERR_04, otrzymano: %s", errResp.Code)
		}
	})

	// ============================================================
	// ERR_15 — Unknown Target (user does not exist)
	// ============================================================
	t.Run("ERR_15_UnknownTarget", func(t *testing.T) {

		conn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
		if err != nil {
			t.Fatalf("Błąd połączenia (Dial): %v", err)
		}
		defer conn.CloseWithError(0, "")

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			t.Fatalf("Błąd otwarcia strumienia: %v", err)
		}

		buf := make([]byte, 1024)

		// HELLO
		hello := protocol.HelloFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeHello, MsgID: 1},
			Version:   "1.0",
		}
		hb, _ := json.Marshal(hello)
		stream.Write(hb)
		stream.Read(buf)

		// AUTH (correct)
		authFrame := protocol.AuthFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeAuth, MsgID: 2},
			User:      "alice",
			Pass:      "secret",
		}
		ab, _ := json.Marshal(authFrame)
		stream.Write(ab)
		stream.Read(buf)

		// DATA → ghostuser (does not exist)
		dataFrame := protocol.DataFrame{
			BaseFrame: protocol.BaseFrame{Type: protocol.FrameTypeData, MsgID: 3},
			Target:    "ghostuser",
			Payload:   "SGVsbG8=",
			MAC:       "dummyMAC",
		}
		db, _ := json.Marshal(dataFrame)
		stream.Write(db)

		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Błąd odczytu odpowiedzi: %v", err)
		}

		var errResp protocol.ErrorFrame
		json.Unmarshal(buf[:n], &errResp)

		if errResp.Code != protocol.ErrUnknownTarget {
			t.Errorf("Oczekiwano ERR_15, otrzymano: %s", errResp.Code)
		}
	})
}
