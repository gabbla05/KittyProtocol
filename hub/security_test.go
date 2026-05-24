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

func TestSecurityScenarios(t *testing.T) {
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

	t.Run("ERR_06_ReplayAttack", func(t *testing.T) {
		conn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
		if err != nil {
			t.Fatalf("Błąd połączenia: %v", err)
		}
		defer conn.CloseWithError(0, "")

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			t.Fatalf("Błąd strumienia: %v", err)
		}

		// 1. HELLO i AUTH
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
		stream.Read(buf) // Odbiór MEOW_OK

		// 2. Pierwsza ramka DATA
		msgID := time.Now().UnixMilli()
		dataFrame := protocol.DataFrame{
			BaseFrame: protocol.BaseFrame{Type: "DATA", MsgID: msgID},
			Target:    "bob",
			Payload:   "SGVsbG8=",
			MAC:       "dummyMAC",
		}
		db, _ := json.Marshal(dataFrame)
		stream.Write(db)
		stream.Read(buf) // Huba odpowiada (MEOW_OK lub ERR_15) i zapisuje MsgID w cache'u anty-replay

		// 3. Replay (ponowne wysłanie tej samej ramki z tym samym MsgID)
		stream.Write(db)
		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Błąd odczytu: %v", err)
		}

		// 4. Weryfikacja
		var errResp protocol.ErrorFrame
		json.Unmarshal(buf[:n], &errResp)

		if errResp.Code != "ERR_06" {
			t.Errorf("Oczekiwano ERR_06 (Replay detected), otrzymano: %s", errResp.Code)
		}
	})

	t.Run("ERR_02_Injection", func(t *testing.T) {
		conn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
		if err != nil {
			t.Fatalf("Błąd połączenia: %v", err)
		}
		defer conn.CloseWithError(0, "")

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			t.Fatalf("Błąd strumienia: %v", err)
		}

		// 1. Wstrzyknięcie złośliwych danych zamiast struktury JSON
		badData := []byte(`{DROP TABLE users; HACK THE PLANET}`)
		stream.Write(badData)

		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Błąd odczytu: %v", err)
		}

		// 2. Weryfikacja błędu z parsera JSON
		var errResp protocol.ErrorFrame
		json.Unmarshal(buf[:n], &errResp)

		if errResp.Code != "ERR_02" {
			t.Errorf("Oczekiwano ERR_02 (błąd formatu), otrzymano: %s", errResp.Code)
		}
	})
}
