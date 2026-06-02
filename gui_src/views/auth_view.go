package views

import (
	"errors"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gabbla05/KittyProtocol/client"
	"github.com/gabbla05/KittyProtocol/client/api"
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

	logo := canvas.NewImageFromResource(resources.LogoZNapisemPng)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(300, 150))

	var isConnected bool
	var disconnectedChan chan struct{}

	ensureConnection := func(hubAddr string) error {
		if isConnected {
			return nil
		}

		if err := s.Client.Connect(hubAddr); err != nil {
			s.Client = api.NewKittyClient()
			return err
		}

		disconnectedChan = make(chan struct{})
		go s.Client.StartReceiverLoop(disconnectedChan)

		select {
		case helloRes := <-s.Client.HelloResult():
			if !helloRes.OK {
				s.Client.Close()
				s.Client = api.NewKittyClient()
				return errors.New("HELLO failed: " + helloRes.Desc)
			}
		case <-time.After(authTimeout):
			s.Client.Close()
			s.Client = api.NewKittyClient()
			return errors.New("HELLO timeout")
		}

		isConnected = true
		return nil
	}

	loginBtn := widget.NewButton("Login", func() {

		client.LoadEnv()
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
			if err := ensureConnection(hubAddr); err != nil {
				dialog.ShowError(err, s.Window)
				return
			}

			if err := s.Client.SendAuth(username, password); err != nil {
				dialog.ShowError(err, s.Window)
				s.Client.Close()
				s.Client = api.NewKittyClient()
				isConnected = false
				return
			}

			select {
			case authRes := <-s.Client.AuthResult():
				if !authRes.OK {
					dialog.ShowInformation("Error", "[Client] AUTH error: "+authRes.Desc, s.Window)
					return
				}

				fyne.Do(func() {
					log.Println("User authorized, switching to main menu!")

					guiUI := &state.GuiUI{State: s}
					s.UI = guiUI
					s.App = app.NewApp(s.Client, guiUI, disconnectedChan)
					s.App.InitSecretStoreForUser(s.Client.User(), []byte(password))

					// --- PODPIĘCIE ZDARZEŃ CZATU ---
					go func() {
						for ev := range s.App.ChatEvents() {
							switch ev.Type {

							case "request":
								fyne.Do(func() {
									cnf := dialog.NewConfirm(
										"Incoming Chat Request",
										"User "+ev.From+" wants to chat with you.",
										func(accept bool) {
											if accept {
												if err := s.App.AcceptChat(ev.From); err != nil {
													dialog.ShowError(err, s.Window)
													return
												}
												s.SwitchToChat(ev.From)
											} else {
												s.App.RefuseChat(ev.From, "User refused via GUI")
											}
										},
										s.Window,
									)
									cnf.SetConfirmText("Accept")
									cnf.SetDismissText("Refuse")
									cnf.Show()
								})

							case "accept":
								fyne.Do(func() {
									s.SwitchToChat(ev.From)
								})

							case "refuse":
								fyne.Do(func() {
									dialog.ShowInformation("Chat Refused", ev.From+" refused: "+ev.Reason, s.Window)
								})

							case "end":
								fyne.Do(func() {
									dialog.ShowInformation("Chat Ended", ev.Reason, s.Window)
									s.SwitchView(GetMenuView(s))
								})

							case "message":
								if s.UI.OnMessage != nil {
									fyne.Do(func() {
										s.UI.OnMessage(ev.From, ev.Text)
									})
								}
							}
						}
					}()

					s.Window.Canvas().Refresh(s.Window.Content())
					s.SwitchView(GetMenuView(s))
				})

			case <-time.After(authTimeout):
				dialog.ShowInformation("Error", "AUTH timeout", s.Window)
				s.Client.Close()
				s.Client = api.NewKittyClient()
				isConnected = false
			}
		}()
	})

	registerBtn := widget.NewButton("Register", func() {
		client.LoadEnv()

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
			if err := ensureConnection(hubAddr); err != nil {
				dialog.ShowError(err, s.Window)
				return
			}

			if err := s.Client.SendRegister(username, password); err != nil {
				dialog.ShowError(err, s.Window)
				s.Client.Close()
				s.Client = api.NewKittyClient()
				isConnected = false
				return
			}

			select {
			case regRes := <-s.Client.RegisterResult():
				if !regRes.OK {
					dialog.ShowInformation("Error", "[Client] REGISTER error: "+regRes.Desc, s.Window)
					return
				}
				dialog.ShowInformation("Success", "[Client] REGISTER OK — you can log in now.", s.Window)

			case <-time.After(authTimeout):
				dialog.ShowInformation("Error", "REGISTER timeout", s.Window)
				s.Client.Close()
				s.Client = api.NewKittyClient()
				isConnected = false
			}
		}()
	})

	formContent := container.New(&formLayout{},
		container.NewCenter(logo),
		userEntry,
		passEntry,
		loginBtn,
		widget.NewSeparator(),
		registerBtn,
	)

	return container.NewCenter(formContent)
}
