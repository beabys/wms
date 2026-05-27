package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/notification-service/internal/application/notification/command"
)

// mockSender implements Sender for testing.
type mockSender struct {
	sendEmailFn func(ctx context.Context, cmd command.SendEmailCommand) (command.SendEmailResult, error)
	sendSMSFn   func(ctx context.Context, cmd command.SendSMSCommand) (command.SendSMSResult, error)
}

func (m *mockSender) SendEmail(ctx context.Context, cmd command.SendEmailCommand) (command.SendEmailResult, error) {
	return m.sendEmailFn(ctx, cmd)
}

func (m *mockSender) SendSMS(ctx context.Context, cmd command.SendSMSCommand) (command.SendSMSResult, error) {
	return m.sendSMSFn(ctx, cmd)
}

func TestSendEmail_Success(t *testing.T) {
	mock := &mockSender{
		sendEmailFn: func(_ context.Context, cmd command.SendEmailCommand) (command.SendEmailResult, error) {
			return command.SendEmailResult{MessageID: "email-abc123", Status: "sent"}, nil
		},
	}
	svc := NewNotificationService(mock)

	result, err := svc.SendEmail(context.Background(), command.SendEmailCommand{
		To:      "test@example.com",
		Subject: "Hello",
		Body:    "World",
	})
	require.NoError(t, err)
	assert.Equal(t, "email-abc123", result.MessageID)
	assert.Equal(t, "sent", result.Status)
}

func TestSendEmail_EmptyTo(t *testing.T) {
	mock := &mockSender{}
	svc := NewNotificationService(mock)

	_, err := svc.SendEmail(context.Background(), command.SendEmailCommand{
		To:      "",
		Subject: "Hello",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "recipient is required")
}

func TestSendEmail_EmptySubject(t *testing.T) {
	mock := &mockSender{}
	svc := NewNotificationService(mock)

	_, err := svc.SendEmail(context.Background(), command.SendEmailCommand{
		To:      "test@example.com",
		Subject: "",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "subject is required")
}

func TestSendSMS_Success(t *testing.T) {
	mock := &mockSender{
		sendSMSFn: func(_ context.Context, cmd command.SendSMSCommand) (command.SendSMSResult, error) {
			return command.SendSMSResult{MessageID: "sms-abc123", Status: "sent"}, nil
		},
	}
	svc := NewNotificationService(mock)

	result, err := svc.SendSMS(context.Background(), command.SendSMSCommand{
		To:      "+1234567890",
		Message: "Hello SMS",
	})
	require.NoError(t, err)
	assert.Equal(t, "sms-abc123", result.MessageID)
	assert.Equal(t, "sent", result.Status)
}

func TestSendSMS_EmptyTo(t *testing.T) {
	mock := &mockSender{}
	svc := NewNotificationService(mock)

	_, err := svc.SendSMS(context.Background(), command.SendSMSCommand{
		To:      "",
		Message: "Hello",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "recipient is required")
}

func TestSendSMS_EmptyMessage(t *testing.T) {
	mock := &mockSender{}
	svc := NewNotificationService(mock)

	_, err := svc.SendSMS(context.Background(), command.SendSMSCommand{
		To:      "+1234567890",
		Message: "",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "message is required")
}
