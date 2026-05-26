package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"

	"github.com/beabys/wms/customer-service-bff/internal/application"
)

// withToken creates a context with JWT in gRPC metadata.
func withToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

// Client implements application.CustomerService by calling the customer gRPC service.
type Client struct {
	conn            *grpc.ClientConn
	customerClient  customerv1.CustomerServiceClient
}

// NewClient creates a new gRPC client.
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		conn:           conn,
		customerClient: customerv1.NewCustomerServiceClient(conn),
	}
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateCustomer sends a create customer request to the gRPC service.
func (c *Client) CreateCustomer(ctx context.Context, req application.CreateCustomerRequest, token string) (*application.Customer, error) {
	ctx = withToken(ctx, token)
	protoReq := &customerv1.CreateCustomerRequest{
		CompanyName: req.CompanyName,
		VatNumber:   req.VATNumber,
		Address: &commonv1.Address{
			Line1:      req.Line1,
			Line2:      req.Line2,
			City:       req.City,
			PostalCode: req.PostalCode,
			Country:    req.Country,
		},
		RateCardId: req.RateCardID,
	}

	resp, err := c.customerClient.CreateCustomer(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("grpc create customer: %w", err)
	}

	return protoCustomerToDTO(resp.GetCustomer()), nil
}

// GetCustomer retrieves a customer by ID.
func (c *Client) GetCustomer(ctx context.Context, id string, token string) (*application.Customer, error) {
	ctx = withToken(ctx, token)
	resp, err := c.customerClient.GetCustomer(ctx, &customerv1.GetCustomerRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("grpc get customer: %w", err)
	}
	return protoCustomerToDTO(resp.GetCustomer()), nil
}

// ListCustomers lists customers with optional status filter.
func (c *Client) ListCustomers(ctx context.Context, status string, page, pageSize int32, token string) (*application.ListResult, error) {
	ctx = withToken(ctx, token)
	req := &customerv1.ListCustomersRequest{}
	if page > 0 || pageSize > 0 {
		req.Pagination = &commonv1.Pagination{Page: page, Limit: pageSize}
	}

	resp, err := c.customerClient.ListCustomers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc list customers: %w", err)
	}

	customers := make([]application.Customer, len(resp.GetCustomers()))
	for i, p := range resp.GetCustomers() {
		customers[i] = *protoCustomerToDTO(p)
	}

	total := int32(0)
	pageOut := page
	pageSizeOut := pageSize
	if resp.GetPagination() != nil {
		total = resp.GetPagination().GetTotal()
		pageOut = resp.GetPagination().GetPage()
		pageSizeOut = resp.GetPagination().GetLimit()
	}

	return &application.ListResult{
		Customers: customers,
		Total:     total,
		Page:      pageOut,
		PageSize:  pageSizeOut,
	}, nil
}

// ApproveCustomer approves a pending customer.
func (c *Client) ApproveCustomer(ctx context.Context, id string, token string) (*application.Customer, error) {
	ctx = withToken(ctx, token)
	resp, err := c.customerClient.ApproveCustomer(ctx, &customerv1.ApproveCustomerRequest{CustomerId: id})
	if err != nil {
		return nil, fmt.Errorf("grpc approve customer: %w", err)
	}
	return protoCustomerToDTO(resp.GetCustomer()), nil
}

// SuspendCustomer suspends an active customer.
func (c *Client) SuspendCustomer(ctx context.Context, id, reason string, token string) (*application.Customer, error) {
	ctx = withToken(ctx, token)
	resp, err := c.customerClient.SuspendCustomer(ctx, &customerv1.SuspendCustomerRequest{
		CustomerId: id,
		Reason:     reason,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc suspend customer: %w", err)
	}
	return protoCustomerToDTO(resp.GetCustomer()), nil
}

// InviteCustomer sends an invitation email.
func (c *Client) InviteCustomer(ctx context.Context, email string, token string) (string, error) {
	ctx = withToken(ctx, token)
	resp, err := c.customerClient.InviteCustomer(ctx, &customerv1.InviteCustomerRequest{Email: email})
	if err != nil {
		return "", fmt.Errorf("grpc invite customer: %w", err)
	}
	return resp.GetInvitationLink(), nil
}

// ---------------------------------------------------------------------------
// Mapping
// ---------------------------------------------------------------------------

func protoCustomerToDTO(p *customerv1.Customer) *application.Customer {
	if p == nil {
		return nil
	}
	return &application.Customer{
		ID:            p.GetId(),
		CompanyName:   p.GetCompanyName(),
		VATNumber:     p.GetVatNumber(),
		Address:       fmt.Sprintf("%s, %s, %s %s, %s", p.GetAddress().GetLine1(), p.GetAddress().GetLine2(), p.GetAddress().GetCity(), p.GetAddress().GetPostalCode(), p.GetAddress().GetCountry()),
		Status:        p.GetStatus(),
		RateCardID:    p.GetRateCardId(),
		CreditBalance: 0, // would parse from proto
		CreatedAt:     p.GetCreatedAt(),
	}
}
