package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/clientutils"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

func main() {
	// Konfiguracja TLS (InsecureSkipVerify dla testów lokalnych)
	tlsConf := &tls.Config{InsecureSkipVerify: true, NextProtos: []string{"kitty-quic"}}
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

	// 1. Logowanie (AUTH) - Task 6
	fmt.Print("Login: ")
	user, _ := reader.ReadString('\n')
	fmt.Print("Hasło: ")
	pass, _ := reader.ReadString('\n')

	auth := protocol.UniversalFrame{
		Type:  "AUTH",
		MsgID: time.Now().UnixMilli(),
		User:  user[:len(user)-1],
		Pass:  pass[:len(pass)-1],
	}
	stream.Write(auth.ToJSON())

	// Goroutine do odbierania wiadomości od innych przez Hub
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stream.Read(buf)
			if err != nil {
				return
			}
			frame, _ := protocol.ParseFrame(buf[:n])
			if frame != nil && frame.Type == "DATA" {
				fmt.Printf("\n[%s]: %s\n> ", frame.Sender, frame.Payload)
			}
		}
	}()

	fmt.Print("Do kogo piszesz?: ")
	target, _ := reader.ReadString('\n')
	target = target[:len(target)-1]

	// Pętla wysyłania - Task 14
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		// Task 14: Przycinanie wiadomości [cite: 664]
		safeText := clientutils.TruncateMessage(text)

		msgID := time.Now().UnixMilli()
		data := protocol.UniversalFrame{
			Type:    "DATA",
			MsgID:   msgID,
			Target:  target,
			Payload: safeText,
		}
		stream.Write(data.ToJSON())

		// Task 14: ACK Timer 5s [cite: 601, 664]
		go func(id int64) {
			time.Sleep(5 * time.Second)
			// Tu w Etapie 2 dojdzie sprawdzanie mapy potwierdzeń
			fmt.Printf("\n[Status] Wiadomość %d: Wysłano (oczekiwanie na ACK...)\n> ", id)
		}(msgID)
	}
}
