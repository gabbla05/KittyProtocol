package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
	"github.com/gabbla05/KittyProtocol/gui_src/views"
)

func main() {
	// 1. Inicjalizacja GUI
	a := app.New()
	w := a.NewWindow("Meowssenger")
	
	// Ustawienie początkowego rozmiaru okna, żeby nie było za małe
	w.Resize(fyne.NewSize(400, 300))

	// 2. Inicjalizacja prawdziwego klienta protokołu!
	kittyClient := api.NewKittyClient()

	// 3. Tworzymy obiekt stanu z podpiętym klientem
	stateObj := &state.UIState{
		Window: w,
		Client: kittyClient,
	}
	
	w.SetContent(views.GetAuthView(stateObj))

	w.ShowAndRun()
}