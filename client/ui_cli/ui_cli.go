// ui_cli.go
// CLI implementation of the UI interface used by the App layer.
// This package contains only user interaction logic (stdin/stdout).
// No networking, cryptography or protocol logic is present here.

package ui_cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gabbla05/KittyProtocol/client/api"
)

// CliUI implements the UI interface required by the App layer.
// It is intentionally minimal and synchronous.
type CliUI struct {
	client *api.KittyClient
	reader *bufio.Reader
}

// NewCliUI creates a new CLI frontend bound to a KittyClient instance.
func NewCliUI(c *api.KittyClient) *CliUI {
	return &CliUI{
		client: c,
		reader: bufio.NewReader(os.Stdin),
	}
}

// --- UI interface methods ---

// ReadLine reads a single line from stdin, trimming whitespace.
func (ui *CliUI) ReadLine() string {
	fmt.Print("> ")
	line, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// Println prints a line to stdout.
func (ui *CliUI) Println(v ...any) {
	fmt.Println(v...)
}

// Printf prints a formatted line to stdout.
func (ui *CliUI) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

// --- Additional helpers used by App ---

// ReadCredentials prompts the user for login and password.
func (ui *CliUI) ReadCredentials() (string, string) {
	fmt.Print("Login: ")
	user, _ := ui.reader.ReadString('\n')

	fmt.Print("Hasło: ")
	pass, _ := ui.reader.ReadString('\n')

	return strings.TrimSpace(user), strings.TrimSpace(pass)
}

// ReadSharedSecret prompts the user for the E2EE shared secret.
// It enforces non-empty input.
func (ui *CliUI) ReadSharedSecret() []byte {
	for {
		fmt.Print("Wspólny sekret (K_AB) dla tej rozmowy: ")
		secret, _ := ui.reader.ReadString('\n')
		secret = strings.TrimSpace(secret)

		if secret != "" {
			return []byte(secret)
		}
		fmt.Println("[UI] Sekret nie może być pusty.")
	}
}

// --- ACK event handlers ---

// OnDelivered is called when the AckManager reports successful delivery.
func (ui *CliUI) OnDelivered(msgID int64) {
	fmt.Printf("\n[Delivered] msg_id=%d\n> ", msgID)
}

// OnTimeout is called when the AckManager reports a delivery timeout.
func (ui *CliUI) OnTimeout(msgID int64) {
	fmt.Printf("\n[Timeout] msg_id=%d not delivered\n> ", msgID)
}
