package ui_cli

import "fmt"

// Prompt prints the standard CLI prompt symbol.
func (ui *CliUI) Prompt() {
	fmt.Print(PromptSymbol)
}

// Println prints a line to stdout.
func (ui *CliUI) Println(v ...any) {
	fmt.Println(v...)
}

// Print prints without newline.
func (ui *CliUI) Print(v ...any) {
	fmt.Print(v...)
}

// Printf prints formatted text.
func (ui *CliUI) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}
