package sender

import (
	"context"
	"fmt"

	"github.com/beabys/wms/notification-service/internal/application/notification/command"
)

// SMTPSender is a stub implementation that will be replaced in a future phase.
type SMTPSender struct {
	host     string
	port     int
	username string
	password string
}

// NewSMTPSender creates a new SMTPSender stub.
func NewSMTPSender(host string, port int, username, password string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

// SendEmail is not implemented.
func (s *SMTPSender) SendEmail(_ context.Context, _ command.SendEmailCommand) (command.SendEmailResult, error) {
	return command.SendEmailResult{}, fmt.Errorf("smtp sender not implemented")
}

// SendSMS is not implemented.
func (s *SMTPSender) SendSMS(_ context.Context, _ command.SendSMSCommand) (command.SendSMSResult, error) {
	return command.SendSMSResult{}, fmt.Errorf("smtp sender not implemented")
}
