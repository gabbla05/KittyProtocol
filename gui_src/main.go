package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/gui_src/resources"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
	"github.com/gabbla05/KittyProtocol/gui_src/theme"
	"github.com/gabbla05/KittyProtocol/gui_src/views"
)

func main() {
    a := app.New()
    
    // Jeśli PinkTheme poprawnie implementuje fyne.Theme, użyj tego:
    a.Settings().SetTheme(&theme.PinkTheme{})
    
    w := a.NewWindow("Meowssenger")
    w.SetIcon(resources.LogoIkonaPng)
    // Ustawienie stałego rozmiaru
    w.Resize(fyne.NewSize(400, 650))
    w.SetFixedSize(true) // TO BLOKUJE ZMIANĘ ROZMIARU

	// Inicjalizacja klienta
	kittyClient := api.NewKittyClient()

	// Inicjalizacja stanu
	stateObj := &state.UIState{
		Window: w,
		Client: kittyClient,
	}
	
	w.SetContent(views.GetAuthView(stateObj))

	w.ShowAndRun()
}