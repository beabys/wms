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

func TestCustomerRepositoryCreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(testLogger(), db)

	customer := &model.Customer{
		ID:             uuid.New().String(),
		CompanyName:    "Integration Test Corp",
		Email:          "inttest@" + uuid.New().String() + ".com",
		Status:         model.CustomerStatusPending,
		CompanyAdminID: uuid.New().String(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), customer.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, customer.ID, got.ID)
	assert.Equal(t, customer.CompanyName, got.CompanyName)
	assert.Equal(t, customer.Email, got.Email)
	assert.Equal(t, model.CustomerStatusPending, got.Status)
}

func TestCustomerRepositoryGetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(testLogger(), db)

	customer := &model.Customer{
		ID:             uuid.New().String(),
		CompanyName:    "Email Test",
		Email:          "emailtest@" + uuid.New().String() + ".com",
		Status:         model.CustomerStatusActive,
		CompanyAdminID: uuid.New().String(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	got, err := repo.GetByEmail(context.Background(), customer.Email)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, customer.ID, got.ID)
}

func TestCustomerRepositoryUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(testLogger(), db)

	customer := &model.Customer{
		ID:             uuid.New().String(),
		CompanyName:    "Original Name",
		Email:          "update@" + uuid.New().String() + ".com",
		Status:         model.CustomerStatusPending,
		CompanyAdminID: uuid.New().String(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	customer.CompanyName = "Updated Name"
	customer.UpdatedAt = time.Now()
	err = repo.Update(context.Background(), customer)
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), customer.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.CompanyName)
}

func TestCustomerRepositoryUpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(testLogger(), db)

	customer := &model.Customer{
		ID:             uuid.New().String(),
		CompanyName:    "Status Test",
		Email:          "status@" + uuid.New().String() + ".com",
		Status:         model.CustomerStatusPending,
		CompanyAdminID: uuid.New().String(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)

	err = repo.UpdateStatus(context.Background(), customer.ID, model.CustomerStatusActive, "")
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), customer.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CustomerStatusActive, got.Status)
}

func TestCustomerRepositoryList(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCustomerRepository(testLogger(), db)

	for i := 0; i < 3; i++ {
		c := &model.Customer{
			ID:             uuid.New().String(),
			CompanyName:    "List Test " + uuid.New().String(),
			Email:          "list" + uuid.New().String() + "@test.com",
			Status:         model.CustomerStatusActive,
			CompanyAdminID: uuid.New().String(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		require.NoError(t, repo.Create(context.Background(), c))
	}

	customers, total, err := repo.List(context.Background(), 1, 10, "")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, 3)
	assert.GreaterOrEqual(t, len(customers), 3)
}
