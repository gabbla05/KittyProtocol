package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

var msgCounter = 1

func main() {
	// TLS configuration (skips verification for self-signed test certs)
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"kitty-quic"},
	}

	fmt.Println("Connecting to Hub...")
	ctx := context.Background()
	conn, err := quic.DialAddr(ctx, "127.0.0.1:9999", tlsConf, nil)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// 1. Send HELLO
	helloMsg := protocol.UniversalFrame{Type: "HELLO", MsgID: msgCounter, Version: "1.0"}
	stream.Write(helloMsg.ToJSON())

	buf := make([]byte, 1024)
	stream.Read(buf) // Wait for MEOW_OK
	msgCounter++

	// 2. Prompt for credentials
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	login, _ := reader.ReadString('\n')
	login = login[:len(login)-1] // Remove trailing newline (\n)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	// 3. Send AUTH
	authMsg := protocol.UniversalFrame{Type: "AUTH", MsgID: msgCounter, User: login, Passw: password}
	stream.Write(authMsg.ToJSON())

	n, _ := stream.Read(buf)
	fmt.Println("\nAuthentication response from Hub:", string(buf[:n]))
	msgCounter++

	fmt.Println("Session initialized. Exiting (Stage 1 complete).")
}
