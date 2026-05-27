package sender

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/beabys/wms/notification-service/internal/application/notification/command"
)

// ConsoleSender logs messages to the console via zap.
type ConsoleSender struct {
	logger *zap.Logger
}

// NewConsoleSender creates a new ConsoleSender.
func NewConsoleSender(logger *zap.Logger) *ConsoleSender {
	return &ConsoleSender{logger: logger}
}

// SendEmail logs an email message and returns a success result.
func (s *ConsoleSender) SendEmail(_ context.Context, cmd command.SendEmailCommand) (command.SendEmailResult, error) {
	messageID := fmt.Sprintf("email-%x", time.Now().UnixNano())

	s.logger.Info("sending email",
		zap.String("message_id", messageID),
		zap.String("to", cmd.To),
		zap.String("subject", cmd.Subject),
		zap.Int("body_length", len(cmd.Body)),
		zap.String("template_id", cmd.TemplateID),
	)

	return command.SendEmailResult{
		MessageID: messageID,
		Status:    "sent",
	}, nil
}

// SendSMS logs an SMS message and returns a success result.
func (s *ConsoleSender) SendSMS(_ context.Context, cmd command.SendSMSCommand) (command.SendSMSResult, error) {
	messageID := fmt.Sprintf("sms-%x", time.Now().UnixNano())

	s.logger.Info("sending sms",
		zap.String("message_id", messageID),
		zap.String("to", cmd.To),
		zap.Int("message_length", len(cmd.Message)),
	)

	return command.SendSMSResult{
		MessageID: messageID,
		Status:    "sent",
	}, nil
}
