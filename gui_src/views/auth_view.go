package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetAuthView(s *state.UIState) fyne.CanvasObject {
	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Username")

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Password")

	loginBtn := widget.NewButton("Login", func() {
		// TU PODPIĘCIE: zakładamy sukces logowania i przełączamy widok
		s.Window.SetContent(GetMenuView(s))
	})

	return container.NewVBox(
		widget.NewLabel("Witaj w Meowssenger"),
		userEntry,
		passEntry,
		loginBtn,
	)
}
