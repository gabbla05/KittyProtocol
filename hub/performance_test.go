package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// W testach wydajnościowych używamy BenchmarkXxx zamiast TestXxx
func BenchmarkHubRouting(b *testing.B) {
	// Definiujemy, na ilu rdzeniach chcemy testować serwer
	cores := []int{1, 2, 4, 8, 16}
	maxCores := runtime.NumCPU()

	for _, c := range cores {
		if c > maxCores {
			continue // Pomijamy testowanie na większej liczbie rdzeni, niż maszyna fizycznie posiada
		}

		// Uruchomienie sub-benchmarku dla danej liczby rdzeni
		b.Run(fmt.Sprintf("Cores_%d", c), func(b *testing.B) {
			// Kluczowe dla zadania: narzucenie Hubowi limitu wątków procesora!
			runtime.GOMAXPROCS(c)

			globalSessions = protection.NewSessionManager()

			_ = os.MkdirAll("../certs", 0755)
			tlsConf, _ := certmanager.SetupTLSConfig("../certs/cert.pem", "../certs/key.pem")

			listener, err := quic.ListenAddr("127.0.0.1:0", tlsConf, nil)
			if err != nil {
				b.Fatalf("listen error: %v", err)
			}
			defer listener.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				for {
					conn, err := listener.Accept(ctx)
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

			// ==========================================
			// Podłączenie Alice (Nadawca)
			// ==========================================
			aliceConn, _ := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
			defer aliceConn.CloseWithError(0, "")
			aliceStream, _ := aliceConn.OpenStreamSync(context.Background())

			authAlice := protocol.AuthFrame{
				BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: time.Now().UnixMilli()},
				User:      "alice",
				Pass:      "secret",
			}
			ab, _ := json.Marshal(authAlice)
			aliceStream.Write(ab)
			buf := make([]byte, 1024)
			aliceStream.Read(buf)

			// ==========================================
			// Podłączenie Boba (Odbiorca)
			// ==========================================
			bobConn, _ := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
			defer bobConn.CloseWithError(0, "")
			bobStream, _ := bobConn.OpenStreamSync(context.Background())

			authBob := protocol.AuthFrame{
				BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: time.Now().UnixMilli()},
				User:      "bob",
				Pass:      "secret",
			}
			bb, _ := json.Marshal(authBob)
			bobStream.Write(bb)
			bobStream.Read(buf)

			// ========================================================
			// NADPISANIE RATE LIMITERA NA POTRZEBY TESTU WYDAJNOŚCI
			// Dzięki temu nie musimy ingerować w kod produkcyjny!
			// ========================================================
			if aliceSess, ok := globalSessions.Get("alice"); ok {
				aliceSess.Limiter = protection.NewRateLimiter(9999999)
			}
			if bobSess, ok := globalSessions.Get("bob"); ok {
				bobSess.Limiter = protection.NewRateLimiter(9999999)
			}

			// Szablon wiadomości
			dataFrame := protocol.DataFrame{
				BaseFrame: protocol.BaseFrame{Type: "DATA", MsgID: 0},
				Target:    "bob",
				Payload:   "SGVsbG8gQm9iIQ==",
				MAC:       "dummy_mac",
			}

			// Zaczynamy mierzyć czas! (pomijamy czas łączenia QUIC i logowania)
			b.ResetTimer()

			var wg sync.WaitGroup
			wg.Add(2)

			// Konsument Boba (musi asynchronicznie odbierać zroutowane wiadomości z Huba)
			go func() {
				defer wg.Done()
				decoder := json.NewDecoder(bobStream)
				for i := 0; i < b.N; i++ {
					var dummy map[string]interface{}
					if err := decoder.Decode(&dummy); err != nil {
						break
					}
				}
			}()

			// Konsument Alice (musi odbierać potwierdzenia MEOW_OK, by nie zablokować bufora)
			go func() {
				defer wg.Done()
				decoder := json.NewDecoder(aliceStream)
				for i := 0; i < b.N; i++ {
					var dummy map[string]interface{}
					if err := decoder.Decode(&dummy); err != nil {
						break
					}
				}
			}()

			// Producent Alice (floodowanie Huba wiadomościami DATA)
			for i := 0; i < b.N; i++ {
				// WAŻNE: Aktualizacja MsgID, żeby nasz własny mechanizm ReplayProtection nie zablokował testu!
				dataFrame.BaseFrame.MsgID = int64(i + 1)
				mb, _ := json.Marshal(dataFrame)
				aliceStream.Write(mb)
			}

			// Czekamy, aż Bob odbierze wszystkie b.N wiadomości
			wg.Wait()
		})
	}
}
