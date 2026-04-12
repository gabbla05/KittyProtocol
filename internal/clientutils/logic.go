package clientutils

import (
	"fmt"
	"time"
)

const MaxPayloadSize = 2048

// TruncateMessage przycina tekst, aby zmieścił się w limicie protokołu.
func TruncateMessage(input string) string {
	if len(input) > MaxPayloadSize {
		fmt.Printf("[System] Wiadomość zbyt długa (%d bajtów). Przycinanie...\n", len(input))
		return input[:MaxPayloadSize]
	}
	return input
}

// StartAckTimer to pomocnicza funkcja dla Twojego timera 5s.
func StartAckTimer(msgID int64, onTimeout func()) *chan bool {
	done := make(chan bool)
	go func() {
		select {
		case <-done:
			return
		case <-time.After(5 * time.Second):
			onTimeout()
		}
	}()
	return &done
}
