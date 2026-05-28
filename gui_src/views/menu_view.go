package views

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gabbla05/KittyProtocol/client/api"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

func GetMenuView(s *state.UIState) fyne.CanvasObject {
	title := widget.NewLabel("Main Menu")
	title.TextStyle = fyne.TextStyle{Bold: true}

	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("Target username (e.g. bob)")

	secretEntry := widget.NewPasswordEntry()
	secretEntry.SetPlaceHolder("E2EE Shared Secret (min. 16 chars)")

	// --- 1. NASŁUCHIWANIE W MENU ---
	// Złapmy powiadomienia z tła (np. o zaproszeniu), żeby wyświetlić powiadomienie
	// bezpośrednio w GUI, zamiast szukać go w czarnej konsoli!
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

		err := s.App.Client().SendGetStatus(target)
		if err != nil {
			dialog.ShowError(err, s.Window)
		}
	})

	setSecretBtn := widget.NewButton("Save Shared Secret", func() {
		target := targetEntry.Text
		secret := secretEntry.Text
		if target == "" || secret == "" {
			dialog.ShowInformation("Error", "Enter target and secret!", s.Window)
			return
		}
		if len(secret) < 16 {
			dialog.ShowInformation("Error", "Secret must be at least 16 characters long!", s.Window)
			return
		}

		s.App.Secrets().Set(target, []byte(secret))
		s.App.Client().SetSharedSecretForPeer(target, []byte(secret))
		dialog.ShowInformation("Success", "Shared secret saved for "+target, s.Window)
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
				dialog.ShowInformation("Error", "No shared secret with "+target+"! Please save it first.", s.Window)
			} else {
				dialog.ShowError(err, s.Window)
			}
			return
		}

		dialog.ShowInformation("Chat Request", "CHAT_REQUEST sent to "+target+".\nWait for them to accept, then you can go to Chat.", s.Window)
		s.SwitchView(GetChatView(s, target))
	})

	acceptChatBtn := widget.NewButton("Accept Incoming Chat", func() {
		target := targetEntry.Text
		if target == "" {
			dialog.ShowInformation("Error", "Enter the username of the person inviting you!", s.Window)
			return
		}

		err := s.App.AcceptChat(target)
		if err != nil {
			dialog.ShowError(err, s.Window)
			return
		}

		dialog.ShowInformation("Success", "Chat accepted with "+target, s.Window)
		s.SwitchView(GetChatView(s, target))
	})

	// --- 2. NOWY GUZIK ODPOWIADAJĄCY KOMENDZIE /refuse ---
	refuseChatBtn := widget.NewButton("Refuse Incoming Chat", func() {
		target := targetEntry.Text
		if target == "" {
			dialog.ShowInformation("Error", "Enter the username of the person inviting you!", s.Window)
			return
		}

		// TUTAJ POPRAWKA: Dodaliśmy drugi argument (powód odrzucenia)
		err := s.App.RefuseChat(target, "Declined via GUI")
		if err != nil {
			dialog.ShowError(err, s.Window)
			return
		}

		dialog.ShowInformation("Success", "Chat request from "+target+" has been refused.", s.Window)
	})

	logoutBtn := widget.NewButton("Logout", func() {
		s.App.Client().SendBye()
		s.App.Client().Close()
		s.SwitchView(GetAuthView(s))
	})
	logoutBtn.Importance = widget.DangerImportance

	return container.NewVBox(
		title,
		widget.NewLabel("1. Set target (and E2EE secret if chatting):"),
		targetEntry,
		secretEntry,
		setSecretBtn,
		widget.NewSeparator(),
		widget.NewLabel("2. Actions:"),
		statusBtn,
		requestChatBtn,
		acceptChatBtn,
		refuseChatBtn, // Nasz dodany guzik!
		widget.NewSeparator(),
		logoutBtn,
	)
}