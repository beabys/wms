//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompanyRoleRepositoryCreateAndList(t *testing.T) {
	db := setupTestDB(t)
	customerID := createTestCustomer(t, db)

	roleRepo := NewCompanyRoleRepository(testLogger(), db)

	role := &model.CompanyRole{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		Name:        "admin",
		Permissions: []string{"*"},
		IsDefault:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := roleRepo.Create(context.Background(), role)
	require.NoError(t, err)

	roles, err := roleRepo.ListByCustomer(context.Background(), customerID)
	require.NoError(t, err)
	assert.Len(t, roles, 1)
	assert.Equal(t, "admin", roles[0].Name)
}

func TestCompanyRoleRepositoryAssignAndGetPermissions(t *testing.T) {
	db := setupTestDB(t)
	customerID := createTestCustomer(t, db)
	userID := uuid.New().String()

	roleRepo := NewCompanyRoleRepository(testLogger(), db)

	role := &model.CompanyRole{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		Name:        "viewer",
		Permissions: []string{"orders:read", "inventory:read"},
		IsDefault:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := roleRepo.Create(context.Background(), role)
	require.NoError(t, err)

	err = roleRepo.AssignUserRole(context.Background(), userID, role.ID, customerID, "admin-1")
	require.NoError(t, err)

	permissions, err := roleRepo.GetUserPermissions(context.Background(), userID, customerID)
	require.NoError(t, err)
	assert.Contains(t, permissions, "orders:read")
	assert.Contains(t, permissions, "inventory:read")
}
