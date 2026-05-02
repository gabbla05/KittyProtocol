// client/main.go
package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/quic-go/quic-go"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"kitty-quic-v1"},
	}

	reader := bufio.NewReader(os.Stdin)
	state := StateDisconnected

	var conn *quic.Conn
	var stream *quic.Stream
	var err error

	// Kanał na sygnały systemowe (SIGINT, SIGTERM)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Osobna gorutyna do obsługi sygnałów – delikatne zamknięcie klienta.
	go func() {
		sig := <-sigCh
		fmt.Println("\n[System] Caught signal:", sig)

		// Jeśli mamy aktywny stream, spróbujmy wysłać BYE.
		if stream != nil {
			sendBye(stream)
			stream.Close()
		}
		// Conn nie musi być zamykany CloseWithError – idle timeout po stronie Huba
		// i tak go posprząta. Tutaj po prostu kończymy proces.
		fmt.Println("[System] Closing session due to signal.")
		os.Exit(0)
	}()

	for {
		switch state {
		case StateDisconnected:
			fmt.Println("[System] Connecting to Hub...")
			conn, err = quic.DialAddr(context.Background(), "127.0.0.1:9999", tlsConf, nil)
			if err != nil {
				fmt.Println("[System] Connection error:", err)
				time.Sleep(2 * time.Second)
				continue
			}
			stream, err = conn.OpenStreamSync(context.Background())
			if err != nil {
				fmt.Println("[System] Stream error:", err)
				state = StateDisconnected
				continue
			}
			sendHello(stream)
			state = StateHandshaking

		case StateHandshaking:
			success, errCode := waitForHelloOK(stream)
			if success {
				state = StateAuthenticating
			} else {
				if errCode == "ERR_03" || errCode == "CONNECTION_LOST" {
					fmt.Println("[System] Handshake failed. Reconnecting...")
					state = StateDisconnected
				} else {
					return
				}
			}

		case StateAuthenticating:
			user, pass := readCredentials(reader)
			sendAuth(stream, user, pass)
			success, errCode := waitForAuthOK(stream)
			if success {
				// Po poprawnym AUTH przechodzimy do wyboru adresata
				state = StateSelectingTarget
			} else {
				fmt.Printf("[System] Auth failed with code: %s. Reconnecting...\n", errCode)
				state = StateDisconnected
			}

		case StateSelectingTarget:
			fmt.Println("[System] Session established.")

			// 1) Najpierw wybieramy adresata
			var target string
			for {
				target = readTarget(reader) // użyjemy Twojej funkcji z input.go
				target = strings.TrimSpace(target)

				if len(target) == 0 {
					fmt.Println("[System] Target cannot be empty.")
					continue
				}

				break
			}

			// 2) Dopiero teraz startujemy ping + odbiornik
			disconnected := make(chan struct{})
			startPingLoop(stream)
			pending, mu := startReceiverLoop(stream, disconnected)

			// Po wybraniu targetu jesteśmy w pełni „Established”
			state = StateEstablished

			// 3) Główna pętla wysyłania
			for {
				select {
				case <-disconnected:
					fmt.Println("[System] Session closed. Returning to disconnected state.")
					state = StateDisconnected
					return
				default:
				}

				fmt.Print("> ")
				text, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						// Ctrl+D – traktujemy jak /quit
						fmt.Println("\n[System] EOF detected (Ctrl+D). Sending BYE and exiting.")
						if stream != nil {
							sendBye(stream)
							stream.Close()
						}
						return
					}
					fmt.Println("[System] Read error:", err)
					continue
				}

				if len(text) <= 1 {
					continue
				}
				text = strings.TrimSpace(text)

				if text == "/quit" {
					// Łagodne zakończenie: BYE + zamknięcie strumienia.
					sendBye(stream)
					stream.Close()
					// NIE używamy CloseWithError – nie chcemy Application error po stronie Huba.
					time.Sleep(200 * time.Millisecond)
					fmt.Println("[System] Closing session by user request.")
					state = StateDisconnected
					return
				}

				sendMessage(stream, target, text, pending, mu)
			}
		}
	}
}
