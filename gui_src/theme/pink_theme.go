package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/gabbla05/KittyProtocol/gui_src/resources"
	"image/color"
)

type PinkTheme struct{}

var _ fyne.Theme = (*PinkTheme)(nil)

func (p PinkTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
    // Force Light Mode: ignorujemy 'v' (variant) i zawsze zwracamy jasne kolory
    
	if n == theme.ColorNameBackground {
		return color.NRGBA{R: 255, G: 240, B: 245, A: 255} 
	}
	if n == theme.ColorNameForeground {
		return color.NRGBA{R: 50, G: 50, B: 50, A: 255}
	}
	if n == theme.ColorNameInputBackground {
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	}
	// Kolor przycisków i akcentów: hot pink
	if n == theme.ColorNamePrimary {
		return color.NRGBA{R: 255, G: 105, B: 180, A: 255}
	}
	// Button styling
	if n == theme.ColorNameButton {
		return color.NRGBA{R: 255, G: 105, B: 180, A: 255}
	}
	// Hover state dla przycisków
	if n == theme.ColorNameHover {
		return color.NRGBA{R: 255, G: 80, B: 160, A: 255}
	}
	// Fokus dla input fields
	if n == theme.ColorNameFocus {
		return color.NRGBA{R: 255, G: 105, B: 180, A: 255}
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
	// Powiększenie tekstu dla lepszej czytelności
	if n == theme.SizeNameText {
		return 16
	}
	
	// Usunięto SizeNameInputHeight, ponieważ nie jest dostępna w Twojej wersji Fyne.
	// Wysokość pól możesz kontrolować bezpośrednio w kodzie widoków 
	// używając .SetMinSize(fyne.NewSize(x, 40)) w razie potrzeby.
	
	return theme.DefaultTheme().Size(n)
}