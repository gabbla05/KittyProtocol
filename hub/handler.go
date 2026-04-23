package main

import (
	"context"
	"encoding/json" // Dodano dla json.Marshal
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// handleClient manages the entire lifecycle of a single client connection.
// It implements HELLO → AUTH, rate limiting, and basic DATA handling.
func handleClient(conn *quic.Conn) {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		fmt.Println("Stream error:", err)
		return
	}
	defer stream.Close()

	var authTimer *protection.AuthTimer
	var session *protection.Session

	buf := make([]byte, 4096)

	for {
		n, err := stream.Read(buf)
		if err != nil {
			// Connection closed by client or network error.
			return
		}

		// TASK 8: Wstępna walidacja typu i pobranie MsgID [cite: 78, 120]
		typeName, _, perr := protocol.GetFrameType(buf[:n])
		if perr != nil {
			sendError(stream, "ERR_02", perr.Error())
			continue
		}

		switch typeName {

		case "HELLO":
			authTimer = handleHELLO(stream, conn)

		case "AUTH":
			// TASK 10: Użycie konkretnej struktury AuthFrame [cite: 83]
			var frame protocol.AuthFrame
			if err := json.Unmarshal(buf[:n], &frame); err != nil {
				sendError(stream, "ERR_02", "Invalid AUTH format")
				continue
			}

			if authTimer != nil {
				authTimer.Stop()
			}
			if !auth.CheckCredentials(frame.User, frame.Pass) {
				sendError(stream, "ERR_04", "Authentication failed")
				return
			}

			session = protection.NewSession(frame.User, conn)
			globalSessions.Add(frame.User, session)

			// TASK 10: Wysłanie dedykowanej ramki MeowOkFrame
			ok := protocol.MeowOkFrame{
				BaseFrame: protocol.BaseFrame{
					Type:  "MEOW_OK",
					MsgID: frame.MsgID,
				},
				Status: "Logged in",
			}
			b, _ := json.Marshal(ok)
			stream.Write(b)

		case "PING":
			if session != nil {
				session.LastActive = time.Now()
			}

		case "DATA":
			// TASK 10: Użycie konkretnej struktury DataFrame [cite: 90-94]
			var frame protocol.DataFrame
			if err := json.Unmarshal(buf[:n], &frame); err != nil {
				sendError(stream, "ERR_02", "Invalid DATA format")
				continue
			}

			if session == nil {
				sendError(stream, "ERR_01", "DATA before AUTH")
				continue
			}
			if !session.Limiter.Allow() {
				sendError(stream, "ERR_07", "Rate limit exceeded")
				continue
			}

			// For now, we just acknowledge DATA locally.
			// Later, MB will route this to the target user.
			ack := protocol.MeowOkFrame{
				BaseFrame: protocol.BaseFrame{
					Type:  "MEOW_OK",
					MsgID: frame.MsgID,
				},
				Status: "Delivered (mock)",
			}
			b, _ := json.Marshal(ack)
			stream.Write(b)
		}
	}
}
