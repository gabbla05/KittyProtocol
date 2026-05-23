package main

import (
	"fmt"
	"time"
)

func logInfo(msg string, args ...any) {
	fmt.Printf("[INFO] %s: %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}

func logWarn(msg string, args ...any) {
	fmt.Printf("[WARN] %s: %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}

func logError(msg string, args ...any) {
	fmt.Printf("[ERROR] %s: %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}
