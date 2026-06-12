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

// companyRoleRow maps to the company_roles table in PostgreSQL.
type companyRoleRow struct {
	ID          string    `db:"id"`
	CustomerID  string    `db:"customer_id"`
	Name        string    `db:"name"`
	Permissions []string  `db:"permissions"`
	IsDefault   bool      `db:"is_default"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// CompanyRoleRepository implements repository.CompanyRoleRepository using PostgreSQL.
type CompanyRoleRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewCompanyRoleRepository creates a new CompanyRoleRepository.
func NewCompanyRoleRepository(log logger.Logger, d database.Database) *CompanyRoleRepository {
	pg, ok := d.GetDBImpl().(*database.Postgres)
	if !ok {
		return &CompanyRoleRepository{db: nil, log: log}
	}
	return &CompanyRoleRepository{db: pg.DB, log: log}
}

// Create inserts a new company role into the database.
func (r *CompanyRoleRepository) Create(ctx context.Context, role *model.CompanyRole) error {
	query := `INSERT INTO company_roles (id, customer_id, name, permissions, is_default, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		role.ID,
		role.CustomerID,
		role.Name,
		role.Permissions,
		role.IsDefault,
		role.CreatedAt,
		role.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create company role: %w", err)
	}
	return nil
}

// ListByCustomer retrieves all roles for a customer.
func (r *CompanyRoleRepository) ListByCustomer(ctx context.Context, customerID string) ([]*model.CompanyRole, error) {
	var rows []companyRoleRow
	err := r.db.SelectContext(ctx, &rows, "SELECT * FROM company_roles WHERE customer_id = $1 ORDER BY created_at", customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list company roles: %w", err)
	}

	roles := make([]*model.CompanyRole, 0, len(rows))
	for i := range rows {
		roles = append(roles, rowToCompanyRole(&rows[i]))
	}
	return roles, nil
}

// GetUserPermissions retrieves aggregated permissions for a user in a customer.
func (r *CompanyRoleRepository) GetUserPermissions(ctx context.Context, userID, customerID string) ([]string, error) {
	query := `SELECT DISTINCT cr.permissions
	          FROM company_user_roles cur
	          JOIN company_roles cr ON cur.role_id = cr.id
	          WHERE cur.user_id = $1 AND cur.customer_id = $2`

	var permissionsList []string
	var result []string

	rows, err := r.db.QueryContext(ctx, query, userID, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&permissionsList); err != nil {
			return nil, fmt.Errorf("failed to scan permissions: %w", err)
		}
		result = append(result, permissionsList...)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	// Deduplicate
	seen := make(map[string]bool)
	unique := make([]string, 0, len(result))
	for _, p := range result {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}

	return unique, nil
}

// AssignUserRole assigns a role to a user within a customer.
func (r *CompanyRoleRepository) AssignUserRole(ctx context.Context, userID, roleID, customerID, assignedBy string) error {
	query := `INSERT INTO company_user_roles (user_id, role_id, customer_id, assigned_by, created_at)
	          VALUES ($1, $2, $3, $4, $5)
	          ON CONFLICT (user_id, role_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		roleID,
		customerID,
		assignedBy,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to assign user role: %w", err)
	}
	return nil
}

// Ensure CompanyRoleRepository implements repository.CompanyRoleRepository.
var _ repository.CompanyRoleRepository = (*CompanyRoleRepository)(nil)

func rowToCompanyRole(row *companyRoleRow) *model.CompanyRole {
	return &model.CompanyRole{
		ID:          row.ID,
		CustomerID:  row.CustomerID,
		Name:        row.Name,
		Permissions: row.Permissions,
		IsDefault:   row.IsDefault,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
