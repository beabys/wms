package logger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLogger(t *testing.T) {
	l, err := NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)
	require.NotNil(t, l)
	assert.NotNil(t, l.log)
}

func TestNewZapLoggerCustomPaths(t *testing.T) {
	l, err := NewZapLogger([]string{"stdout"}, []string{"stderr"}, zapcore.InfoLevel)
	require.NoError(t, err)
	require.NotNil(t, l)
}

func TestZapLoggerGetLogger(t *testing.T) {
	l, err := NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)
	zapLog := l.GetLogger()
	_, ok := zapLog.(*zap.Logger)
	assert.True(t, ok)
}

func TestZapLoggerDebug(t *testing.T) {
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	l := &ZapLogger{log: zap.New(core)}

	l.Debug("debug msg", LogField{Key: "key", Value: "val"})

	var result map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "debug msg", result["msg"])
	assert.Equal(t, "debug", result["level"])
	assert.Equal(t, "val", result["key"])
}

func TestZapLoggerInfo(t *testing.T) {
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	l := &ZapLogger{log: zap.New(core)}

	l.Info("info msg", LogField{Key: "count", Value: 42})

	var result map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "info msg", result["msg"])
	assert.Equal(t, float64(42), result["count"])
}

func TestZapLoggerWarn(t *testing.T) {
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.WarnLevel)
	l := &ZapLogger{log: zap.New(core)}

	l.Warn("warn msg")
	assert.Contains(t, buf.String(), "warn msg")
}

func TestZapLoggerError(t *testing.T) {
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.ErrorLevel)
	l := &ZapLogger{log: zap.New(core)}

	l.Error("error msg", assert.AnError, LogField{Key: "extra", Value: "info"})

	var result map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "error msg", result["msg"])
	assert.Equal(t, "info", result["extra"])
	assert.Contains(t, result["error"], assert.AnError.Error())
}

func TestZapLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.WarnLevel)
	l := &ZapLogger{log: zap.New(core)}

	l.Debug("should be hidden")
	l.Info("also hidden")
	l.Warn("visible warn")

	assert.NotContains(t, buf.String(), "should be hidden")
	assert.NotContains(t, buf.String(), "also hidden")
	assert.Contains(t, buf.String(), "visible warn")
}

func TestZapLoggerWith(t *testing.T) {
	l, err := NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)
	assert.NotNil(t, l)

	withLogger := l.With(LogField{Key: "with_key", Value: "with_val"})
	assert.NotNil(t, withLogger)
}

func TestZapLoggerWithEmpty(t *testing.T) {
	l, err := NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)

	withLogger := l.With()
	assert.NotNil(t, withLogger)
}

func TestLogFieldsToZapFields(t *testing.T) {
	lfs := []LogField{
		{Key: "str", Value: "hello"},
		{Key: "int", Value: 42},
		{Key: "bool", Value: true},
		{Key: "float", Value: 3.14},
		{Key: "bytes", Value: []byte("data")},
		{Key: "duration", Value: 5 * time.Second},
		{Key: "time", Value: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	fields := logFieldsToZapFields(lfs)
	assert.Len(t, fields, 7)
}

func TestZapLoggerFatal(t *testing.T) {
	// Fatal normally calls os.Exit(1). Use WriteThenPanic hook to test.
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.FatalLevel)
	fatalLogger := &ZapLogger{log: zap.New(core, zap.WithFatalHook(zapcore.WriteThenPanic))}

	assert.Panics(t, func() {
		fatalLogger.Fatal("fatal msg", LogField{Key: "reason", Value: "test"})
	})

	assert.Contains(t, buf.String(), "fatal msg")
	assert.Contains(t, buf.String(), "test")
}

func TestInterceptorLoggerWithZap(t *testing.T) {
	l, err := NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)

	interceptor := InterceptorLogger(l.GetLogger().(*zap.Logger))
	assert.NotNil(t, interceptor)
}

func TestInterceptorLoggerPanicOnUnsupported(t *testing.T) {
	assert.Panics(t, func() {
		InterceptorLogger("unsupported")
	})
}

func TestZapLoggerMultipleFields(t *testing.T) {
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	l := &ZapLogger{log: zap.New(core)}

	l.Info("test",
		LogField{Key: "k1", Value: "v1"},
		LogField{Key: "k2", Value: "v2"},
	)

	var result map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "v1", result["k1"])
	assert.Equal(t, "v2", result["k2"])
}
