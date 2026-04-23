package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/quic-go/quic-go"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"kitty-quic"},
	}

	reader := bufio.NewReader(os.Stdin)
	state := StateDisconnected

	var conn *quic.Conn
	var stream *quic.Stream
	var err error

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
				state = StateEstablished
			} else {
				fmt.Printf("[System] Auth failed with code: %s. Reconnecting...\n", errCode)
				state = StateDisconnected
			}

		case StateEstablished:
			fmt.Println("[System] Session established.")
			startPingLoop(stream)
			pending, mu := startReceiverLoop(stream)
			target := readTarget(reader)
			for {
				fmt.Print("> ")
				text, _ := reader.ReadString('\n')
				if len(text) <= 1 {
					continue
				}
				text = text[:len(text)-1]
				sendMessage(stream, target, text, pending, mu)
			}
		}
	}
}
