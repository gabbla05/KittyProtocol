package main

import (
	"context"
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

        frame, perr := protocol.ParseFrame(buf[:n])
        if perr != nil {
            sendError(stream, "ERR_02", perr.Error())
            continue
        }

        switch frame.Type {

        case "HELLO":
            authTimer = handleHELLO(stream, conn)

        case "AUTH":
            if authTimer != nil {
                authTimer.Stop()
            }
            if !auth.CheckCredentials(frame.User, frame.Pass) {
                sendError(stream, "ERR_04", "Authentication failed")
                return
            }

            session = protection.NewSession(frame.User, conn)
            globalSessions.Add(frame.User, session)

            ok := protocol.UniversalFrame{
                Type:   "MEOW_OK",
                MsgID:  frame.MsgID,
                Status: "Logged in",
            }
            stream.Write(ok.ToJSON())

        case "PING":
            if session != nil {
                session.LastActive = time.Now()
            }

        case "DATA":
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
            ack := protocol.UniversalFrame{
                Type:   "MEOW_OK",
                MsgID:  frame.MsgID,
                Status: "Delivered (mock)",
            }
            stream.Write(ack.ToJSON())
        }
    }
}
