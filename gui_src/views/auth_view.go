package views

import (
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas" // Dodane dla obsługi obrazków
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gabbla05/KittyProtocol/client/app"
	"github.com/gabbla05/KittyProtocol/gui_src/resources"
	"github.com/gabbla05/KittyProtocol/gui_src/state"
)

const authTimeout = 5 * time.Second

func GetAuthView(s *state.UIState) fyne.CanvasObject {
	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Username")

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Password")

	// --- Logo z napisem ---
	logo := canvas.NewImageFromResource(resources.LogoZNapisemPng)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(300, 150)) // Dostosuj rozmiar w razie potrzeby

	loginBtn := widget.NewButton("Login", func() {
		username := userEntry.Text
		password := passEntry.Text

		if username == "" || password == "" {
			dialog.ShowInformation("Error", "Enter username and password!", s.Window)
			return
		}

		hubAddr := os.Getenv("KITTY_HUB_ADDR")
		if hubAddr == "" {
			hubAddr = "127.0.0.1:9999"
		}

		go func() {
			if err := s.Client.Connect(hubAddr); err != nil {
				dialog.ShowError(err, s.Window)
				return
			}

			disconnectedChan := make(chan struct{})
			go s.Client.StartReceiverLoop(disconnectedChan)

			select {
			case helloRes := <-s.Client.HelloResult():
				if !helloRes.OK {
					dialog.ShowInformation("Error", "[Client] HELLO failed: "+helloRes.Desc, s.Window)
					s.Client.Close()
					return
				}
			case <-time.After(authTimeout):
				dialog.ShowInformation("Error", "[Client] HELLO timeout", s.Window)
				s.Client.Close()
				return
			}

			if err := s.Client.SendAuth(username, password); err != nil {
				dialog.ShowError(err, s.Window)
				return
			}

			select {
			case authRes := <-s.Client.AuthResult():
				if !authRes.OK {
					dialog.ShowInformation("Error", "[Client] AUTH error: "+authRes.Desc, s.Window)
					s.Client.Close()
					return
				}

				go s.Client.StartPingLoop()
				log.Println("User authorized, switching to main menu!")

				guiUI := &state.GuiUI{}
				s.UI = guiUI
				s.App = app.NewApp(s.Client, guiUI, disconnectedChan)
				s.App.InitSecretStoreForUser(s.Client.User(), []byte(password))

				s.Window.Canvas().Refresh(s.Window.Content())
				s.SwitchView(GetMenuView(s))

			case <-time.After(authTimeout):
				dialog.ShowInformation("Error", "AUTH timeout", s.Window)
				s.Client.Close()
				return
			}
		}()
	})

	registerBtn := widget.NewButton("Register", func() {
		username := userEntry.Text
		password := passEntry.Text

		if username == "" || password == "" {
			dialog.ShowInformation("Error", "Enter username and password for registration!", s.Window)
			return
		}

		hubAddr := os.Getenv("KITTY_HUB_ADDR")
		if hubAddr == "" {
			hubAddr = "127.0.0.1:9999"
		}

		go func() {
			if err := s.Client.Connect(hubAddr); err != nil {
				dialog.ShowError(err, s.Window)
				return
			}

			disconnectedChan := make(chan struct{})
			go s.Client.StartReceiverLoop(disconnectedChan)

			select {
			case helloRes := <-s.Client.HelloResult():
				if !helloRes.OK {
					dialog.ShowInformation("Error", "[Client] HELLO failed: "+helloRes.Desc, s.Window)
					s.Client.Close()
					return
				}
			case <-time.After(authTimeout):
				dialog.ShowInformation("Error", "[Client] HELLO timeout", s.Window)
				s.Client.Close()
				return
			}

			if err := s.Client.SendRegister(username, password); err != nil {
				dialog.ShowError(err, s.Window)
				return
			}

			select {
			case regRes := <-s.Client.RegisterResult():
				if !regRes.OK {
					dialog.ShowInformation("Error", "[Client] REGISTER error: "+regRes.Desc, s.Window)
					s.Client.Close()
					return
				}

				dialog.ShowInformation("Success", "[Client] REGISTER OK — you can log in now.", s.Window)
				s.Client.Close()

			case <-time.After(authTimeout):
				dialog.ShowInformation("Error", "REGISTER timeout", s.Window)
				s.Client.Close()
				return
			}
		}()
	})

	// STYLOWANY UKŁAD: Logo na górze, potem pola formularza
	return container.NewPadded(container.NewCenter(container.NewVBox(
		logo,
		layout.NewSpacer(),
		userEntry,
		passEntry,
		loginBtn,
		widget.NewSeparator(),
		registerBtn,
	)))
}