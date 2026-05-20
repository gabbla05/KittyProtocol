package ui_cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gabbla05/KittyProtocol/client/api"
)

// CliUI implements the UI interface used by the App layer.
type CliUI struct {
	client *api.KittyClient
	reader *bufio.Reader
}

func NewCliUI(c *api.KittyClient) *CliUI {
	return &CliUI{
		client: c,
		reader: bufio.NewReader(os.Stdin),
	}
}

// --- UI interface methods ---

func (ui *CliUI) ReadLine() string {
	fmt.Print("> ")
	line, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func (ui *CliUI) Println(v ...any) {
	fmt.Println(v...)
}

func (ui *CliUI) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

// --- Additional helpers used by App ---

func (ui *CliUI) ReadCredentials() (string, string) {
	fmt.Print("Login: ")
	user, _ := ui.reader.ReadString('\n')

	fmt.Print("Hasło: ")
	pass, _ := ui.reader.ReadString('\n')

	return strings.TrimSpace(user), strings.TrimSpace(pass)
}

func (ui *CliUI) ReadSharedSecret() []byte {
	for {
		fmt.Print("Wspólny sekret (K_AB) dla tej rozmowy: ")
		secret, _ := ui.reader.ReadString('\n')
		secret = strings.TrimSpace(secret)

		if secret != "" {
			return []byte(secret)
		}
		fmt.Println("[Client: UI-cli] Sekret nie może być pusty.")
	}
}

// --- ACK event handlers ---

func (ui *CliUI) OnDelivered(msgID int64) {
	fmt.Printf("\n[Delivered] msg_id=%d\n> ", msgID)
}

func (ui *CliUI) OnTimeout(msgID int64) {
	fmt.Printf("\n[Timeout] msg_id=%d not delivered\n> ", msgID)
}
