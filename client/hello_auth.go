package main

import (
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// sendHello sends the initial HELLO frame to the Hub.
func sendHello(stream *quic.Stream) {
    hello := protocol.UniversalFrame{
        Type:  "HELLO",
        MsgID: time.Now().UnixMilli(),
    }
    stream.Write(hello.ToJSON())
}

// waitForHelloOK waits for MEOW_OK after HELLO and prints the server status.
func waitForHelloOK(stream *quic.Stream) bool {
    buf := make([]byte, 4096)
    n, err := stream.Read(buf)
    if err != nil {
        fmt.Println("Read error after HELLO:", err)
        return false
    }
    frame, err := protocol.ParseFrame(buf[:n])
    if err != nil {
        fmt.Println("Parse error after HELLO:", err)
        return false
    }
    if frame.Type != "MEOW_OK" {
        fmt.Println("Unexpected response after HELLO:", string(buf[:n]))
        return false
    }
    fmt.Println("[Server]:", frame.Status)
    return true
}

// sendAuth sends the AUTH frame with username and password.
func sendAuth(stream *quic.Stream, user, pass string) {
    auth := protocol.UniversalFrame{
        Type:  "AUTH",
        MsgID: time.Now().UnixMilli(),
        User:  user,
        Pass:  pass,
    }
    stream.Write(auth.ToJSON())
}

// waitForAuthOK waits for MEOW_OK or ERROR after AUTH.
func waitForAuthOK(stream *quic.Stream) bool {
    buf := make([]byte, 4096)
    n, err := stream.Read(buf)
    if err != nil {
        fmt.Println("Read error after AUTH:", err)
        return false
    }
    frame, err := protocol.ParseFrame(buf[:n])
    if err != nil {
        fmt.Println("Parse error after AUTH:", err)
        return false
    }
    if frame.Type == "ERROR" {
        fmt.Printf("[Auth ERROR] %s: %s\n", frame.Code, frame.Desc)
        return false
    }
    fmt.Println("[Server]:", frame.Status)
    return true
}
