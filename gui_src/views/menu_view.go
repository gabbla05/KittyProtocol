package views

import (
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/gui_src/resources"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetMenuView(s *state.UIState) fyne.CanvasObject {
	// --- Logo ---
	logo := canvas.NewImageFromResource(resources.LogoZNapisemPng)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(300, 150))

	// Tytuł
	titleText := canvas.NewText("Main Menu", color.NRGBA{R: 255, G: 105, B: 180, A: 255})
	titleText.TextSize = 32
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	// --- Pola wejściowe ---
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("Target username")

	secretEntry := widget.NewPasswordEntry()
	secretEntry.SetPlaceHolder("E2EE Shared Secret")

	// --- Przyciski ---
	statusBtn := widget.NewButton("Check Status", func() {
		target := targetEntry.Text
		if target == "" {
			dialog.ShowInformation("Error", "Enter target username!", s.Window)
			return
		}
		s.App.Client().OnStatus(func(user, status string) {
			dialog.ShowInformation("Status", "User "+user+" is: "+status, s.Window)
		})
		_ = s.App.Client().SendGetStatus(target)
	})

	setSecretBtn := widget.NewButton("Save Shared Secret", func() {
		target := targetEntry.Text
		secret := secretEntry.Text
		if target == "" || len(secret) < 16 {
			dialog.ShowInformation("Error", "Secret must be min 16 chars!", s.Window)
			return
		}
		s.App.Secrets().Set(target, []byte(secret))
		s.App.Client().SetSharedSecretForPeer(target, []byte(secret))
		dialog.ShowInformation("Success", "Secret saved!", s.Window)
	})

	requestChatBtn := widget.NewButton("Request Chat", func() {
		target := targetEntry.Text
		secret := secretEntry.Text
		if target == "" || secret == "" || len(secret) < 16 {
			dialog.ShowInformation("Error", "Invalid data: Target username and 16-character Shared Secret are required.", s.Window)
			return
		}
		if !s.App.Client().HasSharedSecret(target) {
			dialog.ShowInformation("Error", "No shared secret saved for "+target+". Fill Target username and Shared Secret before requesting chat.", s.Window)
			return
		}

		err := s.App.StartChatRequest(target)
		if err != nil {
			if errors.Is(err, api.ErrNoSharedSecret) {
				dialog.ShowInformation("Error", "No secret saved for "+target, s.Window)
			} else {
				dialog.ShowError(err, s.Window)
			}
			return
		}

		// 🔥 NIE przełączamy widoku tutaj!
		// Czekamy na event "accept" z ChatEvents()
		dialog.ShowInformation("Chat Request Sent", "Waiting for "+target+" to accept...", s.Window)
	})

	logoutBtn := widget.NewButton("Logout", func() {
		s.App.Client().SendBye()
		s.App.Client().Close()
		s.SwitchView(GetAuthView(s))
	})

	// --- Finalny układ ---
	menuContent := container.New(&formLayout{},
		container.NewCenter(logo),
		container.NewCenter(titleText),
		widget.NewLabel("1. Configuration:"),
		targetEntry,
		secretEntry,
		setSecretBtn,
		widget.NewSeparator(),
		widget.NewLabel("2. Actions:"),
		statusBtn,
		requestChatBtn,
		widget.NewSeparator(),
		logoutBtn,
	)

	return container.NewCenter(container.NewPadded(menuContent))
}
