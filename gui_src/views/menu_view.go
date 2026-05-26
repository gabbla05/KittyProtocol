package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetMenuView(s *state.UIState) fyne.CanvasObject {
	openChat := widget.NewButton("Open Chat", func() {
		s.Window.SetContent(GetChatView(s))
	})

	return container.NewVBox(
		widget.NewLabel("Menu"),
		openChat,
	)
}
