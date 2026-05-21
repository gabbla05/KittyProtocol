package main

import (
    "encoding/json"
    "fmt"

    "github.com/gabbla05/KittyProtocol/protocol"
)

// handleGetStatus processes a GET_STATUS frame.
//
// This handler performs a pure presence check:
//   - parses the incoming frame,
//   - verifies the target username,
//   - checks whether the target user has an active session,
//   - returns a STATUS_RES frame with "online", "offline", or "no_target".
//
// The Hub does NOT maintain chat‑level state and does NOT track
// which users are currently engaged in a conversation. GET_STATUS
// is strictly a presence query and does not affect routing logic.
func (c *clientContext) handleGetStatus(raw []byte) {
    frame, err := protocol.ParseGetStatusFrame(raw)
    if err != nil {
        sendError(c.stream, "ERR_02", err.Error())
        return
    }

    // Empty target means: "client has no active chat partner".
    // The Hub simply acknowledges this state with a STATUS_RES
    // containing status = "no_target". No session fields are modified.
    if frame.Target == "" {
        res := protocol.StatusResFrame{
            BaseFrame: protocol.BaseFrame{
                Type:  "STATUS_RES",
                MsgID: frame.MsgID,
            },
            Target: "",
            Status: "no_target",
        }

        b, _ := json.Marshal(res)
        c.stream.Write(b)
        return
    }

    // Standard presence check.
    online := globalSessions.IsOnline(frame.Target)
    status := "offline"
    if online {
        status = "online"
    }

    res := protocol.StatusResFrame{
        BaseFrame: protocol.BaseFrame{
            Type:  "STATUS_RES",
            MsgID: frame.MsgID,
        },
        Target: frame.Target,
        Status: status,
    }

    b, err := json.Marshal(res)
    if err != nil {
        sendError(c.stream, "ERR_02", "Failed to marshal STATUS_RES")
        return
    }

    if _, err := c.stream.Write(b); err != nil {
        fmt.Println("[Hub: Status] Failed to send STATUS_RES:", err)
    }
}
