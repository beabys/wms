package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockLogger implements Logger interface for testing.
type mockLogger struct {
	debugMsgs []string
	infoMsgs  []string
	warnMsgs  []string
	errorMsgs []string
	fatalMsgs []string
}

func (m *mockLogger) GetLogger() any      { return m }
func (m *mockLogger) Debug(s string, lf ...LogField)   { m.debugMsgs = append(m.debugMsgs, s) }
func (m *mockLogger) Info(s string, lf ...LogField)    { m.infoMsgs = append(m.infoMsgs, s) }
func (m *mockLogger) Warn(s string, lf ...LogField)    { m.warnMsgs = append(m.warnMsgs, s) }
func (m *mockLogger) Error(s string, err error, lf ...LogField) { m.errorMsgs = append(m.errorMsgs, s) }
func (m *mockLogger) Fatal(s string, lf ...LogField)   { m.fatalMsgs = append(m.fatalMsgs, s) }

func TestMockImplementsLogger(t *testing.T) {
	var _ Logger = (*mockLogger)(nil)
}

func TestMockLoggerDebug(t *testing.T) {
	m := &mockLogger{}
	m.Debug("debug message", LogField{Key: "k", Value: "v"})
	assert.Len(t, m.debugMsgs, 1)
	assert.Equal(t, "debug message", m.debugMsgs[0])
}

func TestMockLoggerInfo(t *testing.T) {
	m := &mockLogger{}
	m.Info("info message")
	assert.Len(t, m.infoMsgs, 1)
	assert.Equal(t, "info message", m.infoMsgs[0])
}

func TestMockLoggerWarn(t *testing.T) {
	m := &mockLogger{}
	m.Warn("warn message")
	assert.Len(t, m.warnMsgs, 1)
	assert.Equal(t, "warn message", m.warnMsgs[0])
}

func TestMockLoggerError(t *testing.T) {
	m := &mockLogger{}
	m.Error("error message", assert.AnError)
	assert.Len(t, m.errorMsgs, 1)
	assert.Equal(t, "error message", m.errorMsgs[0])
}

func TestMockLoggerFatal(t *testing.T) {
	m := &mockLogger{}
	// Note: Fatal would call os.Exit in real impl, but our mock doesn't
	m.Fatal("fatal message")
	assert.Len(t, m.fatalMsgs, 1)
	assert.Equal(t, "fatal message", m.fatalMsgs[0])
}

func TestMockLoggerFields(t *testing.T) {
	m := &mockLogger{}
	m.Info("test", LogField{Key: "key1", Value: "val1"}, LogField{Key: "key2", Value: 2})
	assert.Len(t, m.infoMsgs, 1)
}
