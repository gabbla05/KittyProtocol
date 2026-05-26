package hub

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

var globalMsgID int64 = time.Now().UnixNano()

func BenchmarkHubRouting(b *testing.B) {
	cores := []int{1, 2, 4, 8, 16}
	maxCores := runtime.NumCPU()

	hostname, _ := os.Hostname()

	for _, c := range cores {
		if c > maxCores {
			continue
		}

		b.Run(fmt.Sprintf("Cores_%d", c), func(b *testing.B) {
			old := runtime.GOMAXPROCS(c)
			defer runtime.GOMAXPROCS(old)

			if globalSessions == nil {
				globalSessions = protection.NewSessionManager()
			}

			_ = os.MkdirAll("../certs", 0755)
			tlsConf, _ := certmanager.SetupTLSConfig("../certs/cert.pem", "../certs/key.pem")

			listener, err := quic.ListenAddr("127.0.0.1:0", tlsConf, &quic.Config{
				MaxIdleTimeout: 2 * time.Second,
			})
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

			// --- PARAMETRY OBCIĄŻENIA ---
			numClients := 50     // liczba par Alice/Bob
			msgsPerClient := 200 // ile wiadomości wysyła jedna Alice
			totalOps := numClients * msgsPerClient

			b.ResetTimer()
			start := time.Now()

			var wg sync.WaitGroup
			wg.Add(numClients)

			for clientID := 0; clientID < numClients; clientID++ {
				go func(id int) {
					defer wg.Done()

					userAlice := fmt.Sprintf("alice_%d", id)
					userBob := fmt.Sprintf("bob_%d", id)

					// wyczyść ewentualne stare sesje
					globalSessions.Remove(userAlice)
					globalSessions.Remove(userBob)

					// ALICE
					aliceConn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
					if err != nil {
						return
					}
					defer aliceConn.CloseWithError(0, "")
					aliceStream, err := aliceConn.OpenStreamSync(context.Background())
					if err != nil {
						return
					}

					atomic.AddInt64(&globalMsgID, 1)
					authAlice := protocol.AuthFrame{
						BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: globalMsgID},
						User:      userAlice,
						Pass:      "secret",
					}
					ab, _ := json.Marshal(authAlice)
					aliceStream.Write(append(ab, '\n'))

					buf := make([]byte, 1024)
					aliceStream.Read(buf)

					// BOB
					bobConn, err := quic.DialAddr(context.Background(), listener.Addr().String(), clientTLS, nil)
					if err != nil {
						return
					}
					defer bobConn.CloseWithError(0, "")
					bobStream, err := bobConn.OpenStreamSync(context.Background())
					if err != nil {
						return
					}

					atomic.AddInt64(&globalMsgID, 1)
					authBob := protocol.AuthFrame{
						BaseFrame: protocol.BaseFrame{Type: "AUTH", MsgID: globalMsgID},
						User:      userBob,
						Pass:      "password",
					}
					bb, _ := json.Marshal(authBob)
					bobStream.Write(append(bb, '\n'))
					bobStream.Read(buf)

					// wyłącz rate limiting
					if s, ok := globalSessions.Get(userAlice); ok {
						s.Limiter = protection.NewRateLimiter(9999999)
					}
					if s, ok := globalSessions.Get(userBob); ok {
						s.Limiter = protection.NewRateLimiter(9999999)
					}

					dataFrame := protocol.DataFrame{
						BaseFrame: protocol.BaseFrame{Type: "DATA"},
						Target:    userBob,
						Payload:   "SGVsbG8gQm9iIQ==",
						MAC:       "dummy_mac",
					}

					// consumer BOB
					go func() {
						tmp := make([]byte, 4096)
						for i := 0; i < msgsPerClient; i++ {
							if _, err := bobStream.Read(tmp); err != nil {
								return
							}
						}
					}()

					// producer ALICE
					for i := 0; i < msgsPerClient; i++ {
						atomic.AddInt64(&globalMsgID, 1)
						dataFrame.MsgID = globalMsgID
						mb, _ := json.Marshal(dataFrame)
						if _, err := aliceStream.Write(append(mb, '\n')); err != nil {
							return
						}
					}
				}(clientID)
			}

			wg.Wait()
			total := time.Since(start)

			saveToHistory(hostname, maxCores, c, totalOps, total)
		})
	}
}

func saveToHistory(hostname string, maxCores int, testCores int, totalOps int, duration time.Duration) {
	filename := "../markdowns/benchmark_history.md"

	// sprawdź, czy plik już istnieje
	_, err := os.Stat(filename)
	isNewFile := os.IsNotExist(err)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// w benchmarku nie robimy panic/fatal – po prostu pomijamy zapis
		return
	}
	defer f.Close()

	// jeśli nowy plik – nagłówek + header tabeli
	if isNewFile {
		f.WriteString("# KittyProtocol - Hub Routing Performance History\n\n")
		f.WriteString("| Date | PC Name | Max Cores | Used Cores | Packets | Duration | Latency/pkt | Throughput |\n")
		f.WriteString("|--- |--- |--- |--- |--- |--- |--- |--- |\n")
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	latencyNs := float64(duration.Nanoseconds()) / float64(totalOps)
	throughputMsgSec := float64(totalOps) / duration.Seconds()

	row := fmt.Sprintf(
		"| %s | %s | %d | %d | %d | %s | %.2f ns | %.2f msg/s |\n",
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
