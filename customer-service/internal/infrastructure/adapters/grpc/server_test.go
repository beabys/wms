package grpcdapter

import (
	"context"
	"fmt"
	"testing"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/transformer"
	"github.com/beabys/wms/customer-service/internal/application/customer/usecase"
	grpcmocks "github.com/beabys/wms/customer-service/mocks/infrastructure/adapters/grpc"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCustomerServer_GetCustomer(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("GetCustomer", mock.Anything, "cust-1").Return(&command.CustomerResult{
		ID:          "cust-1",
		CompanyName: "Test Corp",
		Email:       "admin@testcorp.com",
		Status:      "active",
	}, nil)

	s := NewCustomerServer(mockUC)
	resp, err := s.GetCustomer(context.Background(), &customerv1.GetCustomerRequest{Id: "cust-1"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "cust-1", resp.GetCustomer().GetId())
	assert.Equal(t, "Test Corp", resp.GetCustomer().GetCompanyName())
	assert.Equal(t, "admin@testcorp.com", resp.GetCustomer().GetEmail())
	assert.Equal(t, "active", resp.GetCustomer().GetStatus())
}

func TestCustomerServer_GetCustomerByAdminID_Success(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("GetCustomerByAdminID", mock.Anything, "admin-user-1").Return(&command.CustomerResult{
		ID:             "cust-1",
		CompanyName:    "Test Corp",
		Email:          "admin@testcorp.com",
		Status:         "active",
		CompanyAdminID: "admin-user-1",
	}, nil)

	s := NewCustomerServer(mockUC)
	ctx := context.WithValue(context.Background(), contextKey("auth_claims"), &Claims{UserID: "admin-user-1"})
	resp, err := s.GetCustomerByAdminID(ctx, &customerv1.GetCustomerByAdminIDRequest{AdminUserId: "admin-user-1"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetCustomer())
	assert.Equal(t, "cust-1", resp.GetCustomer().GetId())
	assert.Equal(t, "Test Corp", resp.GetCustomer().GetCompanyName())
	assert.Equal(t, "admin@testcorp.com", resp.GetCustomer().GetEmail())
	assert.Equal(t, "active", resp.GetCustomer().GetStatus())
	assert.Equal(t, "admin-user-1", resp.GetCustomer().GetCompanyAdminId())
}

func TestCustomerServer_GetCustomerByAdminID_NotFound(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("GetCustomerByAdminID", mock.Anything, "nonexistent").Return(nil, fmt.Errorf("customer not found"))

	s := NewCustomerServer(mockUC)
	ctx := context.WithValue(context.Background(), contextKey("auth_claims"), &Claims{UserID: "nonexistent"})
	resp, err := s.GetCustomerByAdminID(ctx, &customerv1.GetCustomerByAdminIDRequest{AdminUserId: "nonexistent"})
	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestNewCustomerServer(t *testing.T) {
	uc := &usecase.CustomerUseCase{}
	s := NewCustomerServer(uc)
	assert.NotNil(t, s)
}

func TestCustomerProtoConversion(t *testing.T) {
	r := &command.CustomerResult{
		ID:          "cust-1",
		CompanyName: "Test Corp",
		Email:       "admin@testcorp.com",
		Phone:       "1234567890",
		Status:      "active",
		CreatedAt:   "2024-06-10T12:00:00Z",
	}
	p := transformer.CustomerResultToProto(r)
	assert.NotNil(t, p)
	assert.Equal(t, "cust-1", p.GetId())
	assert.Equal(t, "Test Corp", p.GetCompanyName())
	assert.Equal(t, "admin@testcorp.com", p.GetEmail())
	assert.Equal(t, "1234567890", p.GetPhone())
	assert.Equal(t, "active", p.GetStatus())
	require.NotNil(t, p.GetCreatedAt())
	assert.Equal(t, int64(1718020800), p.GetCreatedAt().GetSeconds())
	assert.Equal(t, int32(0), p.GetCreatedAt().GetNanos())
}

func TestCustomerProtoConversion_NoCreatedAt(t *testing.T) {
	r := &command.CustomerResult{
		ID:          "cust-2",
		CompanyName: "No Date Corp",
		Email:       "nodate@test.com",
		Status:      "pending",
	}
	p := transformer.CustomerResultToProto(r)
	assert.NotNil(t, p)
	assert.Nil(t, p.GetCreatedAt())
}

func TestCustomerProtoConversion_Nil(t *testing.T) {
	assert.Nil(t, transformer.CustomerResultToProto(nil))
}

func TestMapError_NotFound(t *testing.T) {
	err := mapError(assert.AnError)
	assert.Error(t, err)
}

func TestRegisterCustomer_Success(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("RegisterCustomer", mock.Anything, mock.AnythingOfType("command.RegisterCustomerCommand")).Return(&command.RegisterCustomerResult{
		CustomerID:  "cust-123",
		AccessToken: "jwt-token-xxx",
	}, nil)

	s := NewCustomerServer(mockUC)
	resp, err := s.RegisterCustomer(context.Background(), &customerv1.RegisterCustomerRequest{
		Token:       "invite-token-abc",
		CompanyName: "New Corp",
		Email:       "admin@newcorp.com",
		Password:    "secret123",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "cust-123", resp.GetCustomer().GetId())
	assert.Equal(t, "jwt-token-xxx", resp.GetAccessToken())
}

func TestRegisterCustomer_InvalidToken(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("RegisterCustomer", mock.Anything, mock.AnythingOfType("command.RegisterCustomerCommand")).Return(nil, status.Error(codes.InvalidArgument, "invite token has expired"))

	s := NewCustomerServer(mockUC)
	resp, err := s.RegisterCustomer(context.Background(), &customerv1.RegisterCustomerRequest{
		Token:       "expired-token",
		CompanyName: "New Corp",
		Email:       "admin@newcorp.com",
		Password:    "secret123",
	})
	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestApproveCustomer_Success(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("ApproveCustomer", mock.Anything, mock.MatchedBy(func(cmd command.ApproveCustomerCommand) bool {
		return cmd.ApprovedBy == "admin-123"
	})).Return(&command.CustomerResult{
		ID:     "cust-1",
		Status: "active",
	}, nil)

	s := NewCustomerServer(mockUC)
	ctx := context.WithValue(context.Background(), contextKey("auth_claims"), &Claims{UserID: "admin-123"})
	resp, err := s.ApproveCustomer(ctx, &customerv1.ApproveCustomerRequest{
		Id: "cust-1",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetCustomer())
	assert.Equal(t, "cust-1", resp.GetCustomer().GetId())
	assert.Equal(t, "active", resp.GetCustomer().GetStatus())
}

func TestRejectCustomer_Success(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("RejectCustomer", mock.Anything, mock.MatchedBy(func(cmd command.RejectCustomerCommand) bool {
		return cmd.RejectedBy == "admin-123" && cmd.Reason == "invalid docs"
	})).Return(&command.CustomerResult{
		ID:     "cust-1",
		Status: "rejected",
	}, nil)

	s := NewCustomerServer(mockUC)
	ctx := context.WithValue(context.Background(), contextKey("auth_claims"), &Claims{UserID: "admin-123"})
	resp, err := s.RejectCustomer(ctx, &customerv1.RejectCustomerRequest{
		Id:     "cust-1",
		Reason: "invalid docs",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.GetSuccess())
}

func TestSuspendCustomer_Success(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("SuspendCustomer", mock.Anything, mock.MatchedBy(func(cmd command.SuspendCustomerCommand) bool {
		return cmd.SuspendedBy == "admin-123" && cmd.Reason == "policy violation"
	})).Return(&command.CustomerResult{
		ID:     "cust-1",
		Status: "suspended",
	}, nil)

	s := NewCustomerServer(mockUC)
	ctx := context.WithValue(context.Background(), contextKey("auth_claims"), &Claims{UserID: "admin-123"})
	resp, err := s.SuspendCustomer(ctx, &customerv1.SuspendCustomerRequest{
		Id:     "cust-1",
		Reason: "policy violation",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.GetSuccess())
}

func TestRestoreCustomer_Success(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("RestoreCustomer", mock.Anything, mock.MatchedBy(func(cmd command.RestoreCustomerCommand) bool {
		return cmd.RestoredBy == "admin-123"
	})).Return(&command.CustomerResult{
		ID:     "cust-1",
		Status: "active",
	}, nil)

	s := NewCustomerServer(mockUC)
	ctx := context.WithValue(context.Background(), contextKey("auth_claims"), &Claims{UserID: "admin-123"})
	resp, err := s.RestoreCustomer(ctx, &customerv1.RestoreCustomerRequest{
		Id: "cust-1",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetCustomer())
	assert.Equal(t, "cust-1", resp.GetCustomer().GetId())
	assert.Equal(t, "active", resp.GetCustomer().GetStatus())
}

func TestRestoreCustomer_NotFound(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("RestoreCustomer", mock.Anything, mock.AnythingOfType("command.RestoreCustomerCommand")).Return(nil, fmt.Errorf("customer not found"))

	s := NewCustomerServer(mockUC)
	ctx := context.WithValue(context.Background(), contextKey("auth_claims"), &Claims{UserID: "admin-123"})
	resp, err := s.RestoreCustomer(ctx, &customerv1.RestoreCustomerRequest{
		Id: "nonexistent",
	})
	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestRestoreCustomer_NotSuspended(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("RestoreCustomer", mock.Anything, mock.AnythingOfType("command.RestoreCustomerCommand")).Return(nil, fmt.Errorf("customer is not in suspended status"))

	s := NewCustomerServer(mockUC)
	ctx := context.WithValue(context.Background(), contextKey("auth_claims"), &Claims{UserID: "admin-123"})
	resp, err := s.RestoreCustomer(ctx, &customerv1.RestoreCustomerRequest{
		Id: "cust-1",
	})
	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
}

func TestAuthInterceptor_RejectsUnauthenticated(t *testing.T) {
	interceptor := NewAuthInterceptor(nil, []string{}, &testLogger{})
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	// Call with empty context (no metadata)
	_, err := interceptor.Unary()(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/customer.v1.CustomerService/GetCustomer",
	}, handler)
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestListAuditLogs_Success(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("GetAuditLogs", mock.Anything, mock.AnythingOfType("command.ListAuditLogsQuery")).Return(&command.ListAuditLogsResult{
		Entries: []*command.AuditLogResult{
			{
				ID:          "audit-1",
				CustomerID:  "cust-1",
				Action:      "approved",
				PerformedBy: "admin-1",
				Details:     "",
				CreatedAt:   "2024-01-15T10:00:00Z",
			},
		},
		TotalCount: 1,
		Page:       1,
		PageSize:   20,
	}, nil)

	s := NewCustomerServer(mockUC)
	resp, err := s.ListAuditLogs(context.Background(), &customerv1.ListAuditLogsRequest{
		CustomerId: "cust-1",
		Page:       1,
		PageSize:   20,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetEntries(), 1)
	assert.Equal(t, "audit-1", resp.GetEntries()[0].GetId())
	assert.Equal(t, "cust-1", resp.GetEntries()[0].GetCustomerId())
	assert.Equal(t, "approved", resp.GetEntries()[0].GetAction())
	assert.Equal(t, "admin-1", resp.GetEntries()[0].GetPerformedBy())
	assert.NotZero(t, resp.GetEntries()[0].GetCreatedAt())
}

func TestListAuditLogs_NotFound(t *testing.T) {
	mockUC := grpcmocks.NewCustomerUseCase(t)
	mockUC.On("GetAuditLogs", mock.Anything, mock.AnythingOfType("command.ListAuditLogsQuery")).Return(nil, fmt.Errorf("customer not found"))

	s := NewCustomerServer(mockUC)
	resp, err := s.ListAuditLogs(context.Background(), &customerv1.ListAuditLogsRequest{
		CustomerId: "nonexistent",
	})
	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

// testLogger implements logger.Logger for tests.
type testLogger struct{}

func (l *testLogger) GetLogger() any                            { return nil }
func (l *testLogger) Debug(msg string, fields ...logger.LogField) {}
func (l *testLogger) Info(msg string, fields ...logger.LogField)  {}
func (l *testLogger) Warn(msg string, fields ...logger.LogField)  {}
func (l *testLogger) Error(msg string, err error, fields ...logger.LogField) {}
func (l *testLogger) Fatal(msg string, fields ...logger.LogField) {}
