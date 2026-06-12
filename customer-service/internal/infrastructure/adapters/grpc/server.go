package grpcdapter

import (
	"context"
	"strings"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/transformer"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CustomerUseCase defines the interface for customer use cases.
type CustomerUseCase interface {
	RegisterCustomer(ctx context.Context, cmd command.RegisterCustomerCommand) (*command.RegisterCustomerResult, error)
	GetCustomer(ctx context.Context, id string) (*command.CustomerResult, error)
	GetCustomerByAdminID(ctx context.Context, adminUserID string) (*command.CustomerResult, error)
	UpdateCustomer(ctx context.Context, cmd command.UpdateCustomerCommand) (*command.CustomerResult, error)
	ListCustomers(ctx context.Context, query command.ListCustomersQuery) (*command.ListCustomersResult, error)
	ApproveCustomer(ctx context.Context, cmd command.ApproveCustomerCommand) (*command.CustomerResult, error)
	RejectCustomer(ctx context.Context, cmd command.RejectCustomerCommand) (*command.CustomerResult, error)
	SuspendCustomer(ctx context.Context, cmd command.SuspendCustomerCommand) (*command.CustomerResult, error)
	RestoreCustomer(ctx context.Context, cmd command.RestoreCustomerCommand) (*command.CustomerResult, error)
	GetAuditLogs(ctx context.Context, query command.ListAuditLogsQuery) (*command.ListAuditLogsResult, error)
}

// CustomerServer implements customerv1.CustomerServiceServer.
type CustomerServer struct {
	customerv1.UnimplementedCustomerServiceServer
	customerUC CustomerUseCase
}

// NewCustomerServer creates a new CustomerServer.
func NewCustomerServer(customerUC CustomerUseCase) *CustomerServer {
	return &CustomerServer{
		customerUC: customerUC,
	}
}

// RegisterCustomer registers a new customer using an invite token.
func (s *CustomerServer) RegisterCustomer(ctx context.Context, req *customerv1.RegisterCustomerRequest) (*customerv1.RegisterCustomerResponse, error) {
	result, err := s.customerUC.RegisterCustomer(ctx, command.RegisterCustomerCommand{
		Token:       req.GetToken(),
		CompanyName: req.GetCompanyName(),
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		Phone:       req.GetPhone(),
		VatNumber:   req.GetVatNumber(),
		Address:     req.GetAddress(),
		City:        req.GetCity(),
		PostalCode:  req.GetPostalCode(),
		Country:     req.GetCountry(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.RegisterCustomerResponse{
		Customer:    &customerv1.Customer{Id: result.CustomerID},
		AccessToken: result.AccessToken,
	}, nil
}

// GetCustomerByAdminID retrieves a customer by company_admin_id.
// Uses authenticated user ID from context claims (ignores request field for security).
func (s *CustomerServer) GetCustomerByAdminID(ctx context.Context, req *customerv1.GetCustomerByAdminIDRequest) (*customerv1.GetCustomerByAdminIDResponse, error) {
	claims := ClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	result, err := s.customerUC.GetCustomerByAdminID(ctx, claims.UserID)
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.GetCustomerByAdminIDResponse{
		Customer: transformer.CustomerResultToProto(result),
	}, nil
}

// GetCustomer retrieves a customer by ID.
func (s *CustomerServer) GetCustomer(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
	result, err := s.customerUC.GetCustomer(ctx, req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.GetCustomerResponse{
		Customer: transformer.CustomerResultToProto(result),
	}, nil
}

// UpdateCustomer updates an existing customer.
func (s *CustomerServer) UpdateCustomer(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
	result, err := s.customerUC.UpdateCustomer(ctx, command.UpdateCustomerCommand{
		CustomerID: req.GetId(),
		Phone:      req.GetPhone(),
		Address:    req.GetAddress(),
		City:       req.GetCity(),
		PostalCode: req.GetPostalCode(),
		Country:    req.GetCountry(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.UpdateCustomerResponse{
		Customer: transformer.CustomerResultToProto(result),
	}, nil
}

// ListCustomers retrieves a paginated list of customers.
func (s *CustomerServer) ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
	result, err := s.customerUC.ListCustomers(ctx, command.ListCustomersQuery{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Status:   req.GetStatus(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	customers := make([]*customerv1.Customer, 0, len(result.Customers))
	for _, c := range result.Customers {
		customers = append(customers, transformer.CustomerResultToProto(c))
	}

	totalPages := 0
	if result.PageSize > 0 {
		totalPages = (result.TotalCount + result.PageSize - 1) / result.PageSize
	}

	return &customerv1.ListCustomersResponse{
		Customers: customers,
		Pagination: &commonv1.Pagination{
			Page:       int32(result.Page),
			PageSize:   int32(result.PageSize),
			Total:      int32(result.TotalCount),
			TotalPages: int32(totalPages),
		},
	}, nil
}

// ApproveCustomer approves a customer.
// User identity is extracted from JWT claims in context (server-authoritative).
func (s *CustomerServer) ApproveCustomer(ctx context.Context, req *customerv1.ApproveCustomerRequest) (*customerv1.ApproveCustomerResponse, error) {
	claims := ClaimsFromContext(ctx)
	performedBy := ""
	if claims != nil {
		performedBy = claims.UserID
	}
	result, err := s.customerUC.ApproveCustomer(ctx, command.ApproveCustomerCommand{
		CustomerID: req.GetId(),
		ApprovedBy: performedBy,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.ApproveCustomerResponse{
		Customer: transformer.CustomerResultToProto(result),
	}, nil
}

// RejectCustomer rejects a customer with a reason.
// User identity is extracted from JWT claims in context (server-authoritative).
func (s *CustomerServer) RejectCustomer(ctx context.Context, req *customerv1.RejectCustomerRequest) (*customerv1.RejectCustomerResponse, error) {
	claims := ClaimsFromContext(ctx)
	performedBy := ""
	if claims != nil {
		performedBy = claims.UserID
	}
	_, err := s.customerUC.RejectCustomer(ctx, command.RejectCustomerCommand{
		CustomerID: req.GetId(),
		Reason:     req.GetReason(),
		RejectedBy: performedBy,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.RejectCustomerResponse{
		Success: true,
	}, nil
}

// SuspendCustomer suspends a customer.
// User identity is extracted from JWT claims in context (server-authoritative).
func (s *CustomerServer) SuspendCustomer(ctx context.Context, req *customerv1.SuspendCustomerRequest) (*customerv1.SuspendCustomerResponse, error) {
	claims := ClaimsFromContext(ctx)
	performedBy := ""
	if claims != nil {
		performedBy = claims.UserID
	}
	_, err := s.customerUC.SuspendCustomer(ctx, command.SuspendCustomerCommand{
		CustomerID:  req.GetId(),
		Reason:      req.GetReason(),
		SuspendedBy: performedBy,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.SuspendCustomerResponse{
		Success: true,
	}, nil
}

// RestoreCustomer restores a suspended customer.
// User identity is extracted from JWT claims in context (server-authoritative).
func (s *CustomerServer) RestoreCustomer(ctx context.Context, req *customerv1.RestoreCustomerRequest) (*customerv1.RestoreCustomerResponse, error) {
	claims := ClaimsFromContext(ctx)
	performedBy := ""
	if claims != nil {
		performedBy = claims.UserID
	}
	result, err := s.customerUC.RestoreCustomer(ctx, command.RestoreCustomerCommand{
		CustomerID: req.GetId(),
		RestoredBy: performedBy,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.RestoreCustomerResponse{
		Customer: transformer.CustomerResultToProto(result),
	}, nil
}

// ListAuditLogs retrieves paginated audit logs for a customer.
func (s *CustomerServer) ListAuditLogs(ctx context.Context, req *customerv1.ListAuditLogsRequest) (*customerv1.ListAuditLogsResponse, error) {
	result, err := s.customerUC.GetAuditLogs(ctx, command.ListAuditLogsQuery{
		CustomerID: req.GetCustomerId(),
		Page:       int(req.GetPage()),
		PageSize:   int(req.GetPageSize()),
	})
	if err != nil {
		return nil, mapError(err)
	}

	entries := make([]*customerv1.AuditEntry, 0, len(result.Entries))
	for _, e := range result.Entries {
		entries = append(entries, transformer.AuditLogResultToProto(e))
	}

	totalPages := 0
	if result.PageSize > 0 {
		totalPages = (result.TotalCount + result.PageSize - 1) / result.PageSize
	}

	return &customerv1.ListAuditLogsResponse{
		Entries: entries,
		Pagination: &commonv1.Pagination{
			Page:       int32(result.Page),
			PageSize:   int32(result.PageSize),
			Total:      int32(result.TotalCount),
			TotalPages: int32(totalPages),
		},
	}, nil
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "not found"):
		return status.Error(codes.NotFound, errStr)
	case strings.Contains(errStr, "is required"):
		return status.Error(codes.InvalidArgument, errStr)
	case strings.Contains(errStr, "not in pending"), strings.Contains(errStr, "not in suspended"):
		return status.Error(codes.FailedPrecondition, errStr)
	case strings.Contains(errStr, "already used"):
		return status.Error(codes.InvalidArgument, errStr)
	case strings.Contains(errStr, "has expired"):
		return status.Error(codes.InvalidArgument, errStr)
	case strings.Contains(errStr, "already in use"):
		return status.Error(codes.AlreadyExists, errStr)
	default:
		return status.Error(codes.Internal, "internal error")
	}
}


