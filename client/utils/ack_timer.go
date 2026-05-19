package clientutils

import (
	"time"
)

// DefaultAckTimeout defines how long the client waits for MEOW_OK
// before considering the message undelivered.
const DefaultAckTimeout = 5 * time.Second

// StartAckTimer starts a timeout watcher for a given message ID.
// If the "done" channel is closed before the timeout, the timer exits silently.
// If the timeout expires first, onTimeout is invoked.
func StartAckTimer(msgID int64, done <-chan struct{}, onTimeout func()) {
	go func() {
		select {
		case <-done:
			// ACK received, timer cancelled.
			return
		case <-time.After(DefaultAckTimeout):
			onTimeout()
		}
	}()
}
