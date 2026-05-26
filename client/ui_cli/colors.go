package ui_cli

// ANSI color codes used for CLI output.
// These are UI-only and do not leak into application logic.
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorPink1  = "\x1b[38;5;213m"
	ColorPink2  = "\x1b[38;5;218m"
	ColorPink3  = "\x1b[38;5;212m"
	ColorYellow = "\033[33m"

	PromptSymbol = ColorPink1 + "(=^._.^=) > " + ColorReset
)
