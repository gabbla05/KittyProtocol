package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabbla05/KittyProtocol/client/app"
)

func main() {
	dir := filepath.Dir(app.DefaultSecretStorePath())
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println("Katalog nie istnieje:", dir)
	} else {
		fmt.Println("Katalog istnieje:", dir)
	}
}
