package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func main() {
	tlsConf := &tls.Config{InsecureSkipVerify: true, NextProtos: []string{"kitty-quic"}}
	conn, _ := quic.DialAddr(context.Background(), "127.0.0.1:9999", tlsConf, nil)
	stream, _ := conn.OpenStreamSync(context.Background())

	reader := bufio.NewReader(os.Stdin)

	// 1. HELLO
	hello := protocol.UniversalFrame{Type: "HELLO", MsgID: time.Now().UnixMilli()}
	stream.Write(hello.ToJSON())

	// 2. AUTH
	fmt.Print("Login: ")
	user, _ := reader.ReadString('\n')
	user = user[:len(user)-1]
	auth := protocol.UniversalFrame{Type: "AUTH", MsgID: time.Now().UnixMilli(), User: user, Pass: "secret"}
	stream.Write(auth.ToJSON())

	// Wątek do odbierania wiadomości od innych (via Hub)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, _ := stream.Read(buf)
			frame, _ := protocol.ParseFrame(buf[:n])
			if frame.Type == "DATA" {
				fmt.Printf("\n[%s]: %s\n> ", frame.Sender, frame.Payload)
			}
		}
	}()

	// Pętla wysyłania
	fmt.Print("Target User: ")
	target, _ := reader.ReadString('\n')
	target = target[:len(target)-1]

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		// TODO (Gołson - Task 14): Tu dodaj ucinanie tekstu do 2048 bajtów [cite: 602]

		data := protocol.UniversalFrame{
			Type:    "DATA",
			MsgID:   time.Now().UnixMilli(),
			Target:  target,
			Payload: text, // W Etapie 2 tu będzie szyfrowanie MB [cite: 645]
		}
		stream.Write(data.ToJSON())
	}
}
