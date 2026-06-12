package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	repomocks "github.com/beabys/wms/customer-service/mocks/application/customer/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAssignCompanyRole_Success(t *testing.T) {
	companyRepo := repomocks.NewCompanyRepository(t)
	roleRepo := repomocks.NewCompanyRoleRepository(t)

	companyRepo.On("GetByCustomerID", mock.Anything, "cust-1").Return(&model.Company{
		ID:         "comp-1",
		CustomerID: "cust-1",
		Name:       "Test Corp",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil)
	roleRepo.On("ListByCustomer", mock.Anything, "cust-1").Return([]*model.CompanyRole{}, nil)
	roleRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.CompanyRole")).Return(nil)
	roleRepo.On("AssignUserRole", mock.Anything, "user-1", mock.AnythingOfType("string"), "cust-1", "admin-1").Return(nil)

	uc := &CompanyUseCase{
		logger:      &testLogger{},
		companyRepo: companyRepo,
		roleRepo:    roleRepo,
	}

	err := uc.AssignCompanyRole(context.Background(), command.AssignCompanyRoleCommand{
		CustomerID:  "cust-1",
		UserID:      "user-1",
		RoleName:    "manager",
		Permissions: []string{"orders:read", "inventory:read"},
		AssignedBy:  "admin-1",
	})
	require.NoError(t, err)
}

func TestAssignCompanyRole_EmptyCustomerID(t *testing.T) {
	uc := &CompanyUseCase{logger: &testLogger{}}
	err := uc.AssignCompanyRole(context.Background(), command.AssignCompanyRoleCommand{
		CustomerID: "",
		UserID:     "user-1",
		RoleName:   "admin",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer ID is required")
}

func TestAssignCompanyRole_EmptyUserID(t *testing.T) {
	uc := &CompanyUseCase{logger: &testLogger{}}
	err := uc.AssignCompanyRole(context.Background(), command.AssignCompanyRoleCommand{
		CustomerID: "cust-1",
		UserID:     "",
		RoleName:   "admin",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID is required")
}

func TestAssignCompanyRole_EmptyRoleName(t *testing.T) {
	uc := &CompanyUseCase{logger: &testLogger{}}
	err := uc.AssignCompanyRole(context.Background(), command.AssignCompanyRoleCommand{
		CustomerID: "cust-1",
		UserID:     "user-1",
		RoleName:   "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role name is required")
}

func TestAssignCompanyRole_CustomerNotFound(t *testing.T) {
	companyRepo := repomocks.NewCompanyRepository(t)
	roleRepo := repomocks.NewCompanyRoleRepository(t)

	companyRepo.On("GetByCustomerID", mock.Anything, "nonexistent").Return(nil, nil)

	uc := &CompanyUseCase{
		logger:      &testLogger{},
		companyRepo: companyRepo,
		roleRepo:    roleRepo,
	}

	err := uc.AssignCompanyRole(context.Background(), command.AssignCompanyRoleCommand{
		CustomerID: "nonexistent",
		UserID:     "user-1",
		RoleName:   "admin",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer not found")
}

func TestListCompanyRoles_Success(t *testing.T) {
	roleRepo := repomocks.NewCompanyRoleRepository(t)
	now := time.Now()
	roles := []*model.CompanyRole{
		{
			ID:          "role-1",
			CustomerID:  "cust-1",
			Name:        "admin",
			Permissions: []string{"*"},
			IsDefault:   true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "role-2",
			CustomerID:  "cust-1",
			Name:        "viewer",
			Permissions: []string{"orders:read"},
			IsDefault:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	roleRepo.On("ListByCustomer", mock.Anything, "cust-1").Return(roles, nil)

	uc := &CompanyUseCase{
		logger:   &testLogger{},
		roleRepo: roleRepo,
	}

	rolesResult, err := uc.ListCompanyRoles(context.Background(), "cust-1")
	require.NoError(t, err)
	assert.Len(t, rolesResult, 2)

	// Verify mock expectations
	roleRepo.AssertExpectations(t)
}

func TestListCompanyRoles_EmptyCustomerID(t *testing.T) {
	uc := &CompanyUseCase{logger: &testLogger{}}
	_, err := uc.ListCompanyRoles(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer ID is required")
}

func TestGetUserPermissions_Success(t *testing.T) {
	roleRepo := repomocks.NewCompanyRoleRepository(t)

	roleRepo.On("GetUserPermissions", mock.Anything, "user-1", "cust-1").Return([]string{"orders:read", "orders:write"}, nil)

	uc := &CompanyUseCase{
		logger:   &testLogger{},
		roleRepo: roleRepo,
	}

	perms, err := uc.GetUserPermissions(context.Background(), "user-1", "cust-1")
	require.NoError(t, err)
	assert.Contains(t, perms, "orders:read")
	assert.Contains(t, perms, "orders:write")
}

func TestGetUserPermissions_EmptyUserID(t *testing.T) {
	uc := &CompanyUseCase{logger: &testLogger{}}
	_, err := uc.GetUserPermissions(context.Background(), "", "cust-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID is required")
}

func TestGetUserPermissions_EmptyCustomerID(t *testing.T) {
	uc := &CompanyUseCase{logger: &testLogger{}}
	_, err := uc.GetUserPermissions(context.Background(), "user-1", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer ID is required")
}
