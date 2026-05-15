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

// handleClient handles a single QUIC connection (one logical client).
// It reads JSON frames from the stream, performs initial type validation,
// and delegates logic to the appropriate modules (auth_flow, router, protection).
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

	// Global session cleanup when the handler finishes.
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
			// Connection closed by client or network error.
			return
		}

		fmt.Println("[RAW]", string(buf[:n]))
		fmt.Println("[Hub] STREAM ID:", stream.StreamID())
		fmt.Println("[Hub] RAW:", string(buf[:n]))

		// TASK 8: Initial validation of frame type and presence of base fields (type, msg_id).
		typeName, _, perr := protocol.GetFrameType(buf[:n])
		if perr != nil {
			sendError(stream, "ERR_02", perr.Error())
			continue
		}

		fmt.Println("[Handler] Received frame type:", typeName) // temporary debug log

		// TASK 8: Check if the frame type is known to the protocol.
		if !protocol.IsValidType(typeName) {
			sendError(stream, "ERR_02", "Unknown frame type: "+typeName)
			continue
		}

		switch typeName {

		case "HELLO":
			// HELLO → MEOW_OK("Ready for auth") + start AUTH timer (Task 9).
			authTimer = handleHELLO(stream, conn)

		case "AUTH":
			// TASK 8: Strict parser validating presence of user and pass.
			frame, err := protocol.ParseAuthFrame(buf[:n])
			if err != nil {
				sendError(stream, "ERR_02", err.Error())
				continue
			}

			// Stop AUTH timer (Task 9).
			if authTimer != nil {
				authTimer.Stop()
			}

			// Credentials verification (Task 6).
			if !auth.CheckCredentials(frame.User, frame.Pass) {
				sendError(stream, "ERR_04", "Authentication failed")
				return
			}

			// Create session and register it in globalSessions (Task 9, Task 5).
			session = protection.NewSession(frame.User, conn, stream)
			globalSessions.Add(frame.User, session)
			username = frame.User

			// TASK 10: Send dedicated MeowOkFrame.
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
			// Update session activity (Idle Timeout – Task 9).
			if session != nil {
				session.LastActive = time.Now()
			}

		case "DATA":
			// TASK 8: Strict parser validating payload and MAC.
			frame, err := protocol.ParseDataFrame(buf[:n])
			if err != nil {
				sendError(stream, "ERR_02", err.Error())
				continue
			}

			// Check if target is not empty.
			if frame.Target == "" {
				sendError(stream, "ERR_02", "Missing target")
				continue
			}

			if session == nil {
				sendError(stream, "ERR_01", "DATA before AUTH")
				continue
			}

			// Per-user rate limiting (Task 9).
			if !session.Limiter.Allow() {
				sendError(stream, "ERR_07", "Rate limit exceeded")
				continue
			}

			// --- ERR_06: Replay protection ---
			if session.Replay != nil && session.Replay.MarkAndCheck(frame.MsgID) {
				sendError(stream, "ERR_06", "Replay detected")
				continue
			}
			// -----------------------------------------

			// Sender is active; Hub does not track any "chat readiness" state.
			session.LastActive = time.Now()

			ok := routeData(*frame, session, stream)
			if !ok {
				// routeData already sent ERR_10 / ERR_15 – nothing more to do here.
				continue
			}

			// ACK for sender (application-level delivery confirmation – Task 14 / MB contract).
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

		case "STATUS_RES":

		case "GET_STATUS":
			// Parsujemy ramkę GET_STATUS
			frame, err := protocol.ParseGetStatusFrame(buf[:n])
			if err != nil {
				sendError(stream, "ERR_02", err.Error())
				continue
			}

			// Sprawdzamy, czy użytkownik jest online
			online := globalSessions.IsOnline(frame.Target)

			status := "offline"
			if online {
				status = "online"
			}

			res := protocol.StatusResFrame{
				BaseFrame: protocol.BaseFrame{
					Type:  "STATUS_RES",
					MsgID: frame.MsgID, // echo msg_id z zapytania
				},
				Target: frame.Target,
				Status: status,
			}

			b, err := json.Marshal(res)
			if err != nil {
				sendError(stream, "ERR_02", "Failed to marshal STATUS_RES")
				continue
			}

			if _, err := stream.Write(b); err != nil {
				fmt.Println("[Hub] Failed to send STATUS_RES:", err)
				continue
			}

		case "BYE":
			// Client intentionally ends the session – we do not shut down the whole Hub,
			// only clean up its entry in the SessionManager.
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
