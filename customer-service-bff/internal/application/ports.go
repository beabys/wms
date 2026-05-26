package application

import "context"

// Customer represents a customer DTO for the BFF layer.
type Customer struct {
	ID            string `json:"id"`
	CompanyName   string `json:"company_name"`
	VATNumber     string `json:"vat_number"`
	Address       string `json:"address"` // JSON string for simplicity
	Status        string `json:"status"`
	RateCardID    string `json:"rate_card_id"`
	CreditBalance int64  `json:"credit_balance"`
	CreatedAt     string `json:"created_at"`
}

// CreateCustomerRequest is the BFF-level request DTO.
type CreateCustomerRequest struct {
	CompanyName string `json:"company_name"`
	VATNumber   string `json:"vat_number"`
	Line1       string `json:"line1"`
	Line2       string `json:"line2"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	RateCardID  string `json:"rate_card_id"`
}

// ListResult holds paginated customer results.
type ListResult struct {
	Customers []Customer `json:"customers"`
	Total     int32      `json:"total"`
	Page      int32      `json:"page"`
	PageSize  int32      `json:"page_size"`
}

// CustomerService is the input port for the BFF — defined by consumer (application layer).
type CustomerService interface {
	CreateCustomer(ctx context.Context, req CreateCustomerRequest, token string) (*Customer, error)
	GetCustomer(ctx context.Context, id string, token string) (*Customer, error)
	ListCustomers(ctx context.Context, status string, page, pageSize int32, token string) (*ListResult, error)
	ApproveCustomer(ctx context.Context, id string, token string) (*Customer, error)
	SuspendCustomer(ctx context.Context, id, reason string, token string) (*Customer, error)
	InviteCustomer(ctx context.Context, email string, token string) (string, error)
}
