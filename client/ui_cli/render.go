package ui_cli

// render prints the result of a ui_commands call.
func (ui *CliUI) render(msg string, err error) {
	if err != nil {
		ui.Println(ColorRed + "[ERROR] " + err.Error() + ColorReset)
		return
	}
	if msg != "" {
		ui.Println(ColorGreen + msg + ColorReset)
	}
}
