package ui_cli

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ReadLine reads a single line from stdin and trims whitespace.
func (ui *CliUI) ReadLine() string {
	line, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// ReadCredentials prompts the user for login and password.
// The password is read without echo using term.ReadPassword.
func (ui *CliUI) ReadCredentials() (string, string) {
	fmt.Print(ColorBlue + " -> Login: " + ColorReset)
	user, _ := ui.reader.ReadString('\n')
	user = strings.TrimSpace(user)

	fmt.Print(ColorBlue + " -> Password: " + ColorReset)
	bytePass, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	pass := strings.TrimSpace(string(bytePass))

	return user, pass
}

// ReadSharedSecret prompts the user for an E2EE shared secret.
// It enforces non-empty input and loops until valid.
func (ui *CliUI) ReadSharedSecret() []byte {
	for {
		fmt.Print(ColorYellow + " -> Shared secret (K_AB): " + ColorReset)
		secret, _ := ui.reader.ReadString('\n')
		secret = strings.TrimSpace(secret)

		if secret != "" {
			return []byte(secret)
		}
		fmt.Println(ColorRed + "[UI] Secret cannot be empty." + ColorReset)
	}
}
