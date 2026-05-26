package ui_cli

// printMenu prints the main command menu.
func (ui *CliUI) printMenu() {
	ui.Println(ColorPink2 + "\n  ======================" + ColorReset)
	ui.Println(ColorPink2 + " | " + ColorGreen + "Available commands:  " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /status <user>    " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /secret <user>    " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /chat <user>      " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /accept <user>    " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /refuse <user>    " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /msg <text>       " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /end              " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /logout           " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /menu             " + ColorPink2 + "|")
	ui.Println(" | " + ColorGreen + "->" + ColorReset + " /help             " + ColorPink2 + "|")
	ui.Println(ColorPink2 + "  ======================\n" + ColorReset)
}
