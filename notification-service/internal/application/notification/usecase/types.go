package usecase

import (
	"context"

	"github.com/beabys/wms/notification-service/internal/application/notification/command"
)

// Sender sends notification messages.
type Sender interface {
	SendEmail(ctx context.Context, cmd command.SendEmailCommand) (command.SendEmailResult, error)
	SendSMS(ctx context.Context, cmd command.SendSMSCommand) (command.SendSMSResult, error)
}

// NotificationServiceHandler handles notification operations.
type NotificationServiceHandler interface {
	SendEmail(ctx context.Context, cmd command.SendEmailCommand) (command.SendEmailResult, error)
	SendSMS(ctx context.Context, cmd command.SendSMSCommand) (command.SendSMSResult, error)
}

// NotificationService implements NotificationServiceHandler.
type NotificationService struct {
	sender Sender
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(sender Sender) *NotificationService {
	return &NotificationService{sender: sender}
}
