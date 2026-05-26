package ui_commands

import "github.com/gabbla05/KittyProtocol/client/app"

// EndChat handles the logic for the "/end" command.
// It terminates the currently active chat session by delegating
// the operation to the App layer.
//
// This function is UI‑agnostic and returns a message and/or error
// for the UI layer to render.
//
// Returns:
//   - string: confirmation message
//   - error: non-nil if ending the chat failed
func EndChat(a *app.App) (string, error) {
	err := a.EndChat("user ended chat")
	if err != nil {
		return "", err
	}
	return "Chat ended.", nil
}
