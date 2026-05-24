package hub

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func TestNegativeScenarios(t *testing.T) {
	// Inicjalizacja globalnej mapy sesji
	globalSessions = protection.NewSessionManager()

	// Przygotowanie certyfikatów
	err := os.MkdirAll("../certs", 0755)
	if err != nil {
		t.Fatalf("Nie udało się utworzyć folderu certs: %v", err)
	}
	tlsConf, err := certmanager.SetupTLSConfig("../certs/cert.pem", "../certs/key.pem")
	if err != nil {
		t.Fatalf("Błąd konfiguracji TLS: %v", err)
	}

	// Uruchomienie listenera
	listener, err := quic.ListenAddr("127.0.0.1:0", tlsConf, nil)
	if err != nil {
		t.Fatalf("Błąd uruchamiania listenera: %v", err)
	}
	defer listener.Close()

	// Nasłuchiwanie
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

		hello := protocol.HelloFrame{
			BaseFrame: protocol.BaseFrame{Type: "HELLO", MsgID: time.Now().UnixMilli()},
			Version:   "1.0",
		}
		hb, _ := json.Marshal(hello)
		stream.Write(hb)

		buf := make([]byte, 1024)
		stream.Read(buf)

		authFrame := protocol.AuthFrame{
			BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: time.Now().UnixMilli()},
			User:      "alice",
			Pass:      "wrongpassword",
		}
		ab, _ := json.Marshal(authFrame)
		stream.Write(ab)

		// Poprawka dla EOF
		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Błąd odczytu odpowiedzi: %v", err)
		}
		if n == 0 {
			t.Fatalf("Otrzymano EOF bez danych")
		}

		var errResp protocol.ErrorFrame
		json.Unmarshal(buf[:n], &errResp)

		if errResp.Code != "ERR_04" {
			t.Errorf("Oczekiwano ERR_04, otrzymano: %s", errResp.Code)
		}
	})

	t.Run("ERR_15_UserOffline", func(t *testing.T) {
		conn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
		if err != nil {
			t.Fatalf("Błąd połączenia (Dial): %v", err)
		}
		defer conn.CloseWithError(0, "")

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			t.Fatalf("Błąd otwarcia strumienia: %v", err)
		}

		hello := protocol.HelloFrame{
			BaseFrame: protocol.BaseFrame{Type: "HELLO", MsgID: time.Now().UnixMilli()},
			Version:   "1.0",
		}
		hb, _ := json.Marshal(hello)
		stream.Write(hb)
		buf := make([]byte, 1024)
		stream.Read(buf)

		authFrame := protocol.AuthFrame{
			BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: time.Now().UnixMilli()},
			User:      "alice",
			Pass:      "secret",
		}
		ab, _ := json.Marshal(authFrame)
		stream.Write(ab)
		stream.Read(buf)

		dataFrame := protocol.DataFrame{
			BaseFrame: protocol.BaseFrame{Type: "DATA", MsgID: time.Now().UnixMilli()},
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

		if errResp.Code != "ERR_15" {
			t.Errorf("Oczekiwano ERR_15, otrzymano: %s", errResp.Code)
		}
	})
}
