package api

import (
	"fmt"
	"sync"
)

type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

type Logger interface {
	Log(level LogLevel, msg string)
}

type defaultLogger struct{}

func (defaultLogger) Log(level LogLevel, msg string) {
	// Domyślnie nic — API jest UI-agnostic.
}

type logManager struct {
	mu     sync.Mutex
	logger Logger
}

var globalLogger = &logManager{
	logger: defaultLogger{},
}

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
