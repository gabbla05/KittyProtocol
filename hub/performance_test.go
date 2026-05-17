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

func BenchmarkHubRouting(b *testing.B) {
	cores := []int{1, 2, 4, 8, 16}
	maxCores := runtime.NumCPU()

	for _, c := range cores {
		if c > maxCores {
			continue
		}

		b.Run(fmt.Sprintf("Cores_%d", c), func(b *testing.B) {
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

			// Podłączenie Alice
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

			// Podłączenie Boba
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

			b.ResetTimer()
			startTime := time.Now()

			var wg sync.WaitGroup
			wg.Add(2)

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

			for i := 0; i < b.N; i++ {
				// POPRAWKA: Używamy nanosekund, aby MsgID było zawsze unikalne w trakcie benchmarku
				dataFrame.BaseFrame.MsgID = time.Now().UnixNano() + int64(i)
				mb, _ := json.Marshal(dataFrame)
				aliceStream.Write(mb)
			}

			wg.Wait()
			totalDuration := time.Since(startTime)

			saveToHistory(c, b.N, totalDuration)
		})
	}
}

func saveToHistory(cores int, totalOps int, duration time.Duration) {
	filename := "benchmark_history.md"

	_, err := os.Stat(filename)
	isNewFile := os.IsNotExist(err)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[Performance] Błąd zapisu do pliku logów: %v\n", err)
		return
	}
	defer f.Close()

	if isNewFile {
		f.WriteString("# KittyProtocol - Hub Routing Performance History\n\n")
		f.WriteString("| Execution Time | CPU Cores | Total Packets Routed | Combined Duration | Latency Per Packet | Throughput |\n")
		f.WriteString("|--- |--- |--- |--- |--- |--- |\n")
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	latencyNs := float64(duration.Nanoseconds()) / float64(totalOps)
	throughputMsgSec := float64(totalOps) / duration.Seconds()

	row := fmt.Sprintf("| %s | %d cores | %d | %s | %.2f ns/op | %.2f msg/s |\n",
		timestamp,
		cores,
		totalOps,
		duration.Round(time.Millisecond),
		latencyNs,
		throughputMsgSec,
	)

	f.WriteString(row)
}
