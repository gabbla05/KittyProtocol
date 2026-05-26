package api

import (
	"fmt"
	"sync"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

// Logger is a minimal logging interface used by KittyClient.
//
// UI layers (CLI, GUI, Wails frontend) are expected to provide their own
// implementation and install it via SetLogger.
type Logger interface {
	Log(level LogLevel, msg string)
}

type defaultLogger struct{}

// Log implements Logger but intentionally discards all messages.
// This keeps the API layer UI-agnostic by default.
func (defaultLogger) Log(level LogLevel, msg string) {}

type logManager struct {
	mu     sync.Mutex
	logger Logger
}

var globalLogger = &logManager{
	logger: defaultLogger{},
}

// SetLogger installs a process-wide logger used by the client API.
// It is safe to call from multiple goroutines, but typically configured
// once during application startup.
func SetLogger(l Logger) {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()
	globalLogger.logger = l
}

func log(level LogLevel, format string, args ...any) {
	globalLogger.mu.Lock()
	l := globalLogger.logger
	globalLogger.mu.Unlock()

	l.Log(level, fmt.Sprintf(format, args...))
}
