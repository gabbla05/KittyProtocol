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

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[36m"
	ColorPink   = "\033[95m"
	ColorYellow = "\033[33m"

	Prompt = ColorPink + "(=^._.^=) > " + ColorReset
)

func (ui *CliUI) Prompt() {
	fmt.Print(Prompt)
}

// CliUI implements the UI interface required by the App layer.
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
	line, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// Println prints a line to stdout.
func (ui *CliUI) Println(v ...any) {
	fmt.Println(v...)
}

// Print prints without newline.
func (ui *CliUI) Print(v ...any) {
	fmt.Print(v...)
}

// Printf prints formatted text.
func (ui *CliUI) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

// --- Additional helpers used by App ---

// ReadCredentials prompts the user for login and password.
func (ui *CliUI) ReadCredentials() (string, string) {
	fmt.Print(ColorBlue + "Login: " + ColorReset)
	user, _ := ui.reader.ReadString('\n')

	fmt.Print(ColorBlue + "Hasło: " + ColorReset)
	pass, _ := ui.reader.ReadString('\n')

	return strings.TrimSpace(user), strings.TrimSpace(pass)
}

// ReadSharedSecret prompts the user for the E2EE shared secret.
func (ui *CliUI) ReadSharedSecret() []byte {
	for {
		fmt.Print(ColorYellow + "Wspólny sekret (K_AB): " + ColorReset)
		secret, _ := ui.reader.ReadString('\n')
		secret = strings.TrimSpace(secret)

		if secret != "" {
			return []byte(secret)
		}
		fmt.Println(ColorRed + "[UI] Sekret nie może być pusty." + ColorReset)
	}
}

// --- ACK event handlers ---

func (ui *CliUI) OnDelivered(msgID int64) {
	fmt.Printf(ColorGreen+"\n[Delivered] msg_id=%d\n"+ColorReset, msgID)
	ui.Prompt()
}

func (ui *CliUI) OnTimeout(msgID int64) {
	fmt.Printf(ColorRed+"\n[Timeout] msg_id=%d not delivered\n"+ColorReset, msgID)
	ui.Prompt()
}
