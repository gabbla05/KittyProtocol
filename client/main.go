package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/quic-go/quic-go"
)

// main initializes the QUIC client, performs HELLO → AUTH,
// starts background goroutines, and enters the interactive message loop.
func main() {
    tlsConf := &tls.Config{
        InsecureSkipVerify: true, // for local testing only
        NextProtos:         []string{"kitty-quic"},
    }

    conn, err := quic.DialAddr(context.Background(), "127.0.0.1:9999", tlsConf, nil)
    if err != nil {
        fmt.Println("Connection error:", err)
        return
    }
    stream, err := conn.OpenStreamSync(context.Background())
    if err != nil {
        fmt.Println("Stream error:", err)
        return
    }
    defer stream.Close()

    reader := bufio.NewReader(os.Stdin)

    // HELLO phase
    sendHello(stream)
    if !waitForHelloOK(stream) {
        return
    }

    // AUTH phase
    user, pass := readCredentials(reader)
    sendAuth(stream, user, pass)
    if !waitForAuthOK(stream) {
        return
    }

    // Background tasks
    startPingLoop(stream)
    pending, mu := startReceiverLoop(stream)

    // Messaging
    target := readTarget(reader)
    for {
        fmt.Print("> ")
        text, _ := reader.ReadString('\n')
        if len(text) == 0 {
            continue
        }
        text = text[:len(text)-1]
        if text == "" {
            continue
        }
        sendMessage(stream, target, text, pending, mu)
    }
}
