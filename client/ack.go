package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gabbla05/KittyProtocol/internal/clientutils"
	"github.com/gabbla05/KittyProtocol/protocol"
	"github.com/quic-go/quic-go"
)

// startReceiverLoop listens for incoming frames and handles MEOW_OK and ERROR.
// It returns a map used to track pending ACKs.
func startReceiverLoop(stream *quic.Stream) (map[int64]chan struct{}, *sync.Mutex) {
    pending := make(map[int64]chan struct{})
    var mu sync.Mutex

    go func() {
        buf := make([]byte, 4096)
        for {
            n, err := stream.Read(buf)
            if err != nil {
                fmt.Println("\n[Client] Connection closed:", err)
                return
            }
            frame, err := protocol.ParseFrame(buf[:n])
            if err != nil {
                fmt.Println("[Client] Parse error:", err)
                continue
            }

            switch frame.Type {
            case "MEOW_OK":
                mu.Lock()
                if ch, ok := pending[frame.MsgID]; ok {
                    close(ch)
                    delete(pending, frame.MsgID)
                    fmt.Printf("\n[Delivered] msg_id=%d\n> ", frame.MsgID)
                }
                mu.Unlock()
            case "ERROR":
                fmt.Printf("\n[Server ERROR] %s: %s\n> ", frame.Code, frame.Desc)
            }
        }
    }()

    return pending, &mu
}

// sendMessage sends a DATA frame and starts a 5-second ACK timer.
func sendMessage(stream *quic.Stream, target, text string, pending map[int64]chan struct{}, mu *sync.Mutex) {
    safe := clientutils.TruncateMessage(text)
    msgID := time.Now().UnixMilli()

    ch := make(chan struct{})
    mu.Lock()
    pending[msgID] = ch
    mu.Unlock()

    clientutils.StartAckTimer(msgID, ch, func() {
        mu.Lock()
        if _, ok := pending[msgID]; ok {
            delete(pending, msgID)
            fmt.Printf("\n[Timeout] msg_id=%d not delivered\n> ", msgID)
        }
        mu.Unlock()
    })

    frame := protocol.UniversalFrame{
        Type:    "DATA",
        MsgID:   msgID,
        Target:  target,
        Payload: safe,
    }
    stream.Write(frame.ToJSON())
}
