package logger_test

import (
	"bytes"
	"testing"

	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	const logMessage = "log message"

	t.Run("log with exact level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New("TEST: ", "warn", out)
		logg.Warn(logMessage)

		require.Contains(t, out.String(), "[WARN] "+logMessage)
	})

	t.Run("log with higher level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New("TEST: ", "info", out)
		logg.Warn(logMessage)

		require.Contains(t, out.String(), "[WARN] "+logMessage)
	})

	t.Run("log with lower level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New("TEST: ", "error", out)
		logg.Warn(logMessage)

		require.Empty(t, out.String())
	})

	t.Run("log info message at info level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New("TEST: ", "info", out)
		logg.Info(logMessage)

		require.Contains(t, out.String(), "[INFO] "+logMessage)
	})

	t.Run("log debug message at debug level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New("TEST: ", "debug", out)
		logg.Debug(logMessage)

		require.Contains(t, out.String(), "[DEBUG] "+logMessage)
	})

	t.Run("log debug message at info level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New("TEST: ", "info", out)
		logg.Debug(logMessage)

		require.Empty(t, out.String())
	})
}
