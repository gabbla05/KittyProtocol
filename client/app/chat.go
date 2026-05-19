package app

// func (app *App) RunChatSession(target string) {
// 	app.client.SetTarget(target)

// 	secret := app.ui.ReadSharedSecret()
// 	app.client.SetSharedSecret(secret)

// 	for {
// 		msg := app.ui.ReadLine()

// 		switch {
// 		case msg == "/quit":
// 			app.client.SendBye()
// 			return
// 		case msg == "/replay":
// 			app.client.Replay()
// 		default:
// 			app.client.SendMessage(msg)
// 		}
// 	}
// }
