package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetChatView(s *state.UIState) fyne.CanvasObject {
	back := widget.NewButton("Back", func() {
		s.Window.SetContent(GetMenuView(s))
	})

	return container.NewVBox(
		widget.NewLabel("Chat"),
		back,
	)
}
