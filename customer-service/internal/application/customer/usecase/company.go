package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/repository"
	"github.com/beabys/wms/customer-service/internal/application/customer/validator"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	"github.com/beabys/wms/pkg/logger"
	"github.com/google/uuid"
)

// CompanyUseCase handles company role and permission operations.
type CompanyUseCase struct {
	logger    logger.Logger
	companyRepo repository.CompanyRepository
	roleRepo    repository.CompanyRoleRepository
}

// NewCompanyUseCase creates a new CompanyUseCase.
func NewCompanyUseCase(
	log logger.Logger,
	companyRepo repository.CompanyRepository,
	roleRepo repository.CompanyRoleRepository,
) *CompanyUseCase {
	return &CompanyUseCase{
		logger:      log,
		companyRepo: companyRepo,
		roleRepo:    roleRepo,
	}
}

// AssignCompanyRole assigns a role with permissions to a user within a company.
func (uc *CompanyUseCase) AssignCompanyRole(ctx context.Context, cmd command.AssignCompanyRoleCommand) error {
	if err := validator.NotEmpty(cmd.CustomerID, "customer ID"); err != nil {
		return err
	}
	if err := validator.NotEmpty(cmd.UserID, "user ID"); err != nil {
		return err
	}
	if err := validator.NotEmpty(cmd.RoleName, "role name"); err != nil {
		return err
	}

	// Validate the customer exists
	company, err := uc.companyRepo.GetByCustomerID(ctx, cmd.CustomerID)
	if err != nil {
		uc.logger.Error("failed to get company for role assignment", err)
		return fmt.Errorf("customer not found")
	}
	if company == nil {
		return fmt.Errorf("customer not found")
	}

	// Create or get existing role
	roles, err := uc.roleRepo.ListByCustomer(ctx, cmd.CustomerID)
	if err != nil {
		uc.logger.Error("failed to list roles", err)
		return fmt.Errorf("failed to assign role")
	}

	var role *model.CompanyRole
	for _, r := range roles {
		if r.Name == cmd.RoleName {
			role = r
			break
		}
	}

	if role == nil {
		now := time.Now()
		role = &model.CompanyRole{
			ID:          uuid.New().String(),
			CustomerID:  cmd.CustomerID,
			Name:        cmd.RoleName,
			Permissions: cmd.Permissions,
			IsDefault:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := uc.roleRepo.Create(ctx, role); err != nil {
			uc.logger.Error("failed to create role", err)
			return fmt.Errorf("failed to assign role")
		}
	}

	if err := uc.roleRepo.AssignUserRole(ctx, cmd.UserID, role.ID, cmd.CustomerID, cmd.AssignedBy); err != nil {
		uc.logger.Error("failed to assign user role", err)
		return fmt.Errorf("failed to assign role")
	}

	uc.logger.Info("role assigned to user",
		logger.LogField{Key: "user_id", Value: cmd.UserID},
		logger.LogField{Key: "role", Value: cmd.RoleName},
		logger.LogField{Key: "customer_id", Value: cmd.CustomerID},
	)

	return nil
}

// ListCompanyRoles returns all roles for a customer.
func (uc *CompanyUseCase) ListCompanyRoles(ctx context.Context, customerID string) ([]command.CompanyRoleResult, error) {
	if err := validator.NotEmpty(customerID, "customer ID"); err != nil {
		return nil, err
	}

	roles, err := uc.roleRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		uc.logger.Error("failed to list company roles", err)
		return nil, fmt.Errorf("failed to list roles")
	}

	results := make([]command.CompanyRoleResult, 0, len(roles))
	for _, r := range roles {
		results = append(results, command.CompanyRoleResult{
			ID:          r.ID,
			CustomerID:  r.CustomerID,
			Name:        r.Name,
			Permissions: r.Permissions,
			IsDefault:   r.IsDefault,
		})
	}

	return results, nil
}

// GetUserPermissions retrieves aggregated permissions for a user in a customer.
func (uc *CompanyUseCase) GetUserPermissions(ctx context.Context, userID, customerID string) ([]string, error) {
	if err := validator.NotEmpty(userID, "user ID"); err != nil {
		return nil, err
	}
	if err := validator.NotEmpty(customerID, "customer ID"); err != nil {
		return nil, err
	}

	permissions, err := uc.roleRepo.GetUserPermissions(ctx, userID, customerID)
	if err != nil {
		uc.logger.Error("failed to get user permissions", err)
		return nil, fmt.Errorf("failed to get permissions")
	}

	return permissions, nil
}
