package app

// import (
// 	"fmt"
// 	"strings"
// )

// func (app *App) RunMainMenu() {
// 	for {
// 		fmt.Println("Dostępne komendy:")
// 		fmt.Println("  /status <user>")
// 		fmt.Println("  /chat <user>")
// 		fmt.Println("  /quit")

// 		cmd := app.ui.ReadLine()

// 		switch {
// 		case strings.HasPrefix(cmd, "/status "):
// 			app.handleStatus()
// 		case strings.HasPrefix(cmd, "/chat "):
// 			app.handleChat()
// 		case cmd == "/quit":
// 			app.client.SendBye()
// 			return
// 		}
// 	}
// }
