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

	// Channel for OS signals (SIGINT, SIGTERM, SIGQUIT).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Separate goroutine for signal handling – graceful client shutdown.
	go func() {
		sig := <-sigCh
		fmt.Println("\n[System] Caught signal:", sig)

		// If we have an active stream, try to send BYE.
		if stream != nil {
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
				// After successful AUTH we move to target selection.
				state = StateSelectingTarget
			} else {
				fmt.Printf("[System] Auth failed with code: %s. Reconnecting...\n", errCode)
				state = StateDisconnected
			}

		case StateSelectingTarget:
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

			// 2) Only now start ping loop + receiver loop.
			disconnected := make(chan struct{})
			startPingLoop(stream)
			pending, mu := startReceiverLoop(stream, disconnected)

			// After selecting the target we are fully "Established".
			state = StateEstablished

			// 3) Main sending loop.
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
						// Ctrl+D – treat as /quit.
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
					// Graceful termination: BYE + stream close.
					sendBye(stream)
					stream.Close()
					// We do NOT use CloseWithError – we do not want an application error on the Hub side.
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
