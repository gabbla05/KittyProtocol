package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_"github.com/gabbla05/KittyProtocol/gui_src/resources" // Poprawny import
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetChatView(s *state.UIState, target string) fyne.CanvasObject {
	title := widget.NewLabel("Chatting with: " + target)
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Duże, rozwijane pole z historią wiadomości
	chatHistory := widget.NewMultiLineEntry()
	chatHistory.Disable() 
	chatHistory.Wrapping = fyne.TextWrapWord
	chatHistory.SetText("--- Chat Started ---\n")

	// NASŁUCHIWANIE
	s.UI.OnMessage = func(msg string) {
		chatHistory.SetText(chatHistory.Text + msg + "\n")
	}

	msgEntry := widget.NewEntry()
	msgEntry.SetPlaceHolder("Type your message here...")

	sendBtn := widget.NewButton("Send", func() {
		text := msgEntry.Text
		if text == "" {
			return
		}

		err := s.App.SendTextMessage(text)
		if err != nil {
			dialog.ShowError(err, s.Window)
			return
		}

		chatHistory.SetText(chatHistory.Text + "[Ja]: " + text + "\n")
		msgEntry.SetText("")
	})
	sendBtn.Importance = widget.HighImportance

	endChatBtn := widget.NewButton("End Chat", func() {
		s.App.EndChat("User closed chat UI")
		s.UI.OnMessage = nil
		s.SwitchView(GetMenuView(s))
	})
	endChatBtn.Importance = widget.DangerImportance

	// Układ (Border)
	bottomBox := container.NewVBox(msgEntry, sendBtn, widget.NewSeparator(), endChatBtn)
	historyScroll := container.NewVScroll(chatHistory)

	return container.NewBorder(title, bottomBox, nil, nil, historyScroll)
}