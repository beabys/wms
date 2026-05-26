package grpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/usecase"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// Server implements customerv1.CustomerServiceServer.
type Server struct {
	customerv1.UnimplementedCustomerServiceServer

	log             *zap.Logger
	customerService *usecase.CustomerService
}

// NewServer creates a new gRPC customer server.
func NewServer(log *zap.Logger, customerService *usecase.CustomerService) *Server {
	return &Server{
		log:             log,
		customerService: customerService,
	}
}

// CreateCustomer implements customer.v1.CustomerService.
func (s *Server) CreateCustomer(ctx context.Context, req *customerv1.CreateCustomerRequest) (*customerv1.CreateCustomerResponse, error) {
	addr, err := protoToDomainAddress(req.GetAddress())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %v", err)
	}

	cmd := &command.CreateCustomerCommand{
		CompanyName: req.GetCompanyName(),
		VATNumber:   req.GetVatNumber(),
		Address:     addr,
		RateCardID:  req.GetRateCardId(),
	}

	customer, err := s.customerService.CreateCustomer(ctx, cmd)
	if err != nil {
		s.log.Error("create customer failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "create customer: %v", err)
	}

	return &customerv1.CreateCustomerResponse{
		Customer: domainToProtoCustomer(customer),
	}, nil
}

// GetCustomer implements customer.v1.CustomerService.
func (s *Server) GetCustomer(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
	q := &command.GetCustomerQuery{CustomerID: req.GetId()}

	customer, err := s.customerService.GetCustomer(ctx, q)
	if err != nil {
		s.log.Error("get customer failed", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "get customer: %v", err)
	}

	return &customerv1.GetCustomerResponse{
		Customer: domainToProtoCustomer(customer),
	}, nil
}

// ListCustomers implements customer.v1.CustomerService.
func (s *Server) ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
	q := &command.ListCustomersQuery{
		Page:     req.GetPagination().GetPage(),
		PageSize: req.GetPagination().GetLimit(),
	}

	customers, total, err := s.customerService.ListCustomers(ctx, q)
	if err != nil {
		s.log.Error("list customers failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "list customers: %v", err)
	}

	protoCustomers := make([]*customerv1.Customer, len(customers))
	for i, c := range customers {
		protoCustomers[i] = domainToProtoCustomer(c)
	}

	return &customerv1.ListCustomersResponse{
		Customers: protoCustomers,
		Pagination: &commonv1.Pagination{
			Page:  req.GetPagination().GetPage(),
			Limit: req.GetPagination().GetLimit(),
			Total: total,
		},
	}, nil
}

// ApproveCustomer implements customer.v1.CustomerService.
func (s *Server) ApproveCustomer(ctx context.Context, req *customerv1.ApproveCustomerRequest) (*customerv1.ApproveCustomerResponse, error) {
	cmd := &command.ApproveCustomerCommand{CustomerID: req.GetCustomerId()}

	customer, err := s.customerService.ApproveCustomer(ctx, cmd)
	if err != nil {
		s.log.Error("approve customer failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "approve customer: %v", err)
	}

	return &customerv1.ApproveCustomerResponse{
		Customer: domainToProtoCustomer(customer),
	}, nil
}

// SuspendCustomer implements customer.v1.CustomerService.
func (s *Server) SuspendCustomer(ctx context.Context, req *customerv1.SuspendCustomerRequest) (*customerv1.SuspendCustomerResponse, error) {
	cmd := &command.SuspendCustomerCommand{
		CustomerID: req.GetCustomerId(),
		Reason:     req.GetReason(),
	}

	customer, err := s.customerService.SuspendCustomer(ctx, cmd)
	if err != nil {
		s.log.Error("suspend customer failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "suspend customer: %v", err)
	}

	return &customerv1.SuspendCustomerResponse{
		Customer: domainToProtoCustomer(customer),
	}, nil
}

// InviteCustomer implements customer.v1.CustomerService.
func (s *Server) InviteCustomer(ctx context.Context, req *customerv1.InviteCustomerRequest) (*customerv1.InviteCustomerResponse, error) {
	cmd := &command.InviteCustomerCommand{Email: req.GetEmail()}

	invite, err := s.customerService.InviteCustomer(ctx, cmd)
	if err != nil {
		s.log.Error("invite customer failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invite customer: %v", err)
	}

	return &customerv1.InviteCustomerResponse{
		InvitationLink: invite.Token,
	}, nil
}

// ---------------------------------------------------------------------------
// Mapping helpers
// ---------------------------------------------------------------------------

func protoToDomainAddress(addr *commonv1.Address) (model.Address, error) {
	if addr == nil {
		return model.Address{}, fmt.Errorf("address is required")
	}
	return model.NewAddress(addr.GetLine1(), addr.GetLine2(), addr.GetCity(), addr.GetPostalCode(), addr.GetCountry())
}

func domainToProtoCustomer(c *model.Customer) *customerv1.Customer {
	return &customerv1.Customer{
		Id:          c.ID,
		CompanyName: c.CompanyName,
		VatNumber:   c.VATNumber,
		Address: &commonv1.Address{
			Line1:      c.Address.Line1,
			Line2:      c.Address.Line2,
			City:       c.Address.City,
			PostalCode: c.Address.PostalCode,
			Country:    c.Address.Country,
		},
		Status:        string(c.Status),
		RateCardId:    c.RateCardID,
		CreditBalance: fmt.Sprintf("%d", c.CreditBalance),
		CreatedAt:     c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
