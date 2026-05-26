package ui_commands

import "github.com/gabbla05/KittyProtocol/client/app"

// Logout handles the logic for the "/logout" command.
// It gracefully terminates any active chat session and then
// sends a BYE frame to the server via the App layer.
//
// This function does not exit the program itself — the UI layer
// (CLI or GUI) decides what "logout" means in its context.
//
// Returns:
//   - string: confirmation message
//   - error: non-nil if sending BYE failed
func Logout(a *app.App) (string, error) {
	if active, peer := a.ChatState().IsActive(); active && peer != "" {
		_ = a.EndChat("user logout")
	}

	_ = a.Client().SendBye()
	return "Logged out.", nil
}
