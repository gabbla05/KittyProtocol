package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/certmanager"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// GLOBALNY LICZNIK PAKIETÓW - Gwarantuje brak kolizji z Anti-Replay
var globalMsgID int64 = time.Now().UnixNano()

func BenchmarkHubRouting(b *testing.B) {
	cores := []int{1, 2, 4, 8, 16}
	maxCores := runtime.NumCPU()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown-PC"
	}

	for _, c := range cores {
		if c > maxCores {
			continue
		}

		b.Run(fmt.Sprintf("Cores_%d", c), func(b *testing.B) {
			// Zapisujemy poprzednie ustawienia i przywracamy je po teście,
			// żeby framework testowy Go nie rzucał błędem "left GOMAXPROCS".
			oldProcs := runtime.GOMAXPROCS(c)
			defer runtime.GOMAXPROCS(oldProcs)

			// 1. Czysty stan sesji dla każdej fazy benchmarku
			if globalSessions != nil {
				globalSessions.Remove("alice")
				globalSessions.Remove("bob")
			} else {
				globalSessions = protection.NewSessionManager()
			}

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

			// 2. PODŁĄCZENIE ALICE
			aliceConn, _ := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
			defer aliceConn.CloseWithError(0, "")
			aliceStream, _ := aliceConn.OpenStreamSync(context.Background())

			atomic.AddInt64(&globalMsgID, 1)
			authAlice := protocol.AuthFrame{
				BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: globalMsgID},
				User:      "alice",
				Pass:      "secret",
			}
			ab, _ := json.Marshal(authAlice)
			aliceStream.Write(append(ab, '\n'))
			buf := make([]byte, 1024)
			aliceStream.Read(buf)

			// 3. PODŁĄCZENIE BOBA
			bobConn, _ := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
			defer bobConn.CloseWithError(0, "")
			bobStream, _ := bobConn.OpenStreamSync(context.Background())

			atomic.AddInt64(&globalMsgID, 1)
			authBob := protocol.AuthFrame{
				BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: globalMsgID},
				User:      "bob",
				Pass:      "secret",
			}
			bb, _ := json.Marshal(authBob)
			bobStream.Write(append(bb, '\n'))
			bobStream.Read(buf)

			// 4. WYŁĄCZENIE LIMITÓW RUCHU (Rate Limiter)
			if aliceSess, ok := globalSessions.Get("alice"); ok {
				aliceSess.Limiter = protection.NewRateLimiter(9999999)
			}
			if bobSess, ok := globalSessions.Get("bob"); ok {
				bobSess.Limiter = protection.NewRateLimiter(9999999)
			}

			dataFrame := protocol.DataFrame{
				BaseFrame: protocol.BaseFrame{Type: "DATA", MsgID: 0},
				Target:    "bob",
				Payload:   "SGVsbG8gQm9iIQ==",
				MAC:       "dummy_mac",
			}

			// Bezpieczny timeout (10 sekund), zapobiega wiecznemu wiszeniu testu
			bobStream.SetReadDeadline(time.Now().Add(10 * time.Second))

			b.ResetTimer()
			startTime := time.Now()

			var wg sync.WaitGroup
			wg.Add(1)
			var fatalError bool

			// 5A. KONSUMENT BOBA
			go func() {
				defer wg.Done()
				decoder := json.NewDecoder(bobStream)
				for i := 0; i < b.N; i++ {
					var resp map[string]interface{}
					if err := decoder.Decode(&resp); err != nil {
						b.Errorf("\n[BOB] Serwer przerwał routing na pakiecie %d. Błąd: %v", i, err)
						fatalError = true
						return
					}
					if t, ok := resp["type"].(string); ok && t == "ERROR" {
						b.Errorf("\n[BOB] Dostał błąd od Huba: %v", resp)
						fatalError = true
						return
					}
				}
			}()

			// 5B. KONSUMENT ALICE (Czyści potwierdzenia w tle)
			go func() {
				decoder := json.NewDecoder(aliceStream)
				for {
					var dummy map[string]interface{}
					if err := decoder.Decode(&dummy); err != nil {
						return
					}
				}
			}()

			// 5C. PRODUCENT ALICE
			for i := 0; i < b.N; i++ {
				if fatalError {
					break
				}

				// Używamy bezpiecznego, rosnącego o 1 licznika. Nigdy nie wygeneruje kolizji.
				atomic.AddInt64(&globalMsgID, 1)
				dataFrame.BaseFrame.MsgID = globalMsgID

				mb, _ := json.Marshal(dataFrame)
				if _, err := aliceStream.Write(append(mb, '\n')); err != nil {
					break
				}

				// 100 mikrosekund pauzy zapobiega zapchaniu rur sieciowych
				time.Sleep(100 * time.Microsecond)
			}

			wg.Wait()
			totalDuration := time.Since(startTime)

			if !fatalError {
				saveToHistory(hostname, maxCores, c, b.N, totalDuration)
			}
		})
	}
}

func saveToHistory(hostname string, maxCores int, testCores int, totalOps int, duration time.Duration) {
	filename := "benchmark_history.md"
	_, err := os.Stat(filename)
	isNewFile := os.IsNotExist(err)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	if isNewFile {
		f.WriteString("# KittyProtocol - Hub Routing Performance History\n\n")
		f.WriteString("| Date | PC Name | Max Cores | Used Cores | Packets | Duration | Latency/pkt | Throughput |\n")
		f.WriteString("|--- |--- |--- |--- |--- |--- |--- |--- |\n")
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	latencyNs := float64(duration.Nanoseconds()) / float64(totalOps)
	throughputMsgSec := float64(totalOps) / duration.Seconds()

	row := fmt.Sprintf("| %s | %s | %d | %d | %d | %s | %.2f ns | %.2f msg/s |\n",
		timestamp,
		hostname,
		maxCores,
		testCores,
		totalOps,
		duration.Round(time.Millisecond),
		latencyNs,
		throughputMsgSec,
	)
	f.WriteString(row)
}
