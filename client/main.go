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

// Main entry point for the KittyProtocol client (Meowssenger).
// Responsibilities:
//   - establish QUIC + TLS 1.3 connection to the Hub,
//   - perform application-level handshake (HELLO / HELLO_OK),
//   - authenticate user (AUTH),
//   - select target user for the chat,
//   - run the main send loop for plaintext input,
//   - gracefully close the session (BYE + stream close) on exit.
func main() {
	// Minimal TLS configuration for development.
	// In production: proper certificate verification, pinning, etc.
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"kitty-quic-v1"},
	}

	reader := bufio.NewReader(os.Stdin)
	state := StateDisconnected

	var conn *quic.Conn
	var stream *quic.Stream
	var err error

	// OS signal channel (SIGINT, SIGTERM, SIGQUIT) – allows graceful shutdown
	// instead of abruptly killing the process.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Separate goroutine for signal handling – does not block the main loop.
	go func() {
		sig := <-sigCh
		fmt.Println("\n[System] Caught signal:", sig)

		// If we have an active stream, try to send BYE.
		if stream != nil {
			// Try to send BYE before closing the stream.
			sendBye(stream)
			stream.Close()
		}
		// Connection does not need CloseWithError – Hub's idle timeout
		// will clean it up anyway. We simply exit the process.
		fmt.Println("[System] Closing session due to signal.")
		os.Exit(0)
	}()

	for {
		switch state {
		case StateDisconnected:
			// No active connection – try to establish a new QUIC session to the Hub.
			fmt.Println("[System] Connecting to Hub...")
			conn, err = quic.DialAddr(context.Background(), "127.0.0.1:9999", tlsConf, nil)
			if err != nil {
				fmt.Println("[System] Connection error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			// Open a bidirectional stream – all KittyProtocol frames go over this stream.
			stream, err = conn.OpenStreamSync(context.Background())
			if err != nil {
				fmt.Println("[System] Stream error:", err)
				state = StateDisconnected
				continue
			}

			// HELLO – start of the application-level protocol.
			sendHello(stream)
			state = StateHandshaking

		case StateHandshaking:
			// Wait for HELLO_OK or an error (e.g. ERR_03).
			success, errCode := waitForHelloOK(stream)
			if success {
				state = StateAuthenticating
			} else {
				if errCode == "ERR_03" || errCode == "CONNECTION_LOST" {
					fmt.Println("[System] Handshake failed. Reconnecting...")
					state = StateDisconnected
				} else {
					// Other error – terminate the client.
					return
				}
			}

		case StateAuthenticating:
			// AUTH phase – read username/password from stdin.
			user, pass := readCredentials(reader)
			sendAuth(stream, user, pass)

			// Wait for AUTH_OK or ERR_xx.
			success, errCode := waitForAuthOK(stream)
			if success {
				// After successful AUTH we move to target selection.
				state = StateSelectingTarget
			} else {
				fmt.Printf("[System] Auth failed with code: %s. Reconnecting...\n", errCode)
				state = StateDisconnected
			}

		case StateSelectingTarget:
			// After successful AUTH – select the target user for the chat.
			fmt.Println("[System] Session established.")

			// 1) First, select the recipient.
			var target string
			for {
				target = readTarget(reader)
				target = strings.TrimSpace(target)

				if len(target) == 0 {
					fmt.Println("[System] Target cannot be empty.")
					continue
				}

				break
			}

			// Channel used to signal that the session/stream has been closed.
			disconnected := make(chan struct{})

			// Periodic PING loop – keeps the session alive on the Hub side.
			startPingLoop(stream)

			// Receiver loop – handles MEOW_OK, ERROR and DATA (with E2EE + replay protection).
			pending, mu := startReceiverLoop(stream, disconnected)

			// After selecting the target we are fully "Established".
			state = StateEstablished

			// Main sending loop – reads user input and sends DATA frames.
			for {
				select {
				case <-disconnected:
					// Hub closed the connection or an error occurred on the stream.
					fmt.Println("[System] Session closed. Returning to disconnected state.")
					state = StateDisconnected
					return
				default:
				}

				fmt.Print("> ")
				text, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						// User closed stdin (Ctrl+D) – send BYE and exit.
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

				text = strings.TrimSpace(text)
				if len(text) == 0 {
					continue
				}
				text = strings.TrimSpace(text)

				// Local command to terminate the session.
				if text == "/quit" {
					// Graceful termination: BYE + stream close.
					sendBye(stream)
					stream.Close()
					// We do NOT use CloseWithError – we do not want an application error on the Hub side.
					time.Sleep(200 * time.Millisecond)
					fmt.Println("[System] Closing session by user request.")
					state = StateDisconnected
					return
				}

				if after, ok := strings.CutPrefix(text, "/status "); ok {
					targetUser := strings.TrimSpace(after)
					if targetUser == "" {
						fmt.Println("[System] Usage: /status <username>")
						continue
					}
					sendGetStatus(stream, targetUser)
					continue
				}

				// Production message sending:
				//   - generates msg_id (timestamp),
				//   - performs E2EE (AEAD + HMAC) in cryptoee,
				//   - sends DATA frame,
				//   - starts ACK timer and tracks delivery in the pending map.
				sendMessage(stream, target, text, pending, mu)
			}
		}
	}
}
