package ui_cli

import (
	"fmt"

	"github.com/gabbla05/KittyProtocol/client/api"
)

// CliLogger is a simple stdout logger used by the CLI UI.
// It implements api.Logger and is installed via api.SetLogger().
type CliLogger struct{}

// Log prints a formatted log message based on severity.
func (CliLogger) Log(level api.LogLevel, msg string) {
	switch level {
	case api.LogError:
		fmt.Printf("[ERROR] %s\n", msg)
	case api.LogWarn:
		fmt.Printf("[WARN] %s\n", msg)
	case api.LogInfo:
		fmt.Printf("[INFO] %s\n", msg)
	case api.LogDebug:
		fmt.Printf("[DEBUG] %s\n", msg)
	}
}
