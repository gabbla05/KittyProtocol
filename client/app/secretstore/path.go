package secretstore

import (
	"os"
	"path/filepath"
)

// PathForUser returns ~/.kitty/<kittyUser>/secrets.json.enc
// or ./kitty/<kittyUser>/secrets.json.enc if $HOME is unavailable.
func PathForUser(kittyUser string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".", "kitty", kittyUser, "secrets.json.enc")
	}
	return filepath.Join(home, ".kitty", kittyUser, "secrets.json.enc")
}
