package ui_commands

import (
	"errors"
	"strings"

	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/client/app"
)

// ChatRequest handles the logic for the "/chat <user>" command.
// It extracts the target username from the CLI input, validates it,
// and delegates the actual chat request initiation to the App layer.
//
// This function is UI‑agnostic: it does not print anything, does not
// read from stdin, and does not depend on terminal colors. Instead,
// it returns a user-facing message (string) and/or an error, which
// the UI layer (CLI or GUI) is responsible for rendering.
//
// Returns:
//   - string: a human-readable message for the UI
//   - error: non-nil if the operation failed
func ChatRequest(line string, a *app.App) (string, error) {
	user := strings.TrimSpace(strings.TrimPrefix(line, "/chat "))
	user = strings.ToLower(user)

	if user == "" {
		return "Usage: /chat <user>", nil
	}

	err := a.StartChatRequest(user)
	if err != nil {
		if errors.Is(err, api.ErrNoSharedSecret) {
			return "No shared secret with " + user + ". Use /secret " + user, nil
		}
		return "", err
	}

	return "CHAT_REQUEST sent to " + user, nil
}

// ChatAccept handles the logic for the "/accept <user>" command.
// It validates the username and delegates the acceptance of a chat
// request to the App layer.
//
// UI‑agnostic: returns a message and/or error for the UI to render.
func ChatAccept(line string, a *app.App) (string, error) {
	user := strings.TrimSpace(strings.TrimPrefix(line, "/accept "))
	user = strings.ToLower(user)

	if user == "" {
		return "Usage: /accept <user>", nil
	}

	err := a.AcceptChat(user)
	if err != nil {
		return "", err
	}

	return "Chat accepted with " + user, nil
}

// ChatRefuse handles the logic for the "/refuse <user>" command.
// It validates the username and delegates the refusal of a chat
// request to the App layer.
//
// UI‑agnostic: returns a message and/or error for the UI to render.
func ChatRefuse(line string, a *app.App) (string, error) {
	user := strings.TrimSpace(strings.TrimPrefix(line, "/refuse "))
	user = strings.ToLower(user)

	if user == "" {
		return "Usage: /refuse <user>", nil
	}

	err := a.RefuseChat(user, "user refused")
	if err != nil {
		return "", err
	}

	return "Chat refused with " + user, nil
}
