package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// CliUI implements AckEventHandler and provides a simple terminal UI.
type CliUI struct {
	client *KittyClient
	reader *bufio.Reader
}

// NewCliUI creates a new CLI wrapper for KittyClient.
func NewCliUI(c *KittyClient) *CliUI {
	return &CliUI{
		client: c,
		reader: bufio.NewReader(os.Stdin),
	}
}

// OnDelivered is called by AckManager when a message is delivered.
func (ui *CliUI) OnDelivered(msgID int64) {
	fmt.Printf("\n[Delivered] msg_id=%d\n> ", msgID)
}

// OnTimeout is called by AckManager when a message times out.
func (ui *CliUI) OnTimeout(msgID int64) {
	fmt.Printf("\n[Timeout] msg_id=%d not delivered\n> ", msgID)
}

// ReadCredentials prompts the user for login and password.
func (ui *CliUI) ReadCredentials() (string, string) {
	fmt.Print("Login: ")
	user, _ := ui.reader.ReadString('\n')

	fmt.Print("Hasło: ")
	pass, _ := ui.reader.ReadString('\n')

	return strings.TrimSpace(user), strings.TrimSpace(pass)
}

// ReadTarget prompts the user for the chat target.
func (ui *CliUI) ReadTarget() string {
	for {
		fmt.Print("Do kogo piszesz?: ")
		target, _ := ui.reader.ReadString('\n')
		target = strings.TrimSpace(target)

		if target != "" {
			return target
		}
		fmt.Println("[Client: UI-cli] Target cannot be empty.")
	}
}

func (ui *CliUI) ReadSharedSecret() string {
	for {
		fmt.Print("Wspólny sekret (K_AB) dla tej rozmowy: ")
		secret, _ := ui.reader.ReadString('\n')
		secret = strings.TrimSpace(secret)

		if secret != "" {
			return secret
		}
		fmt.Println("[Client: UI-cli] Sekret nie może być pusty.")
	}
}

// RunSendLoop starts the main CLI loop for sending messages.
func (ui *CliUI) RunSendLoop(disconnected chan struct{}) {
	for {
		select {
		case <-disconnected:
			fmt.Println("[Client: UI-cli] Session closed. Returning to disconnected state.")
			return

		default:
		}

		fmt.Print("> ")
		text, err := ui.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\n[Client: UI-cli] EOF detected. Sending BYE and exiting.")
				_ = ui.client.SendBye()
				ui.client.Close()
				return
			}
			fmt.Println("[Client: UI-cli] Read error:", err)
			continue
		}

		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		// Local commands
		if text == "/quit" {
			_ = ui.client.SendBye()
			ui.client.Close()
			fmt.Println("[Client: UI-cli] Closing session by user request.")
			return
		}

		if after, ok := strings.CutPrefix(text, "/status "); ok {
			target := strings.TrimSpace(after)
			if target == "" {
				fmt.Println("[Client: UI-cli] Usage: /status <user>")
				continue
			}
			_ = ui.client.SendGetStatus(target)
			continue
		}

		if text == "/replay" {
			ui.client.mu.Lock()
			frame := ui.client.lastFrame
			ui.client.mu.Unlock()

			if frame == nil {
				fmt.Println("[Client] No message to replay.")
				continue
			}

			_, err := ui.client.stream.Write(frame)
			if err != nil {
				fmt.Println("[Client: UI-cli] Replay send error:", err)
			} else {
				fmt.Println("[Client] Replay sent.")
			}

			continue
		}

		// Normal message
		const MaxMessageLen = 2000

		// Truncate if too long
		if len(text) > MaxMessageLen {
			fmt.Printf("[Client] Message too long (%d chars). Truncated to %d.\n",
				len(text), MaxMessageLen)
			text = text[:MaxMessageLen]
		}

		if err := ui.client.SendMessage(text); err != nil {
			fmt.Println("[Client: UI-cli] Send error:", err)
		}

	}
}
