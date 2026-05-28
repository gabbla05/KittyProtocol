package views

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/gui_src/resources"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetMenuView(s *state.UIState) fyne.CanvasObject {
	// --- Logotyp w rogu ---
	ikona := canvas.NewImageFromResource(resources.LogoIkonaPng)
	ikona.SetMinSize(fyne.NewSize(60, 60))
	header := container.NewHBox(ikona, layout.NewSpacer())

	title := widget.NewLabel("Main Menu")
	title.TextStyle = fyne.TextStyle{Bold: true}

	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("Target username (e.g. bob)")

	secretEntry := widget.NewPasswordEntry()
	secretEntry.SetPlaceHolder("E2EE Shared Secret (min. 16 chars)")

	// --- Nasłuchiwanie powiadomień ---
	s.UI.OnMessage = func(msg string) {
		if strings.Contains(msg, "CHAT_REQUEST received") {
			dialog.ShowInformation("Incoming Chat!", msg, s.Window)
		} else if strings.Contains(msg, "refused") {
			dialog.ShowInformation("Chat Refused", msg, s.Window)
		}
	}

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
			dialog.ShowInformation("Error", "Enter target and valid secret (min 16 chars)!", s.Window)
			return
		}
		s.App.Secrets().Set(target, []byte(secret))
		s.App.Client().SetSharedSecretForPeer(target, []byte(secret))
		dialog.ShowInformation("Success", "Secret saved!", s.Window)
	})

	requestChatBtn := widget.NewButton("Request Chat", func() {
		target := targetEntry.Text
		if target == "" {
			dialog.ShowInformation("Error", "Enter target username!", s.Window)
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
		s.SwitchView(GetChatView(s, target))
	})

	acceptChatBtn := widget.NewButton("Accept Incoming Chat", func() {
		target := targetEntry.Text
		if target == "" {
			dialog.ShowInformation("Error", "Enter username!", s.Window)
			return
		}
		err := s.App.AcceptChat(target)
		if err != nil {
			dialog.ShowError(err, s.Window)
			return
		}
		s.SwitchView(GetChatView(s, target))
	})

	refuseChatBtn := widget.NewButton("Refuse Incoming Chat", func() {
		target := targetEntry.Text
		if target == "" {
			dialog.ShowInformation("Error", "Enter username!", s.Window)
			return
		}
		err := s.App.RefuseChat(target, "Declined via GUI")
		if err != nil {
			dialog.ShowError(err, s.Window)
			return
		}
		dialog.ShowInformation("Success", "Chat refused.", s.Window)
	})

	logoutBtn := widget.NewButton("Logout", func() {
		s.App.Client().SendBye()
		s.App.Client().Close()
		s.SwitchView(GetAuthView(s))
	})
	logoutBtn.Importance = widget.DangerImportance

	// --- Finalny układ ---
	menuContent := container.NewVBox(
		title,
		widget.NewLabel("1. Set target & E2EE secret:"),
		targetEntry,
		secretEntry,
		setSecretBtn,
		widget.NewSeparator(),
		widget.NewLabel("2. Actions:"),
		statusBtn,
		requestChatBtn,
		acceptChatBtn,
		refuseChatBtn,
		widget.NewSeparator(),
		logoutBtn,
	)

	return container.NewBorder(header, nil, nil, nil, container.NewPadded(menuContent))
}