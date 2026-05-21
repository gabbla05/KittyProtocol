package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func TestHappyPathE2E(t *testing.T) {
	// Inicjalizacja globalnej mapy sesji (Task 5 / Task 9)
	globalSessions = protection.NewSessionManager()

	// Przygotowanie certyfikatów TLS (Task 32)
	err := os.MkdirAll("../certs", 0755)
	if err != nil {
		t.Fatalf("Nie udało się utworzyć folderu certs: %v", err)
	}
	tlsConf, err := certmanager.SetupTLSConfig("../certs/cert.pem", "../certs/key.pem")
	if err != nil {
		t.Fatalf("Błąd konfiguracji TLS: %v", err)
	}

	// Uruchomienie listenera Huba
	listener, err := quic.ListenAddr("127.0.0.1:0", tlsConf, nil)
	if err != nil {
		t.Fatalf("Błąd uruchamiania listenera: %v", err)
	}
	defer listener.Close()

	// Nasłuchiwanie w tle (Task 4)
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
		InsecureSkipVerify: true, // Akceptowalne dla lokalnego testu integracyjnego
		NextProtos:         []string{"kitty-quic-v1"},
	}

	// ==========================================
	// KROK 1: Podłączenie i logowanie Alice
	// ==========================================
	aliceConn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
	if err != nil {
		t.Fatalf("Błąd połączenia Alice: %v", err)
	}
	defer aliceConn.CloseWithError(0, "")
	aliceStream, err := aliceConn.OpenStreamSync(context.Background())
	if err != nil {
		t.Fatalf("Błąd strumienia Alice: %v", err)
	}

	// Alice: HELLO
	aliceHello := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{Type: "HELLO", MsgID: time.Now().UnixMilli()},
		Version:   "1.0",
	}
	b, _ := json.Marshal(aliceHello)
	aliceStream.Write(b)
	bufA := make([]byte, 2048)
	aliceStream.Read(bufA) // Odbiór MEOW_OK

	// Alice: AUTH
	aliceAuth := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: time.Now().UnixMilli()},
		User:      "alice",
		Pass:      "secret", // Hasło z mock DB
	}
	b, _ = json.Marshal(aliceAuth)
	aliceStream.Write(b)
	aliceStream.Read(bufA) // Odbiór MEOW_OK("Logged in")

	// ==========================================
	// KROK 2: Podłączenie i logowanie Boba
	// ==========================================
	bobConn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
	if err != nil {
		t.Fatalf("Błąd połączenia Boba: %v", err)
	}
	defer bobConn.CloseWithError(0, "")
	bobStream, err := bobConn.OpenStreamSync(context.Background())
	if err != nil {
		t.Fatalf("Błąd strumienia Boba: %v", err)
	}

	// Bob: HELLO
	bobHello := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{Type: "HELLO", MsgID: time.Now().UnixMilli()},
		Version:   "1.0",
	}
	b, _ = json.Marshal(bobHello)
	bobStream.Write(b)
	bufB := make([]byte, 2048)
	bobStream.Read(bufB)

	// Bob: AUTH
	bobAuth := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: time.Now().UnixMilli()},
		User:      "bob",
		Pass:      "password", // Hasło z mock DB
	}
	b, _ = json.Marshal(bobAuth)
	bobStream.Write(b)
	bobStream.Read(bufB) // Odbiór MEOW_OK("Logged in")

	// ==========================================
	// KROK 3: Alice wysyła DATA do Boba
	// ==========================================
	msgID := time.Now().UnixMilli()
	dataFrame := protocol.DataFrame{
		BaseFrame: protocol.BaseFrame{Type: "DATA", MsgID: msgID},
		Target:    "bob",
		Payload:   "SGVsbG8gQm9iIQ==", // Zakodowane Base64 (np. z E2EE)
		MAC:       "dummy_mac_123",
	}
	b, _ = json.Marshal(dataFrame)
	aliceStream.Write(b)

	// Alice powinna dostać MEOW_OK ("Delivered (mock)") od Huba
	nA, _ := aliceStream.Read(bufA)
	var aliceAck protocol.MeowOkFrame
	json.Unmarshal(bufA[:nA], &aliceAck)

	if aliceAck.Type != "MEOW_OK" {
		t.Errorf("Alice oczekiwała MEOW_OK, otrzymała typ: %s", aliceAck.Type)
	}

	// ==========================================
	// KROK 4: Bob odbiera zroutowaną wiadomość
	// ==========================================
	nB, err := bobStream.Read(bufB)
	if err != nil {
		t.Fatalf("Bob nie odczytał zroutowanej wiadomości: %v", err)
	}

	var bobRecv protocol.DataFrame
	json.Unmarshal(bufB[:nB], &bobRecv)

	// Hub podczas routingu powinien dokleić pole Sender (Task 7)
	if bobRecv.Type != "DATA" || bobRecv.Sender != "alice" || bobRecv.Payload != "SGVsbG8gQm9iIQ==" {
		t.Errorf("Bob otrzymał niepoprawną ramkę: %s", string(bufB[:nB]))
	}
}
