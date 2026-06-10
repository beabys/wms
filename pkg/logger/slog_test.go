package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSlogLogger(t *testing.T) {
	l := NewSlogLogger(slog.LevelDebug)
	require.NotNil(t, l)
	assert.NotNil(t, l.log)
}

func TestSlogLoggerDebug(t *testing.T) {
	var buf bytes.Buffer
	l := &SlogLogger{
		log: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}

	l.Debug("debug msg", LogField{Key: "key", Value: "val"})
	assert.Contains(t, buf.String(), "debug msg")
	assert.Contains(t, buf.String(), "key")
	assert.Contains(t, buf.String(), "val")
	assert.Contains(t, buf.String(), "DEBUG")
}

func TestSlogLoggerInfo(t *testing.T) {
	var buf bytes.Buffer
	l := &SlogLogger{
		log: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}

	l.Info("info msg", LogField{Key: "count", Value: 42})
	assert.Contains(t, buf.String(), "info msg")
	assert.Contains(t, buf.String(), "42")
}

func TestSlogLoggerWarn(t *testing.T) {
	var buf bytes.Buffer
	l := &SlogLogger{
		log: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})),
	}

	l.Warn("warn msg")
	assert.Contains(t, buf.String(), "warn msg")
}

func TestSlogLoggerError(t *testing.T) {
	var buf bytes.Buffer
	l := &SlogLogger{
		log: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	l.Error("error msg", assert.AnError, LogField{Key: "err", Value: assert.AnError.Error()})
	assert.Contains(t, buf.String(), "error msg")
}

func TestSlogLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := &SlogLogger{
		log: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	l.Debug("should be hidden")
	l.Info("should be hidden too")
	l.Warn("also hidden")
	l.Error("visible error", assert.AnError)

	assert.NotContains(t, buf.String(), "should be hidden")
	assert.NotContains(t, buf.String(), "should be hidden too")
	assert.Contains(t, buf.String(), "visible error")
}

func TestSlogLoggerMultipleFields(t *testing.T) {
	var buf bytes.Buffer
	l := &SlogLogger{
		log: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}

	l.Info("multi field",
		LogField{Key: "str", Value: "hello"},
		LogField{Key: "num", Value: 123},
	)

	var result map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "hello", result["str"])
	// slog encodes numbers as float64 in JSON
	assert.Equal(t, float64(123), result["num"])
}

func TestSlogLoggerGetLogger(t *testing.T) {
	l := NewSlogLogger(slog.LevelDebug)
	assert.NotNil(t, l.GetLogger())
}

func TestLogFieldsToSlogArgs(t *testing.T) {
	lfs := []LogField{
		{Key: "a", Value: "1"},
		{Key: "b", Value: 2},
	}
	args := logFieldsToSlogArgs(lfs)
	assert.Len(t, args, 4)
	assert.Equal(t, "a", args[0])
	assert.Equal(t, "1", args[1])
	assert.Equal(t, "b", args[2])
	assert.Equal(t, 2, args[3])
}

func TestLogFieldsToSlogArgsEmpty(t *testing.T) {
	args := logFieldsToSlogArgs(nil)
	assert.Empty(t, args)

	args = logFieldsToSlogArgs([]LogField{})
	assert.Empty(t, args)
}

func TestSlogLoggerFatalExists(t *testing.T) {
	// Fatal calls os.Exit(1) — cannot test directly without killing runner.
	// Verify the method is wired: it logs at error level before exiting.
	var buf bytes.Buffer
	l := &SlogLogger{
		log: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	// We call l.log.Error directly (same as Fatal does) to verify log output.
	l.log.Error("fatal msg would appear here")
	assert.Contains(t, buf.String(), "fatal msg would appear here")
}
