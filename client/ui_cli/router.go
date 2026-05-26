package ui_cli

import (
	"strings"

	"github.com/gabbla05/KittyProtocol/client/app"
	"github.com/gabbla05/KittyProtocol/client/ui_commands"
)

// handleCommand routes a single CLI command.
// Returns true when the menu should exit (logout).
func (ui *CliUI) handleCommand(line string, a *app.App) bool {

	switch {

	// UI-only commands
	case line == "/menu":
		ui.cmdMenu()
		return false

	case line == "/help":
		ui.cmdHelp()
		return false

	// logout
	case line == "/logout":
		msg, err := ui_commands.Logout(a)
		ui.render(msg, err)
		return true

	// status
	case strings.HasPrefix(line, "/status "):
		msg, err := ui_commands.Status(line, a)
		ui.render(msg, err)
		return false

	// secret (requires UI input)
	case strings.HasPrefix(line, "/secret "):
		secret := ui.ReadSharedSecret()
		msg, err := ui_commands.Secret(line, secret, a)
		ui.render(msg, err)
		return false

	// chat
	case strings.HasPrefix(line, "/chat "):
		msg, err := ui_commands.ChatRequest(line, a)
		ui.render(msg, err)
		return false

	case strings.HasPrefix(line, "/accept "):
		msg, err := ui_commands.ChatAccept(line, a)
		ui.render(msg, err)
		return false

	case strings.HasPrefix(line, "/refuse "):
		msg, err := ui_commands.ChatRefuse(line, a)
		ui.render(msg, err)
		return false

	// message
	case strings.HasPrefix(line, "/msg "):
		msg, err := ui_commands.SendMessage(line, a)
		ui.render(msg, err)
		return false

	// end chat
	case line == "/end":
		msg, err := ui_commands.EndChat(a)
		ui.render(msg, err)
		return false

	default:
		ui.Println(ColorYellow + "Unknown command. Use /help." + ColorReset)
		return false
	}
}
