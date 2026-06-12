package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/repository"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// auditLogRow maps to the customer_audit_log table in PostgreSQL.
type auditLogRow struct {
	ID          string    `db:"id"`
	CustomerID  string    `db:"customer_id"`
	Action      string    `db:"action"`
	PerformedBy string    `db:"performed_by"`
	Details     string    `db:"details"`
	CreatedAt   time.Time `db:"created_at"`
}

// AuditLogRepository implements repository.AuditLogRepository using PostgreSQL.
type AuditLogRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewAuditLogRepository creates a new AuditLogRepository.
func NewAuditLogRepository(log logger.Logger, d database.Database) *AuditLogRepository {
	pg, ok := d.GetDBImpl().(*database.Postgres)
	if !ok {
		return &AuditLogRepository{db: nil, log: log}
	}
	return &AuditLogRepository{db: pg.DB, log: log}
}

// Insert adds a new audit log entry.
func (r *AuditLogRepository) Insert(ctx context.Context, log *model.AuditLog) error {
	query := `INSERT INTO customer_audit_log (id, customer_id, action, performed_by, details, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.CustomerID,
		string(log.Action),
		log.PerformedBy,
		log.Details,
		log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}
	return nil
}

// ListByCustomerID retrieves paginated audit logs for a customer.
func (r *AuditLogRepository) ListByCustomerID(ctx context.Context, customerID string, page, pageSize int) ([]*model.AuditLog, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Count
	var total int
	countQuery := "SELECT COUNT(*) FROM customer_audit_log WHERE customer_id = $1"
	if err := r.db.GetContext(ctx, &total, countQuery, customerID); err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Fetch with pagination
	offset := (page - 1) * pageSize
	dataQuery := `SELECT id, customer_id, action, performed_by, details, created_at
	              FROM customer_audit_log
	              WHERE customer_id = $1
	              ORDER BY created_at DESC
	              LIMIT $2 OFFSET $3`

	var rows []auditLogRow
	if err := r.db.SelectContext(ctx, &rows, dataQuery, customerID, pageSize, offset); err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}

	logs := make([]*model.AuditLog, 0, len(rows))
	for i := range rows {
		logs = append(logs, rowToAuditLog(&rows[i]))
	}

	return logs, total, nil
}

// Ensure AuditLogRepository implements repository.AuditLogRepository.
var _ repository.AuditLogRepository = (*AuditLogRepository)(nil)

func rowToAuditLog(row *auditLogRow) *model.AuditLog {
	return &model.AuditLog{
		ID:          row.ID,
		CustomerID:  row.CustomerID,
		Action:      model.AuditAction(row.Action),
		PerformedBy: row.PerformedBy,
		Details:     row.Details,
		CreatedAt:   row.CreatedAt,
	}
}
