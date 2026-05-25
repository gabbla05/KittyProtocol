package ui_cli

import (
	"fmt"

	"github.com/gabbla05/KittyProtocol/client/api"
)

type CliLogger struct{}

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
