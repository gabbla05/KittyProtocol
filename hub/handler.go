// hub/handler.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/auth"
	"github.com/gabbla05/KittyProtocol/internal/protection"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// handleClient obsługuje pojedyncze połączenie QUIC (jeden klient).
// Czyta ramki JSON ze strumienia, wykonuje wstępną walidację typu
// i deleguje logikę do odpowiednich modułów (auth_flow, router, protection).
func handleClient(conn *quic.Conn) {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		fmt.Println("Stream error:", err)
		return
	}
	defer stream.Close()

	var authTimer *protection.AuthTimer
	var session *protection.Session
	var username string

	// Globalne sprzątanie sesji po zakończeniu handlera.
	defer func() {
		if session != nil {
			fmt.Println("[Handler] Cleaning up session for:", username)
			globalSessions.Remove(username)
			if session.CloseFunc != nil {
				session.CloseFunc()
			}
		}
	}()

	buf := make([]byte, 4096)

	for {
		n, err := stream.Read(buf)
		if err != nil {
			// Połączenie zamknięte przez klienta lub błąd sieci.
			return
		}

		fmt.Println("[RAW]", string(buf[:n]))
		fmt.Println("[Hub] STREAM ID:", stream.StreamID())
		fmt.Println("[Hub] RAW:", string(buf[:n]))

		// TASK 8: Wstępna walidacja typu i obecności pól bazowych (type, msg_id)
		typeName, _, perr := protocol.GetFrameType(buf[:n])
		if perr != nil {
			sendError(stream, "ERR_02", perr.Error())
			continue
		}

		fmt.Println("[Handler] Received frame type:", typeName) // tymczasowo

		// TASK 8: Sprawdzenie czy typ ramki jest znany protokołowi
		if !protocol.IsValidType(typeName) {
			sendError(stream, "ERR_02", "Unknown frame type: "+typeName)
			continue
		}

		switch typeName {

		case "HELLO":
			// HELLO → MEOW_OK("Ready for auth") + start timera AUTH (Task 9)
			authTimer = handleHELLO(stream, conn)

		case "AUTH":
			// TASK 8: Rygorystyczny parser sprawdzający obecność user i pass
			frame, err := protocol.ParseAuthFrame(buf[:n])
			if err != nil {
				sendError(stream, "ERR_02", err.Error())
				continue
			}

			// Zatrzymanie timera AUTH (Task 9)
			if authTimer != nil {
				authTimer.Stop()
			}

			// Weryfikacja poświadczeń (Task 6)
			if !auth.CheckCredentials(frame.User, frame.Pass) {
				sendError(stream, "ERR_04", "Authentication failed")
				return
			}

			// Utworzenie sesji i rejestracja w globalSessions (Task 9, Task 5)
			session = protection.NewSession(frame.User, conn, stream)
			globalSessions.Add(frame.User, session)
			username = frame.User

			// TASK 10: Wysłanie dedykowanej ramki MeowOkFrame
			ok := protocol.MeowOkFrame{
				BaseFrame: protocol.BaseFrame{
					Type:  "MEOW_OK",
					MsgID: frame.MsgID,
				},
				Status: "Logged in",
			}
			if b, err := json.Marshal(ok); err == nil {
				stream.Write(b)
			}

		case "PING":
			// Aktualizacja aktywności sesji (Idle Timeout – Task 9)
			if session != nil {
				session.LastActive = time.Now()
			}

		case "DATA":
			// TASK 8: Rygorystyczny parser sprawdzający payload i MAC
			frame, err := protocol.ParseDataFrame(buf[:n])
			if err != nil {
				sendError(stream, "ERR_02", err.Error())
				continue
			}

			// check if target is not empty
			if frame.Target == "" {
				sendError(stream, "ERR_02", "Missing target")
				continue
			}

			if session == nil {
				sendError(stream, "ERR_01", "DATA before AUTH")
				continue
			}

			// Rate limiting per user (Task 9)
			if !session.Limiter.Allow() {
				sendError(stream, "ERR_07", "Rate limit exceeded")
				continue
			}

			// Nadawca jest aktywny i „gotowy do czatu”
			session.LastActive = time.Now()
			session.ReadyForChat = true // <-- NOWE: oznaczamy, że ten user już wszedł w tryb rozmowy

			ok := routeData(*frame, session, stream)
			if !ok {
				// routeData samo wysłało ERR_10 / ERR_15 / ERR_16 – nic więcej nie robimy
				continue
			}

			// ACK dla nadawcy (aplikacyjne potwierdzenie – Task 14/MB kontrakt)
			ack := protocol.MeowOkFrame{
				BaseFrame: protocol.BaseFrame{
					Type:  "MEOW_OK",
					MsgID: frame.MsgID,
				},
				Status: "Delivered (mock)",
			}
			if b, err := json.Marshal(ack); err == nil {
				stream.Write(b)
			}

		case "BYE":
			// Klient świadomie kończy sesję – nie zamykamy całego Huba,
			// tylko sprzątamy jego wpis w SessionManager.
			if session != nil {
				fmt.Println("[Handler] Cleaning up session for:", username)
				globalSessions.Remove(username)
				if session.CloseFunc != nil {
					session.CloseFunc()
				}
			}
			return
		}
	}
}
