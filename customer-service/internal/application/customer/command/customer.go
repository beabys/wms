package command

// RegisterCustomerCommand represents data needed to register a customer.
type RegisterCustomerCommand struct {
	Token       string
	CompanyName string
	Email       string
	Password    string
	Phone       string
	VatNumber   string
	Address     string
	City        string
	PostalCode  string
	Country     string
}

// RegisterCustomerResult contains the result of a registration.
type RegisterCustomerResult struct {
	CustomerID  string
	AccessToken string
}

// GetCustomerQuery retrieves a customer by ID.
type GetCustomerQuery struct {
	CustomerID string
}

// UpdateCustomerCommand updates a customer's fields.
type UpdateCustomerCommand struct {
	CustomerID string
	CompanyName string
	Phone       string
	VatNumber   string
	Address     string
	City        string
	PostalCode  string
	Country     string
}

// ListCustomersQuery represents paginated list filters.
type ListCustomersQuery struct {
	Page     int
	PageSize int
	Status   string
}

// ListCustomersResult contains paginated customer results.
type ListCustomersResult struct {
	Customers  []*CustomerResult
	TotalCount int
	Page       int
	PageSize   int
}

// ApproveCustomerCommand approves a customer.
type ApproveCustomerCommand struct {
	CustomerID string
	ApprovedBy string
}

// RejectCustomerCommand rejects a customer with a reason.
type RejectCustomerCommand struct {
	CustomerID string
	Reason     string
	RejectedBy string
}

// SuspendCustomerCommand suspends a customer.
type SuspendCustomerCommand struct {
	CustomerID string
	Reason     string
	SuspendedBy string
}

// RestoreCustomerCommand restores a suspended customer.
type RestoreCustomerCommand struct {
	CustomerID string
	RestoredBy string
}

// CustomerResult contains customer data.
type CustomerResult struct {
	ID             string
	CompanyName    string
	Email          string
	Phone          string
	VatNumber      string
	Address        string
	City           string
	PostalCode     string
	Country        string
	Status         string
	CompanyAdminID string
	RejectReason   string
	CreatedAt      string
	UpdatedAt      string
}
