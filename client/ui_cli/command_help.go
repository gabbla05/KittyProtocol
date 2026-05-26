package ui_cli

import "github.com/gabbla05/KittyProtocol/client/ui_commands"

// cmdHelp prints a styled help page with detailed command descriptions.
func (ui *CliUI) cmdHelp() {
	ui.Println(ColorPink1 + "\n  ╔══════════════════════════════════════════════════════╗" + ColorReset)
	ui.Println(ColorPink1 + "  ║                    Kitty CLI Help                    ║" + ColorReset)
	ui.Println(ColorPink1 + "  ╚══════════════════════════════════════════════════════╝\n" + ColorReset)

	ui.Println(ColorPink2 + "  Available Commands:" + ColorReset)

	for _, h := range ui_commands.Help() {
		ui.Println()
		ui.Printf(ColorGreen+"    %-22s"+ColorReset+"\n", h.Command)
		ui.Printf("      %s\n", h.Description)
	}

	ui.Println(ColorPink3 + "\n  ════════════════════════════════════════════════════════" + ColorReset)
	ui.Println(ColorPink3 + "  Tip: Use /menu to show the compact command list anytime." + ColorReset)
	ui.Println(ColorPink3 + "  ════════════════════════════════════════════════════════\n" + ColorReset)
}
