package ui_commands

import (
	"strings"

	"github.com/gabbla05/KittyProtocol/client/app"
)

// Status handles the logic for the "/status <user>" command.
// It extracts the username from the CLI input, validates it,
// and delegates the status request to the App layer.
//
// This function does not print anything and does not depend on
// terminal-specific features. It returns a message and/or error
// for the UI layer to render.
//
// Returns:
//   - string: a user-facing message
//   - error: non-nil if sending the request failed
func Status(line string, a *app.App) (string, error) {
	user := strings.TrimSpace(strings.TrimPrefix(line, "/status "))
	user = strings.ToLower(user)

	if user == "" {
		return "Usage: /status <user>", nil
	}

	err := a.Client().SendGetStatus(user)
	if err != nil {
		return "", err
	}

	return "Status request sent.", nil
}
