package model

// CustomerResponse represents a customer in API responses.
type CustomerResponse struct {
	ID             string `json:"id"`
	CompanyName    string `json:"company_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	VatNumber      string `json:"vat_number"`
	Address        string `json:"address"`
	City           string `json:"city"`
	PostalCode     string `json:"postal_code"`
	Country        string `json:"country"`
	Status         string `json:"status"`
	CompanyAdminID string `json:"company_admin_id"`
	CreatedAt      string `json:"created_at"`
}

// RegisterCustomerRequest represents the request to register a new customer.
type RegisterCustomerRequest struct {
	Token       string `json:"token"`
	CompanyName string `json:"company_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	VatNumber   string `json:"vat_number"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
}

// RegisterCustomerResponse represents the response from registering a customer.
type RegisterCustomerResponse struct {
	Customer    *CustomerResponse `json:"customer"`
	AccessToken string            `json:"access_token"`
}

// CustomerListResponse represents a paginated list of customers.
type CustomerListResponse struct {
	Customers  []CustomerResponse `json:"customers"`
	Pagination Pagination         `json:"pagination"`
}

// Pagination contains pagination metadata.
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
}

// UpdateCustomerRequest represents the request body for updating a customer.
type UpdateCustomerRequest struct {
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// AssignCompanyRoleRequest represents the request to assign a company role.
type AssignCompanyRoleRequest struct {
	CustomerID  string   `json:"customer_id"`
	UserID      string   `json:"user_id"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions,omitempty"`
	AssignedBy  string   `json:"assigned_by,omitempty"`
}

// CompanyRoleResponse represents a company role in API responses.
type CompanyRoleResponse struct {
	ID          string   `json:"id"`
	CustomerID  string   `json:"customer_id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions,omitempty"`
	IsDefault   bool     `json:"is_default"`
}

// CompanyRoleListResponse represents a list of company roles.
type CompanyRoleListResponse struct {
	Roles []CompanyRoleResponse `json:"roles"`
}

// PermissionsResponse represents a list of permissions.
type PermissionsResponse struct {
	Permissions []string `json:"permissions"`
}

// AuditEntry represents a single audit log entry.
type AuditEntry struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customer_id"`
	Action      string `json:"action"`
	PerformedBy string `json:"performed_by"`
	Details     string `json:"details,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// AuditLogListResponse represents a paginated list of audit entries.
type AuditLogListResponse struct {
	Entries    []AuditEntry `json:"entries"`
	Pagination Pagination   `json:"pagination"`
}
