package notificationgrpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notificationv1 "github.com/beabys/wms/proto/gen/go/notification/v1"

	"github.com/beabys/wms/notification-service/internal/application/notification/command"
	"github.com/beabys/wms/notification-service/internal/application/notification/usecase"
)

// NotificationServer implements the notification.v1.NotificationService gRPC server.
type NotificationServer struct {
	notificationv1.UnimplementedNotificationServiceServer
	svc usecase.NotificationServiceHandler
}

// NewNotificationServer creates a new NotificationServer.
func NewNotificationServer(svc usecase.NotificationServiceHandler) *NotificationServer {
	return &NotificationServer{svc: svc}
}

// SendEmail handles SendEmail RPC.
func (s *NotificationServer) SendEmail(ctx context.Context, req *notificationv1.SendEmailRequest) (*notificationv1.SendEmailResponse, error) {
	cmd := command.SendEmailCommand{
		To:           req.To,
		Subject:      req.Subject,
		Body:         req.Body,
		TemplateID:   req.TemplateId,
		TemplateData: req.TemplateData,
	}

	result, err := s.svc.SendEmail(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "send email: %v", err)
	}

	return &notificationv1.SendEmailResponse{
		MessageId: result.MessageID,
		Status:    result.Status,
	}, nil
}

// SendSMS handles SendSMS RPC.
func (s *NotificationServer) SendSMS(ctx context.Context, req *notificationv1.SendSMSRequest) (*notificationv1.SendSMSResponse, error) {
	cmd := command.SendSMSCommand{
		To:      req.To,
		Message: req.Message,
	}

	result, err := s.svc.SendSMS(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "send sms: %v", err)
	}

	return &notificationv1.SendSMSResponse{
		MessageId: result.MessageID,
		Status:    result.Status,
	}, nil
}
