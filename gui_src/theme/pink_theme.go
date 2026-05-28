package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/gabbla05/KittyProtocol/gui_src/resources"
	"image/color"
)

type PinkTheme struct{}

func (p PinkTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if n == theme.ColorNameBackground { return color.NRGBA{R: 255, G: 240, B: 245, A: 255} }
	if n == theme.ColorNamePrimary { return color.NRGBA{R: 255, G: 105, B: 180, A: 255} }
	return theme.DefaultTheme().Color(n, v)
}

func (p PinkTheme) Font(s fyne.TextStyle) fyne.Resource {
	// Używamy zmiennej wygenerowanej przez fyne bundle
	return resources.QuicksandVariableFontWghtTtf
}

func (p PinkTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (p PinkTheme) Size(n fyne.ThemeSizeName) float32 { return theme.DefaultTheme().Size(n) }