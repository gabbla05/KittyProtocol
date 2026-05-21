// hub/router.go
// Routing logic for DATA frames between active user sessions.

package main

import (
    "encoding/json"
    "fmt"

    "github.com/gabbla05/KittyProtocol/internal/protection"
    "github.com/gabbla05/KittyProtocol/protocol"
    "github.com/quic-go/quic-go"
)

// routeData forwards a DATA frame from the sender to the target session.
//
// This function performs only transport‑level validation:
//   - verifies that the target session exists,
//   - verifies that the target stream is available,
//   - forwards the DATA frame unchanged except for the Sender field.
//
// The Hub does NOT interpret application‑level payloads and does NOT
// enforce chat‑level rules (e.g., active conversation state). All
// application logic is handled by clients. The Hub acts strictly as
// a message router.
func routeData(frame protocol.DataFrame, sender *protection.Session, senderStream *quic.Stream) bool {
    targetSess, ok := globalSessions.Get(frame.Target)
    if !ok {
        sendError(senderStream, "ERR_15", "Receiver offline")
        return false
    }

    if targetSess.Stream == nil {
        sendError(senderStream, "ERR_10", "Receiver stream not available")
        return false
    }

    // Construct the forwarded DATA frame. The Hub does not modify
    // application payloads or metadata beyond setting the Sender field.
    forward := protocol.DataFrame{
        BaseFrame: protocol.BaseFrame{
            Type:  "DATA",
            MsgID: frame.MsgID,
        },
        Sender:  sender.ID,
        Target:  frame.Target,
        Payload: frame.Payload,
        MAC:     frame.MAC,
    }

    fb, err := json.Marshal(forward)
    if err != nil {
        fmt.Println("[Hub: Router] Failed to marshal forwarded DATA:", err)
        sendError(senderStream, "ERR_02", "Failed to marshal forwarded DATA")
        return false
    }

    if _, err := targetSess.Stream.Write(fb); err != nil {
        fmt.Println("[Hub: Router] Failed to deliver DATA to receiver:", err)
        sendError(senderStream, "ERR_10", "Failed to deliver to receiver")
        return false
    }

    return true
}
