package httpadapter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/beabys/wms/customer-service-bff/internal/api/v1"
	"github.com/beabys/wms/customer-service-bff/internal/domain/model"
	"github.com/beabys/wms/pkg/logger"
)

// CustomerHandlerService defines what the HTTP handler needs from the use case.
// Interface defined by consumer (handler layer).
type CustomerHandlerService interface {
	RegisterCustomer(ctx context.Context, req *model.RegisterCustomerRequest) (*model.RegisterCustomerResponse, error)
	GetCustomer(ctx context.Context, id string) (*model.CustomerResponse, error)
	GetMyCustomer(ctx context.Context, adminUserID string) (*model.CustomerResponse, error)
	ListCustomers(ctx context.Context, status *string, page, pageSize *int) (*model.CustomerListResponse, error)
	UpdateCustomer(ctx context.Context, id string, req *model.UpdateCustomerRequest) (*model.CustomerResponse, error)
	ApproveCustomer(ctx context.Context, id, approvedBy string) (*model.CustomerResponse, error)
	RejectCustomer(ctx context.Context, id, reason string) error
	SuspendCustomer(ctx context.Context, id, reason string) error
	RestoreCustomer(ctx context.Context, id string, restoredBy string) (*model.CustomerResponse, error)
	AssignCompanyRole(ctx context.Context, req *model.AssignCompanyRoleRequest) error
	ListCompanyRoles(ctx context.Context, customerID string) (*model.CompanyRoleListResponse, error)
	GetUserPermissions(ctx context.Context, userID, customerID string) (*model.PermissionsResponse, error)
	ListAuditLogs(ctx context.Context, customerID string, page, pageSize int) (*model.AuditLogListResponse, error)
}

// Server implements v1.ServerInterface for the customer BFF HTTP API.
type Server struct {
	customerSvc CustomerHandlerService
	logger      logger.Logger
}

// NewServer creates a new Server with the given customer service.
func NewServer(customerSvc CustomerHandlerService, log logger.Logger) *Server {
	return &Server{
		customerSvc: customerSvc,
		logger:      log,
	}
}

// compile-time check
var _ v1.ServerInterface = (*Server)(nil)

// writeSuccess writes a success JSON response with status 200.
func (s *Server) writeSuccess(w http.ResponseWriter, data interface{}) {
	successResponseJSON(w, http.StatusOK, data)
}

// writeCreated writes a success JSON response with status 201.
func (s *Server) writeCreated(w http.ResponseWriter, data interface{}) {
	successResponseJSON(w, http.StatusCreated, data)
}

// writeError writes an error JSON response with the given status and message.
func (s *Server) writeError(w http.ResponseWriter, status int, msg string) {
	s.logger.Warn("request failed",
		logger.LogField{Key: "status", Value: status},
		logger.LogField{Key: "error", Value: msg},
	)
	errorResponseJSON(w, status, fmt.Errorf("%s", msg))
}
