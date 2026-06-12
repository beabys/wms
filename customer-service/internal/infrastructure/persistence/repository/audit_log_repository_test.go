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

func TestAuditLogRepository_InsertAndList(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuditLogRepository(testLogger(), db)

	customerID := createTestCustomer(t, db)

	// Insert audit log
	log := &model.AuditLog{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		Action:      model.AuditActionApproved,
		PerformedBy: "admin-123",
		Details:     "",
		CreatedAt:   time.Now(),
	}

	err := repo.Insert(context.Background(), log)
	require.NoError(t, err)

	// List audit logs
	logs, total, err := repo.ListByCustomerID(context.Background(), customerID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, logs, 1)
	assert.Equal(t, log.ID, logs[0].ID)
	assert.Equal(t, log.Action, logs[0].Action)
	assert.Equal(t, log.PerformedBy, logs[0].PerformedBy)
}

func TestAuditLogRepository_ListMultiple(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuditLogRepository(testLogger(), db)

	customerID := createTestCustomer(t, db)

	// Insert multiple audit logs
	for i := 0; i < 3; i++ {
		log := &model.AuditLog{
			ID:          uuid.New().String(),
			CustomerID:  customerID,
			Action:      model.AuditActionApproved,
			PerformedBy: "admin-123",
			CreatedAt:   time.Now(),
		}
		require.NoError(t, repo.Insert(context.Background(), log))
	}

	// List with pagination
	logs, total, err := repo.ListByCustomerID(context.Background(), customerID, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, logs, 2)

	// Page 2
	logs, total, err = repo.ListByCustomerID(context.Background(), customerID, 2, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, logs, 1)
}

func TestAuditLogRepository_ListEmpty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuditLogRepository(testLogger(), db)

	logs, total, err := repo.ListByCustomerID(context.Background(), "nonexistent-id", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, logs)
}

func TestAuditLogRepository_ListWithDefaults(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuditLogRepository(testLogger(), db)

	logs, total, err := repo.ListByCustomerID(context.Background(), "nonexistent", 0, 0)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, logs)
}
