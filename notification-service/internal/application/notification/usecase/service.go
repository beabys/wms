package usecase

import (
	"context"
	"fmt"

	"github.com/beabys/wms/notification-service/internal/application/notification/command"
)

// SendEmail validates and sends an email via the configured sender.
func (s *NotificationService) SendEmail(ctx context.Context, cmd command.SendEmailCommand) (command.SendEmailResult, error) {
	if cmd.To == "" {
		return command.SendEmailResult{}, fmt.Errorf("email recipient is required")
	}
	if cmd.Subject == "" {
		return command.SendEmailResult{}, fmt.Errorf("email subject is required")
	}
	return s.sender.SendEmail(ctx, cmd)
}

// SendSMS validates and sends an SMS via the configured sender.
func (s *NotificationService) SendSMS(ctx context.Context, cmd command.SendSMSCommand) (command.SendSMSResult, error) {
	if cmd.To == "" {
		return command.SendSMSResult{}, fmt.Errorf("sms recipient is required")
	}
	if cmd.Message == "" {
		return command.SendSMSResult{}, fmt.Errorf("sms message is required")
	}
	return s.sender.SendSMS(ctx, cmd)
}
