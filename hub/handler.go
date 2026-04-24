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

		// TASK 8: Wstępna walidacja typu i obecności pól bazowych (type, msg_id)
		typeName, _, perr := protocol.GetFrameType(buf[:n])
		if perr != nil {
			sendError(stream, "ERR_02", perr.Error())
			continue
		}

		// TASK 8: Sprawdzenie czy typ ramki jest znany protokołowi
		if !protocol.IsValidType(typeName) {
			sendError(stream, "ERR_02", "Unknown frame type: "+typeName)
			continue
		}

		switch typeName {

		case "HELLO":
			authTimer = handleHELLO(stream, conn)

		case "AUTH":
			// TASK 8: Rygorystyczny parser sprawdzający obecność user i pass
			frame, err := protocol.ParseAuthFrame(buf[:n])
			if err != nil {
				sendError(stream, "ERR_02", err.Error())
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

			// TASK 10: Wysłanie dedykowanej ramki MeowOkFrame [cite: 133]
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
			// TASK 8: Rygorystyczny parser sprawdzający payload i MAC [cite: 81, 117]
			frame, err := protocol.ParseDataFrame(buf[:n])
			if err != nil {
				sendError(stream, "ERR_02", err.Error())
				continue
			}

			if session == nil {
				sendError(stream, "ERR_01", "DATA before AUTH") // [cite: 278]
				continue
			}
			if !session.Limiter.Allow() {
				sendError(stream, "ERR_07", "Rate limit exceeded") // [cite: 278]
				continue
			}

			// Aplikacyjne potwierdzenie (ACK) - [cite: 48, 55]
			ack := protocol.MeowOkFrame{
				BaseFrame: protocol.BaseFrame{
					Type:  "MEOW_OK",
					MsgID: frame.MsgID,
				},
				Status: "Delivered (mock)",
			}
			b, _ := json.Marshal(ack)
			stream.Write(b)

		case "BYE":
			// Celowe zakończenie sesji [cite: 72]
			return
		}
	}
}
