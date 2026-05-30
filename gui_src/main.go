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

	// Inicjalizacja klienta (linijka 26)
	kittyClient := api.NewKittyClient()

	// Inicjalizacja stanu
	stateObj := &state.UIState{
		Window: w,
		Client: kittyClient,
	}
	
	stateObj.UI = &state.GuiUI{State: stateObj}

	// NOWOŚĆ: Przekazujemy logikę przełączania widoku z parametrem docelowego użytkownika
	stateObj.SwitchToChat = func(target string) {
		stateObj.SwitchView(views.GetChatView(stateObj, target))
	}
	
	w.SetContent(views.GetAuthView(stateObj))

	w.ShowAndRun()
}