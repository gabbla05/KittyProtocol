package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetChatView(s *state.UIState, target string) fyne.CanvasObject {
	// --- Nagłówek ---
	titleText := canvas.NewText("Chatting with: "+target, color.NRGBA{R: 255, G: 105, B: 180, A: 255})
	titleText.TextSize = 22
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	// --- Historia czatu ---
	chatHistory := widget.NewMultiLineEntry()
	chatHistory.Disable()
	chatHistory.Wrapping = fyne.TextWrapWord
	chatHistory.SetText("--- Chat Started ---\n")

	// --- Odbiór wiadomości z backendu ---
	s.UI.OnMessage = func(from string, text string) {
		fyne.Do(func() {
			chatHistory.SetText(chatHistory.Text + "[" + from + "]: " + text + "\n")
		})
	}

	// --- Pole wpisywania wiadomości ---
	msgEntry := widget.NewEntry()
	msgEntry.SetPlaceHolder("Type your message here...")

	sendBtn := widget.NewButton("Send", func() {
		text := msgEntry.Text
		if text == "" {
			return
		}

		if err := s.App.SendTextMessage(text); err != nil {
			dialog.ShowError(err, s.Window)
			return
		}

		// UI thread
		chatHistory.SetText(chatHistory.Text + "[Me]: " + text + "\n")
		msgEntry.SetText("")
	})
	sendBtn.Importance = widget.HighImportance

	endChatBtn := widget.NewButton("End Chat", func() {
		s.App.EndChat("User closed chat UI")
		s.UI.OnMessage = nil
		s.SwitchView(GetMenuView(s))
	})
	endChatBtn.Importance = widget.DangerImportance

	// --- Układ dolny ---
	bottomBox := container.New(&formLayout{},
		msgEntry,
		sendBtn,
		widget.NewSeparator(),
		endChatBtn,
	)

	historyScroll := container.NewVScroll(chatHistory)
	historyScroll.SetMinSize(fyne.NewSize(350, 300))

	// --- Cały widok ---
	return container.NewBorder(
		container.NewPadded(container.NewCenter(titleText)),
		container.NewCenter(container.NewPadded(bottomBox)),
		nil,
		nil,
		container.NewCenter(historyScroll),
	)
}
