package ui_commands

import (
	"bytes"
	"os"
	"strings"

	"github.com/gabbla05/KittyProtocol/client/app"
)

// Secret handles the logic for the "/secret <user> [file:<path>]" command.
// It supports two modes:
//  1. Loading the shared secret from a file (file:<path>)
//  2. Using a secret provided by the UI layer (e.g. typed by the user)
//
// This function performs no I/O other than reading a file when explicitly
// requested. It does not print anything and does not interact with stdin.
// The UI layer is responsible for collecting the secret from the user.
//
// After obtaining the secret, it delegates key derivation and storage to
// the App layer and SecretStore.
//
// Returns:
//   - string: a user-facing message
//   - error: non-nil if key derivation or storage failed
func Secret(line string, secretInput []byte, a *app.App) (string, error) {
	args := strings.Fields(line)
	if len(args) < 2 {
		return "Usage: /secret <user> [file:<path>]", nil
	}

	user := strings.ToLower(args[1])
	var secret []byte

	// Load secret from file if requested
	if len(args) == 3 && strings.HasPrefix(args[2], "file:") {
		path := strings.TrimPrefix(args[2], "file:")
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		secret = bytes.TrimSpace(data)
	} else {
		// Secret provided by UI (CLI or GUI)
		secret = secretInput
	}

	// Derive keys and store secret
	if err := a.Client().SetSharedSecretForPeer(user, secret); err != nil {
		return "", err
	}
	if err := a.Secrets().Set(user, secret); err != nil {
		return "", err
	}

	return "Shared secret configured for " + user, nil
}
