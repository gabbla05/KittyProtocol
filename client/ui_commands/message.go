package ui_commands

import (
	"strings"

	"github.com/gabbla05/KittyProtocol/client/app"
)

// SendMessage handles the logic for the "/msg <text>" command.
// It extracts the message text from the CLI input and delegates
// the actual encrypted message sending to the App layer.
//
// This function is UI‑agnostic: it does not print anything and does
// not interact with stdin. It simply returns a user-facing message
// and/or an error for the UI layer (CLI or GUI) to render.
//
// Returns:
//   - string: confirmation message for the UI
//   - error: non-nil if sending the message failed
func SendMessage(line string, a *app.App) (string, error) {
	text := strings.TrimSpace(strings.TrimPrefix(line, "/msg "))
	if text == "" {
		return "", nil
	}

	err := a.SendTextMessage(text)
	if err != nil {
		return "", err
	}

	return "Message sent.", nil
}
