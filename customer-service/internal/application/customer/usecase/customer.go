package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/ports"
	"github.com/beabys/wms/customer-service/internal/application/customer/repository"
	"github.com/beabys/wms/customer-service/internal/application/customer/validator"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	"github.com/beabys/wms/pkg/logger"
	"github.com/google/uuid"
)

// CustomerUseCase handles customer operations.
type CustomerUseCase struct {
	logger        logger.Logger
	customerRepo  repository.CustomerRepository
	companyRepo   repository.CompanyRepository
	roleRepo      repository.CompanyRoleRepository
	auditLogRepo  repository.AuditLogRepository
	authClient    ports.AuthClient
}

// NewCustomerUseCase creates a new CustomerUseCase.
func NewCustomerUseCase(
	log logger.Logger,
	customerRepo repository.CustomerRepository,
	companyRepo repository.CompanyRepository,
	roleRepo repository.CompanyRoleRepository,
	auditLogRepo repository.AuditLogRepository,
	authClient ports.AuthClient,
) *CustomerUseCase {
	return &CustomerUseCase{
		logger:       log,
		customerRepo: customerRepo,
		companyRepo:  companyRepo,
		roleRepo:     roleRepo,
		auditLogRepo: auditLogRepo,
		authClient:   authClient,
	}
}

