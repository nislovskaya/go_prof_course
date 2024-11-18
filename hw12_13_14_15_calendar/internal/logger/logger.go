package logger

import (
	"io"
	"log"
)

type Logger struct {
	logger *log.Logger
	level  int
}

const (
	LevelError = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

func New(prefix, level string, output io.Writer) *Logger {
	var logLevel int
	switch level {
	case "error":
		logLevel = LevelError
	case "warn":
		logLevel = LevelWarn
	case "info":
		logLevel = LevelInfo
	case "debug":
		logLevel = LevelDebug
	default:
		logLevel = LevelInfo
	}

	return &Logger{
		logger: log.New(output, prefix, log.Ldate|log.Ltime|log.Lshortfile),
		level:  logLevel,
	}
}

func (l *Logger) Error(msg string) {
	if l.level < LevelError {
		return
	}
	l.logger.Printf("[ERROR] %s", msg)
}

func (l *Logger) Warn(msg string) {
	if l.level < LevelWarn {
		return
	}
	l.logger.Printf("[WARN] %s", msg)
}

func (l *Logger) Info(msg string) {
	if l.level < LevelInfo {
		return
	}
	l.logger.Printf("[INFO] %s", msg)
}

func (l *Logger) Debug(msg string) {
	if l.level < LevelDebug {
		return
	}
	l.logger.Printf("[DEBUG] %s", msg)
}
