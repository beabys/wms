package repository

import (
	"context"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// AuditLogRepository defines the interface for audit log persistence.
type AuditLogRepository interface {
	// Insert adds a new audit log entry.
	Insert(ctx context.Context, log *model.AuditLog) error

	// ListByCustomerID retrieves paginated audit logs for a customer.
	ListByCustomerID(ctx context.Context, customerID string, page, pageSize int) ([]*model.AuditLog, int, error)
}