// RegisterCustomer registers a new customer using an invite token.
func (uc *CustomerUseCase) RegisterCustomer(ctx context.Context, cmd command.RegisterCustomerCommand) (*command.RegisterCustomerResult, error) {
	if err := validator.InviteToken(cmd.Token); err != nil {
		return nil, err
	}
	if err := validator.Email(cmd.Email); err != nil {
		return nil, err
	}
	if err := validator.NotEmpty(cmd.Password, "password"); err != nil {
		return nil, err
	}
	if err := validator.CompanyName(cmd.CompanyName); err != nil {
		return nil, err
	}

	// 1. Validate invite token via auth-service
	email, invitedBy, err := uc.authClient.ValidateInvite(ctx, cmd.Token)
	if err != nil {
		uc.logger.Error("invite validation failed", err)
		return nil, fmt.Errorf("invalid invite token")
	}
	uc.logger.Info("invite validated",
		logger.LogField{Key: "email", Value: email},
		logger.LogField{Key: "invited_by", Value: invitedBy},
	)

	// 2. Create user in auth-service (customer_admin role)
	authResp, err := uc.authClient.CreateUser(ctx, command.CreateUserRequest{
		Email:    cmd.Email,
		Password: cmd.Password,
		Name:     cmd.CompanyName + " Admin",
		Role:     "customer_admin",
	})
	if err != nil {
		uc.logger.Error("failed to create user in auth-service", err)
		return nil, fmt.Errorf("failed to create user account")
	}

	// 4. Create customer record
	now := time.Now()
	customer := &model.Customer{
		ID:             uuid.New().String(),
		CompanyName:    cmd.CompanyName,
		Email:          cmd.Email,
		Phone:          cmd.Phone,
		VatNumber:      cmd.VatNumber,
		Address:        cmd.Address,
		City:           cmd.City,
		PostalCode:     cmd.PostalCode,
		Country:        cmd.Country,
		Status:         model.CustomerStatusPending,
		CompanyAdminID: authResp.UserID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := uc.customerRepo.Create(ctx, customer); err != nil {
		uc.logger.Error("failed to create customer", err)
		return nil, fmt.Errorf("failed to create customer")
	}

	// 5. Create default company
	company := &model.Company{
		ID:         uuid.New().String(),
		CustomerID: customer.ID,
		Name:       cmd.CompanyName,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := uc.companyRepo.Create(ctx, company); err != nil {
		uc.logger.Error("failed to create company", err)
		return nil, fmt.Errorf("failed to create company")
	}

	// 6. Create default admin role with full permissions
	role := &model.CompanyRole{
		ID:          uuid.New().String(),
		CustomerID:  customer.ID,
		Name:        "admin",
		Permissions: []string{"*"},
		IsDefault:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.roleRepo.Create(ctx, role); err != nil {
		uc.logger.Error("failed to create default role", err)
		return nil, fmt.Errorf("failed to initialize company roles")
	}

	// 7. Assign admin role to the user
	if err := uc.roleRepo.AssignUserRole(ctx, authResp.UserID, role.ID, customer.ID, authResp.UserID); err != nil {
		uc.logger.Error("failed to assign admin role", err)
		return nil, fmt.Errorf("failed to assign default role")
	}

	uc.logger.Info("customer registered successfully",
		logger.LogField{Key: "customer_id", Value: customer.ID},
		logger.LogField{Key: "company", Value: cmd.CompanyName},
	)

	return &command.RegisterCustomerResult{
		CustomerID:  customer.ID,
		AccessToken: "", // Access token obtained via login flow
	}, nil
}

// GetCustomerByAdminID retrieves a customer by company_admin_id.
func (uc *CustomerUseCase) GetCustomerByAdminID(ctx context.Context, adminUserID string) (*command.CustomerResult, error) {
	if err := validator.NotEmpty(adminUserID, "admin user ID"); err != nil {
		return nil, err
	}

	customer, err := uc.customerRepo.FindByAdminID(ctx, adminUserID)
	if err != nil {
		uc.logger.Error("failed to get customer by admin id", err)
		return nil, fmt.Errorf("customer not found")
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	return customerToResult(customer), nil
}

// GetCustomer retrieves a customer by ID.
func (uc *CustomerUseCase) GetCustomer(ctx context.Context, id string) (*command.CustomerResult, error) {
	if err := validator.NotEmpty(id, "customer ID"); err != nil {
		return nil, err
	}

	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get customer", err)
		return nil, fmt.Errorf("customer not found")
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	return customerToResult(customer), nil
}

// UpdateCustomer updates an existing customer's fields.
func (uc *CustomerUseCase) UpdateCustomer(ctx context.Context, cmd command.UpdateCustomerCommand) (*command.CustomerResult, error) {
	if err := validator.CustomerID(cmd.CustomerID); err != nil {
		return nil, err
	}

	customer, err := uc.customerRepo.GetByID(ctx, cmd.CustomerID)
	if err != nil {
		uc.logger.Error("failed to get customer for update", err)
		return nil, fmt.Errorf("customer not found")
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	if cmd.CompanyName != "" {
		customer.CompanyName = cmd.CompanyName
	}
	if cmd.Phone != "" {
		customer.Phone = cmd.Phone
	}
	if cmd.VatNumber != "" {
		customer.VatNumber = cmd.VatNumber
	}
	if cmd.Address != "" {
		customer.Address = cmd.Address
	}
	if cmd.City != "" {
		customer.City = cmd.City
	}
	if cmd.PostalCode != "" {
		customer.PostalCode = cmd.PostalCode
	}
	if cmd.Country != "" {
		customer.Country = cmd.Country
	}
	customer.UpdatedAt = time.Now()

	if err := uc.customerRepo.Update(ctx, customer); err != nil {
		uc.logger.Error("failed to update customer", err)
		return nil, fmt.Errorf("failed to update customer")
	}

	return customerToResult(customer), nil
}

// ListCustomers retrieves a paginated list of customers.
func (uc *CustomerUseCase) ListCustomers(ctx context.Context, query command.ListCustomersQuery) (*command.ListCustomersResult, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	customers, total, err := uc.customerRepo.List(ctx, query.Page, query.PageSize, query.Status)
	if err != nil {
		uc.logger.Error("failed to list customers", err)
		return nil, fmt.Errorf("failed to list customers")
	}

	results := make([]*command.CustomerResult, 0, len(customers))
	for _, c := range customers {
		results = append(results, customerToResult(c))
	}

	return &command.ListCustomersResult{
		Customers:  results,
		TotalCount: total,
		Page:       query.Page,
		PageSize:   query.PageSize,
	}, nil
}

// ApproveCustomer approves a customer (sets status to active).
func (uc *CustomerUseCase) ApproveCustomer(ctx context.Context, cmd command.ApproveCustomerCommand) (*command.CustomerResult, error) {
	if err := validator.CustomerID(cmd.CustomerID); err != nil {
		return nil, err
	}

	customer, err := uc.customerRepo.GetByID(ctx, cmd.CustomerID)
	if err != nil {
		uc.logger.Error("failed to get customer for approval", err)
		return nil, fmt.Errorf("customer not found")
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	if customer.Status != model.CustomerStatusPending {
		return nil, fmt.Errorf("customer is not in pending status")
	}

	if err := uc.customerRepo.UpdateStatus(ctx, cmd.CustomerID, model.CustomerStatusActive, ""); err != nil {
		uc.logger.Error("failed to approve customer", err)
		return nil, fmt.Errorf("failed to approve customer")
	}

	customer.Status = model.CustomerStatusActive
	customer.UpdatedAt = time.Now()

	if err := uc.logAuditEvent(ctx, cmd.CustomerID, model.AuditActionApproved, cmd.ApprovedBy, ""); err != nil {
		uc.logger.Error("failed to log audit event", err)
	}

	return customerToResult(customer), nil
}

// RejectCustomer rejects a customer with a reason.
func (uc *CustomerUseCase) RejectCustomer(ctx context.Context, cmd command.RejectCustomerCommand) (*command.CustomerResult, error) {
	if err := validator.CustomerID(cmd.CustomerID); err != nil {
		return nil, err
	}
	if err := validator.NotEmpty(cmd.Reason, "reject reason"); err != nil {
		return nil, err
	}

	customer, err := uc.customerRepo.GetByID(ctx, cmd.CustomerID)
	if err != nil {
		uc.logger.Error("failed to get customer for rejection", err)
		return nil, fmt.Errorf("customer not found")
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	if customer.Status != model.CustomerStatusPending {
		return nil, fmt.Errorf("customer is not in pending status")
	}

	if err := uc.customerRepo.UpdateStatus(ctx, cmd.CustomerID, model.CustomerStatusRejected, cmd.Reason); err != nil {
		uc.logger.Error("failed to reject customer", err)
		return nil, fmt.Errorf("failed to reject customer")
	}

	customer.Status = model.CustomerStatusRejected
	customer.RejectReason = cmd.Reason
	customer.UpdatedAt = time.Now()

	if err := uc.logAuditEvent(ctx, cmd.CustomerID, model.AuditActionRejected, cmd.RejectedBy, cmd.Reason); err != nil {
		uc.logger.Error("failed to log audit event", err)
	}

	return customerToResult(customer), nil
}

// SuspendCustomer suspends a customer.
func (uc *CustomerUseCase) SuspendCustomer(ctx context.Context, cmd command.SuspendCustomerCommand) (*command.CustomerResult, error) {
	if err := validator.CustomerID(cmd.CustomerID); err != nil {
		return nil, err
	}

	customer, err := uc.customerRepo.GetByID(ctx, cmd.CustomerID)
	if err != nil {
		uc.logger.Error("failed to get customer for suspension", err)
		return nil, fmt.Errorf("customer not found")
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	if err := uc.customerRepo.UpdateStatus(ctx, cmd.CustomerID, model.CustomerStatusSuspended, cmd.Reason); err != nil {
		uc.logger.Error("failed to suspend customer", err)
		return nil, fmt.Errorf("failed to suspend customer")
	}

	customer.Status = model.CustomerStatusSuspended
	customer.UpdatedAt = time.Now()

	if err := uc.logAuditEvent(ctx, cmd.CustomerID, model.AuditActionSuspended, cmd.SuspendedBy, cmd.Reason); err != nil {
		uc.logger.Error("failed to log audit event", err)
	}

	return customerToResult(customer), nil
}

// RestoreCustomer restores a suspended customer (sets status to active).
func (uc *CustomerUseCase) RestoreCustomer(ctx context.Context, cmd command.RestoreCustomerCommand) (*command.CustomerResult, error) {
	if err := validator.CustomerID(cmd.CustomerID); err != nil {
		return nil, err
	}

	customer, err := uc.customerRepo.GetByID(ctx, cmd.CustomerID)
	if err != nil {
		uc.logger.Error("failed to get customer for restoration", err)
		return nil, fmt.Errorf("customer not found")
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	if customer.Status != model.CustomerStatusSuspended {
		return nil, fmt.Errorf("customer is not in suspended status")
	}

	if err := uc.customerRepo.UpdateStatus(ctx, cmd.CustomerID, model.CustomerStatusActive, ""); err != nil {
		uc.logger.Error("failed to restore customer", err)
		return nil, fmt.Errorf("failed to restore customer")
	}

	customer.Status = model.CustomerStatusActive
	customer.UpdatedAt = time.Now()

	if err := uc.logAuditEvent(ctx, cmd.CustomerID, model.AuditActionRestored, cmd.RestoredBy, ""); err != nil {
		uc.logger.Error("failed to log audit event", err)
	}

	return customerToResult(customer), nil
}

// logAuditEvent creates an audit log entry for a customer action.
func (uc *CustomerUseCase) logAuditEvent(ctx context.Context, customerID string, action model.AuditAction, performedBy, details string) error {
	if uc.auditLogRepo == nil {
		return nil
	}
	log := &model.AuditLog{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		Action:      action,
		PerformedBy: performedBy,
		Details:     details,
		CreatedAt:   time.Now(),
	}
	return uc.auditLogRepo.Insert(ctx, log)
}

// GetAuditLogs retrieves paginated audit logs for a customer.
func (uc *CustomerUseCase) GetAuditLogs(ctx context.Context, query command.ListAuditLogsQuery) (*command.ListAuditLogsResult, error) {
	if err := validator.NotEmpty(query.CustomerID, "customer ID"); err != nil {
		return nil, err
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	entries, total, err := uc.auditLogRepo.ListByCustomerID(ctx, query.CustomerID, query.Page, query.PageSize)
	if err != nil {
		uc.logger.Error("failed to list audit logs", err)
		return nil, fmt.Errorf("failed to list audit logs")
	}

	results := make([]*command.AuditLogResult, 0, len(entries))
	for _, e := range entries {
		results = append(results, &command.AuditLogResult{
			ID:          e.ID,
			CustomerID:  e.CustomerID,
			Action:      string(e.Action),
			PerformedBy: e.PerformedBy,
			Details:     e.Details,
			CreatedAt:   e.CreatedAt.Format(time.RFC3339),
		})
	}

	return &command.ListAuditLogsResult{
		Entries:    results,
		TotalCount: total,
		Page:       query.Page,
		PageSize:   query.PageSize,
	}, nil
}

func customerToResult(c *model.Customer) *command.CustomerResult {
	return &command.CustomerResult{
		ID:             c.ID,
		CompanyName:    c.CompanyName,
		Email:          c.Email,
		Phone:          c.Phone,
		VatNumber:      c.VatNumber,
		Address:        c.Address,
		City:           c.City,
		PostalCode:     c.PostalCode,
		Country:        c.Country,
		Status:         string(c.Status),
		CompanyAdminID: c.CompanyAdminID,
		RejectReason:   c.RejectReason,
		CreatedAt:      c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      c.UpdatedAt.Format(time.RFC3339),
	}
}
