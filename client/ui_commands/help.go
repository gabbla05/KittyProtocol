package ui_commands

// HelpEntry represents a single help item describing a command
// and its purpose. The UI layer is responsible for formatting
// and rendering these entries (e.g., colors, ASCII boxes, etc.).
type HelpEntry struct {
	Command     string
	Description string
}

// Help returns a static list of all supported commands along with
// human-readable descriptions. This data is UI‑agnostic and can be
// rendered differently by CLI and GUI frontends.
//
// The UI layer decides how to present this list (e.g., colored table,
// popup window, tooltip, etc.).
func Help() []HelpEntry {
	return []HelpEntry{
		{"/status <user>", "Check whether a user is online or offline."},
		{"/secret <user> [file:<path>]", "Configure the shared E2EE secret for a peer."},
		{"/chat <user>", "Send a chat request to a user."},
		{"/accept <user>", "Accept an incoming chat request."},
		{"/refuse <user>", "Refuse an incoming chat request."},
		{"/msg <text>", "Send a text message to the active chat partner."},
		{"/end", "End the currently active chat session."},
		{"/logout", "Log out from the server and close the client."},
		{"/menu", "Display the command menu."},
		{"/help", "Display detailed help for all commands."},
	}
}
