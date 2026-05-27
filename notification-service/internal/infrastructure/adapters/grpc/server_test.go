package notificationgrpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	notificationv1 "github.com/beabys/wms/proto/gen/go/notification/v1"

	notificationgrpc "github.com/beabys/wms/notification-service/internal/infrastructure/adapters/grpc"
	senderadapter "github.com/beabys/wms/notification-service/internal/infrastructure/adapters/sender"
	"github.com/beabys/wms/notification-service/internal/application/notification/usecase"
)

const bufSize = 1024 * 1024

func startTestServer(t *testing.T, ctx context.Context) notificationv1.NotificationServiceClient {
	t.Helper()

	logger := zaptest.NewLogger(t)
	sender := senderadapter.NewConsoleSender(logger)
	svc := usecase.NewNotificationService(sender)

	lis := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	notificationSrv := notificationgrpc.NewNotificationServer(svc)
	notificationv1.RegisterNotificationServiceServer(grpcServer, notificationSrv)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Errorf("server serve error: %v", err)
		}
	}()
	t.Cleanup(grpcServer.Stop)

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return notificationv1.NewNotificationServiceClient(conn)
}

func TestSendEmail_Success(t *testing.T) {
	ctx := context.Background()
	client := startTestServer(t, ctx)

	resp, err := client.SendEmail(ctx, &notificationv1.SendEmailRequest{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.MessageId)
	assert.Equal(t, "sent", resp.Status)
}

func TestSendEmail_MissingTo(t *testing.T) {
	ctx := context.Background()
	client := startTestServer(t, ctx)

	_, err := client.SendEmail(ctx, &notificationv1.SendEmailRequest{
		Subject: "Test Subject",
		Body:    "Test Body",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "recipient is required")
}

func TestSendEmail_MissingSubject(t *testing.T) {
	ctx := context.Background()
	client := startTestServer(t, ctx)

	_, err := client.SendEmail(ctx, &notificationv1.SendEmailRequest{
		To:   "test@example.com",
		Body: "Test Body",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "subject is required")
}

func TestSendEmail_WithTemplate(t *testing.T) {
	ctx := context.Background()
	client := startTestServer(t, ctx)

	resp, err := client.SendEmail(ctx, &notificationv1.SendEmailRequest{
		To:         "user@example.com",
		Subject:    "Welcome",
		Body:       "Welcome!",
		TemplateId: "welcome-template",
		TemplateData: map[string]string{
			"name": "John",
		},
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.MessageId)
	assert.Equal(t, "sent", resp.Status)
}

func TestSendSMS_Success(t *testing.T) {
	ctx := context.Background()
	client := startTestServer(t, ctx)

	resp, err := client.SendSMS(ctx, &notificationv1.SendSMSRequest{
		To:      "+1234567890",
		Message: "Hello SMS",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.MessageId)
	assert.Equal(t, "sent", resp.Status)
}

func TestSendSMS_MissingTo(t *testing.T) {
	ctx := context.Background()
	client := startTestServer(t, ctx)

	_, err := client.SendSMS(ctx, &notificationv1.SendSMSRequest{
		Message: "Hello",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "recipient is required")
}

func TestSendSMS_MissingMessage(t *testing.T) {
	ctx := context.Background()
	client := startTestServer(t, ctx)

	_, err := client.SendSMS(ctx, &notificationv1.SendSMSRequest{
		To: "+1234567890",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "message is required")
}
