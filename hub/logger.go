// logger.go
// Lightweight structured logger for the Hub. Uses stdout and RFC3339 timestamps.
// This logger is intentionally simple to avoid external dependencies.

package hub

import (
	"fmt"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

var colorsEnabled = true

func log(level, msg string, args ...any) {
	timestamp := time.Now().Format(time.RFC3339)

	var color string
	if colorsEnabled {
		switch level {
		case "INFO":
			color = colorBlue
		case "WARN":
			color = colorYellow
		case "ERROR":
			color = colorRed
		default:
			color = colorReset
		}
	}

	// Without colors
	if !colorsEnabled {
		fmt.Printf("[%s] %s: %s\n",
			level, timestamp, fmt.Sprintf(msg, args...))
		return
	}

	// With colors
	fmt.Printf("%s[%s] %s: %s%s\n",
		color, level, timestamp, fmt.Sprintf(msg, args...), colorReset)
}

func logInfo(msg string, args ...any)  { log("INFO", msg, args...) }
func logWarn(msg string, args ...any)  { log("WARN", msg, args...) }
func logError(msg string, args ...any) { log("ERROR", msg, args...) }
