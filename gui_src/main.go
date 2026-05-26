package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
	"github.com/gabbla05/KittyProtocol/gui_src/views"
)

func main() {
	// 1. Inicjalizacja GUI
	a := app.New()
	w := a.NewWindow("Meowssenger")

	// 2. Tworzymy obiekt stanu (przekażemy go do widoków, żeby mogły przełączać ekrany)
	stateObj := &state.UIState{Window: w}
	w.SetContent(views.GetAuthView(stateObj))

	w.ShowAndRun()
}
