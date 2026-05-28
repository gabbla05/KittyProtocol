package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/gabbla05/KittyProtocol/gui_src/resources"
)

type PinkTheme struct{}

var _ fyne.Theme = (*PinkTheme)(nil)

// Kolor jasnoróżowy (ten sam dla przycisków i ramek)
var pinkColor = color.NRGBA{R: 255, G: 182, B: 200, A: 255}
// Ciemniejszy róż dla aktywnej ramki (Focus)
var focusPink = color.NRGBA{R: 255, G: 150, B: 180, A: 255}
// Jasnoróżowe tło całej aplikacji
var bgColor = color.NRGBA{R: 255, G: 240, B: 245, A: 255}

func (p PinkTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	// Ikony (w tym oko w polu hasła) i tekst
	if n == theme.ColorNameForeground {
		// Ciemna malina
		return color.NRGBA{R: 150, G: 30, B: 80, A: 255}
	}

	// Jasnoróżowy separator
	if n == theme.ColorNameSeparator {
		return focusPink
	}

	switch n {
	case theme.ColorNameBackground:
		return bgColor
	// TŁO WYSKAKUJĄCYCH OKIENEK (DIALOGÓW) - ZMIANA TUTAJ
	case theme.ColorNameOverlayBackground:
		return bgColor
	// (Opcjonalnie) tło rozwijanych menu, żeby też było spójne
	case theme.ColorNameMenuBackground:
		return bgColor
	case theme.ColorNameInputBackground:
		return color.White
	case theme.ColorNameInputBorder:
		return pinkColor
	case theme.ColorNameFocus:
		return focusPink
	case theme.ColorNameButton:
		return pinkColor
	case theme.ColorNamePrimary:
		return focusPink
	}

	return theme.DefaultTheme().Color(n, v)
}

func (p PinkTheme) Font(s fyne.TextStyle) fyne.Resource {
	return resources.GetFont()
}

func (p PinkTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (p PinkTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameText {
		return 16
	}
	return theme.DefaultTheme().Size(n)
}