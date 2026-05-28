package resources

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed assets/images/logo_ikona.png
var logoIkonaData []byte
var LogoIkonaPng = fyne.NewStaticResource("logo_ikona.png", logoIkonaData)

//go:embed assets/images/logo_z_napisem.png
var logoZNapisemData []byte
var LogoZNapisemPng = fyne.NewStaticResource("logo_z_napisem.png", logoZNapisemData)

//go:embed assets/fonts/Quicksand.ttf
var fontData []byte
var QuicksandVariableFontWghtTtf = fyne.NewStaticResource("Quicksand-VariableFont_wght.ttf", fontData)

func GetLogoIkona() fyne.Resource {
	return fyne.NewStaticResource("logo_ikona.png", logoIkonaData)
}

func GetLogoZNapisem() fyne.Resource {
	return fyne.NewStaticResource("logo_z_napisem.png", logoZNapisemData)
}

func GetFont() fyne.Resource {
	return fyne.NewStaticResource("Quicksand.ttf", fontData)
}