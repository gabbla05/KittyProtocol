package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetChatView(s *state.UIState, target string) fyne.CanvasObject {
	title := widget.NewLabel("Chatting with: " + target)
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Duże, rozwijane pole z historią wiadomości
	chatHistory := widget.NewMultiLineEntry()
	chatHistory.Disable() // Zablokowane do edycji przez użytkownika
	chatHistory.Wrapping = fyne.TextWrapWord
	chatHistory.SetText("--- Chat Started ---\n")

	// NASŁUCHIWANIE: Każda akcja z tła (wiadomości, statusy) wpadnie tutaj!
	s.UI.OnMessage = func(msg string) {
		current := chatHistory.Text
		chatHistory.SetText(current + msg)
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

		// Po udanym wysłaniu, pokazujemy NASZĄ wiadomość u siebie w historii
		chatHistory.SetText(chatHistory.Text + "\n[Ja]: " + text + "\n")
		msgEntry.SetText("")
	})
	sendBtn.Importance = widget.HighImportance

	endChatBtn := widget.NewButton("End Chat", func() {
		s.App.EndChat("User closed chat UI")
		s.UI.OnMessage = nil // Odpinamy nasłuchiwanie by nie było bugów
		s.SwitchView(GetMenuView(s))
	})
	endChatBtn.Importance = widget.DangerImportance

	// Układ (Border) sprawi, że historia zajmie max dostępnego miejsca, a reszta przylgnie do krawędzi
	bottomBox := container.NewVBox(msgEntry, sendBtn, widget.NewSeparator(), endChatBtn)
	
	// Wrzucamy chatHistory w Scroll, żeby dało się przewijać suwakiem
	historyScroll := container.NewVScroll(chatHistory)

	return container.NewBorder(title, bottomBox, nil, nil, historyScroll)
}