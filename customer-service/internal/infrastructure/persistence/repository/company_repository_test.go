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

func TestCompanyRepositoryCreateAndGet(t *testing.T) {
	db := setupTestDB(t)

	// Create a customer first
	customerID := createTestCustomer(t, db)

	company := &model.Company{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Name:       "Test Company",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	companyRepo := NewCompanyRepository(testLogger(), db)
	err := companyRepo.Create(context.Background(), company)
	require.NoError(t, err)

	got, err := companyRepo.GetByID(context.Background(), company.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, company.Name, got.Name)
	assert.Equal(t, company.CustomerID, got.CustomerID)
}

func TestCompanyRepositoryGetByCustomerID(t *testing.T) {
	db := setupTestDB(t)
	customerID := createTestCustomer(t, db)

	company := &model.Company{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Name:       "Company By Customer",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	companyRepo := NewCompanyRepository(testLogger(), db)
	err := companyRepo.Create(context.Background(), company)
	require.NoError(t, err)

	got, err := companyRepo.GetByCustomerID(context.Background(), customerID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, company.ID, got.ID)
}
