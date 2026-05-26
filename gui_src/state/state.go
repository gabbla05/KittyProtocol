package state

import "fyne.io/fyne/v2"

type UIState struct {
	Window fyne.Window
}

func (s *UIState) SwitchView(content fyne.CanvasObject) {
	s.Window.SetContent(content)
}
