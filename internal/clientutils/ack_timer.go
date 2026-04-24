package clientutils

import (
	"time"
)

// StartAckTimer starts a 5-second timer for a given message ID.
// If the "done" channel is not closed within 5 seconds, onTimeout is called.
func StartAckTimer(msgID int64, done <-chan struct{}, onTimeout func()) {
    go func() {
        select {
        case <-done:
            // ACK received, timer cancelled.
            return
        case <-time.After(5 * time.Second):
            onTimeout()
        }
    }()
}
