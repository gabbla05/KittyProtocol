package ui_cli

import "fmt"

// OnDelivered is called when the AckManager reports successful delivery.
// This is UI-only feedback and does not affect application logic.
func (ui *CliUI) OnDelivered(msgID int64) {
	fmt.Printf(ColorGreen+"\n[Delivered] msg_id=%d\n"+ColorReset, msgID)
	ui.Prompt()
}

// OnTimeout is called when the AckManager reports a delivery timeout.
// This is UI-only feedback and does not affect application logic.
func (ui *CliUI) OnTimeout(msgID int64) {
	fmt.Printf(ColorRed+"\n[Timeout] msg_id=%d not delivered\n"+ColorReset, msgID)
	ui.Prompt()
}
