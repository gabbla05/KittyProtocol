package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

var msgCounter = 1

func main() {
	// CLI flag for dynamic port assignment (żeby móc odpalić kilka kotów naraz)
	listenPort := flag.Int("port", 8888, "Port to listen for incoming P2P connections")
	flag.Parse()

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"kitty-quic"},
	}

	// 1. Uruchamiamy nasłuch P2P w tle
	go startP2PListener(*listenPort)

	// 2. Łączymy się z Hubem
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

	// HELLO
	helloMsg := protocol.UniversalFrame{Type: "HELLO", MsgID: msgCounter, Version: "1.0"}
	stream.Write(helloMsg.ToJSON())

	buf := make([]byte, 2048)
	stream.Read(buf) // Wait for MEOW_OK
	msgCounter++

	// AUTH
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	login, _ := reader.ReadString('\n')
	login = login[:len(login)-1]

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	// Wysyłamy AUTH razem z portem, na którym nasłuchuje nasz lokalny serwer P2P
	authMsg := protocol.UniversalFrame{Type: "AUTH", MsgID: msgCounter, User: login, Passw: password, Port: *listenPort}
	stream.Write(authMsg.ToJSON())

	n, _ := stream.Read(buf)
	fmt.Println("\n[Auth status]:", string(buf[:n]))
	msgCounter++

	// LOOKUP
	fmt.Print("Who do you want to chat with? (target username): ")
	target, _ := reader.ReadString('\n')
	target = target[:len(target)-1]

	lookupMsg := protocol.UniversalFrame{Type: "LOOKUP", MsgID: msgCounter, Target: target}
	stream.Write(lookupMsg.ToJSON())

	n, _ = stream.Read(buf)
	resp, err := protocol.ParseFrame(buf[:n])

	if err == nil && resp.Type == "MEOW_OK" {
		fmt.Printf("Found '%s' at %s:%d! Establishing P2P tunnel...\n", target, resp.IP, resp.Port)
		connectToPeer(resp.IP, resp.Port, tlsConf)
	} else {
		fmt.Println("Error:", string(buf[:n]))
	}
}

// startP2PListener uruchamia lokalny serwer QUIC do odbierania wiadomości
func startP2PListener(port int) {
	cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		log.Fatal("P2P Cert error:", err)
	}
	serverTLS := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"kitty-quic"},
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := quic.ListenAddr(addr, serverTLS, nil)
	if err != nil {
		log.Fatal("P2P Listen error:", err)
	}

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			continue
		}

		// Czyste domknięcie bez przekazywania argumentów
		go func() {
			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				return
			}
			buf := make([]byte, 2048)
			for {
				n, err := stream.Read(buf)
				if err != nil {
					return
				}
				frame, _ := protocol.ParseFrame(buf[:n])
				if frame != nil && frame.Type == "DATA" {
					fmt.Printf("\n[P2P Message received]: %s\n> ", frame.Payload)
				}
			}
		}()
	}
}

// connectToPeer otwiera strumień QUIC i rozpoczyna pętlę czatu konsolowego
func connectToPeer(ip string, port int, tlsConf *tls.Config) {
	ctx := context.Background()
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := quic.DialAddr(ctx, addr, tlsConf, nil)
	if err != nil {
		fmt.Println("Could not connect to peer:", err)
		return
	}
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		fmt.Println("Could not open stream:", err)
		return
	}

	// Wysyłamy P2P_HELLO
	helloMsg := protocol.UniversalFrame{Type: "P2P_HELLO", MsgID: msgCounter, Token: "session_token"}
	stream.Write(helloMsg.ToJSON())
	msgCounter++

	fmt.Println("🐈 You are now in a secure P2P tunnel! Type your messages.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		text := scanner.Text()

		if text == "/quit" {
			break
		}

		dataMsg := protocol.UniversalFrame{
			Type:    "DATA",
			MsgID:   msgCounter,
			Payload: text,
		}
		stream.Write(dataMsg.ToJSON())
		msgCounter++
	}
}
