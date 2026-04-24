package main

import (
	"encoding/json" // Dodano dla json.Marshal i json.Unmarshal
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// sendHello wysyła wstępną ramkę HELLO do Huba.
func sendHello(stream *quic.Stream) {
	// TASK 10: Użycie HelloFrame zamiast UniversalFrame
	hello := protocol.HelloFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "HELLO",
			MsgID: time.Now().UnixMilli(),
		},
	}
	b, _ := json.Marshal(hello)
	stream.Write(b)
}

// waitForHelloOK czeka na MEOW_OK po wysłaniu HELLO i wypisuje status serwera.
// waitForHelloOK czeka na MEOW_OK po HELLO i zwraca (sukces, kod_błędu).
func waitForHelloOK(stream *quic.Stream) (bool, string) {
	buf := make([]byte, 4096)
	n, err := stream.Read(buf)
	if err != nil {
		// Sprawdzamy specyficzne błędy sieciowe QUIC
		if qerr, ok := err.(*quic.ApplicationError); ok && qerr.ErrorCode == 0x03 {
			return false, "ERR_03"
		}
		fmt.Println("Read error after HELLO:", err)
		return false, "CONNECTION_LOST"
	}

	// TASK 8: Wstępna walidacja typu i pobranie MsgID
	typeName, _, err := protocol.GetFrameType(buf[:n])
	if err != nil {
		fmt.Println("Parse error after HELLO:", err)
		return false, "PARSE_ERROR"
	}

	// Obsługa ramki ERROR zamiast MEOW_OK (np. błąd formatu HELLO)
	if typeName == "ERROR" {
		var errFrame protocol.ErrorFrame
		json.Unmarshal(buf[:n], &errFrame)
		fmt.Printf("\n[Server ERROR] %s: %s\n", errFrame.Code, errFrame.Desc)
		return false, errFrame.Code
	}

	if typeName != "MEOW_OK" {
		fmt.Println("Unexpected response after HELLO:", string(buf[:n]))
		return false, "UNKNOWN_FRAME"
	}

	// TASK 10: Parsowanie do dedykowanej struktury MeowOkFrame
	var okFrame protocol.MeowOkFrame
	if err := json.Unmarshal(buf[:n], &okFrame); err != nil {
		return false, "PARSE_ERROR"
	}

	fmt.Println("[Server]:", okFrame.Status)
	return true, ""
}

// sendAuth wysyła ramkę AUTH z loginem i hasłem.
func sendAuth(stream *quic.Stream, user, pass string) {
	// TASK 10: Użycie AuthFrame zamiast UniversalFrame
	auth := protocol.AuthFrame{
		BaseFrame: protocol.BaseFrame{
			Type:  "AUTH",
			MsgID: time.Now().UnixMilli(),
		},
		User: user,
		Pass: pass,
	}
	b, _ := json.Marshal(auth)
	stream.Write(b)
}

// waitForAuthOK czeka na MEOW_OK lub ERROR po wysłaniu AUTH.
func waitForAuthOK(stream *quic.Stream) (bool, string) {
	buf := make([]byte, 4096)
	n, err := stream.Read(buf)
	if err != nil {
		fmt.Println("Read error after AUTH:", err)
		return false, "READ_ERROR"
	}

	typeName, _, err := protocol.GetFrameType(buf[:n])
	if err != nil {
		fmt.Println("Parse error after AUTH:", err)
		return false, "PARSE_ERROR"
	}

	if typeName == "ERROR" {
		var errFrame protocol.ErrorFrame
		json.Unmarshal(buf[:n], &errFrame)
		fmt.Printf("[Auth ERROR] %s: %s\n", errFrame.Code, errFrame.Desc)
		return false, errFrame.Code // Zwracamy konkretny kod (np. ERR_03)
	}

	if typeName == "MEOW_OK" {
		var okFrame protocol.MeowOkFrame
		json.Unmarshal(buf[:n], &okFrame)
		fmt.Println("[Server]:", okFrame.Status)
		return true, ""
	}

	return false, "UNKNOWN"
}
