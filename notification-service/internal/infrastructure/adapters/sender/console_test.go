package sender

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/beabys/wms/notification-service/internal/application/notification/command"
)

func newTestLogger(buf *bytes.Buffer) *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(buf),
		zapcore.DebugLevel,
	)
	return zap.New(core)
}

func TestConsoleSender_SendEmail(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)
	sender := NewConsoleSender(logger)

	cmd := command.SendEmailCommand{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	result, err := sender.SendEmail(context.Background(), cmd)
	require.NoError(t, err)
	assert.NotEmpty(t, result.MessageID)
	assert.Equal(t, "sent", result.Status)
	assert.True(t, strings.HasPrefix(result.MessageID, "email-"))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "sending email")
	assert.Contains(t, logOutput, "test@example.com")
	assert.Contains(t, logOutput, "Test Subject")
	assert.Contains(t, logOutput, result.MessageID)
}

func TestConsoleSender_SendSMS(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)
	sender := NewConsoleSender(logger)

	cmd := command.SendSMSCommand{
		To:      "+1234567890",
		Message: "Hello SMS",
	}

	result, err := sender.SendSMS(context.Background(), cmd)
	require.NoError(t, err)
	assert.NotEmpty(t, result.MessageID)
	assert.Equal(t, "sent", result.Status)
	assert.True(t, strings.HasPrefix(result.MessageID, "sms-"))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "sending sms")
	assert.Contains(t, logOutput, "+1234567890")
}

func TestConsoleSender_EmailWithTemplate(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)
	sender := NewConsoleSender(logger)

	cmd := command.SendEmailCommand{
		To:         "user@example.com",
		Subject:    "Welcome",
		Body:       "Welcome to our service!",
		TemplateID: "welcome-template",
		TemplateData: map[string]string{
			"name": "John",
		},
	}

	result, err := sender.SendEmail(context.Background(), cmd)
	require.NoError(t, err)
	assert.NotEmpty(t, result.MessageID)
	assert.Equal(t, "sent", result.Status)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "welcome-template")
}
