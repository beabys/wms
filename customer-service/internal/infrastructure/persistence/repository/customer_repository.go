package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/repository"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// customerRow maps to the customers table in PostgreSQL.
type customerRow struct {
	ID             string    `db:"id"`
	CompanyName    string    `db:"company_name"`
	Email          string    `db:"email"`
	Phone          string    `db:"phone"`
	VatNumber      string    `db:"vat_number"`
	Address        string    `db:"address"`
	City           string    `db:"city"`
	PostalCode     string    `db:"postal_code"`
	Country        string    `db:"country"`
	Status         string    `db:"status"`
	CompanyAdminID string    `db:"company_admin_id"`
	RejectReason   string    `db:"reject_reason"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// CustomerRepository implements repository.CustomerRepository using PostgreSQL.
type CustomerRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(log logger.Logger, d database.Database) *CustomerRepository {
	pg, ok := d.GetDBImpl().(*database.Postgres)
	if !ok {
		return &CustomerRepository{db: nil, log: log}
	}
	return &CustomerRepository{db: pg.DB, log: log}
}

// Create inserts a new customer into the database.
func (r *CustomerRepository) Create(ctx context.Context, customer *model.Customer) error {
	query := `INSERT INTO customers (id, company_name, email, phone, vat_number, address, city, postal_code, country, status, company_admin_id, reject_reason, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.CompanyName,
		customer.Email,
		customer.Phone,
		customer.VatNumber,
		customer.Address,
		customer.City,
		customer.PostalCode,
		customer.Country,
		string(customer.Status),
		customer.CompanyAdminID,
		customer.RejectReason,
		customer.CreatedAt,
		customer.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

// GetByID retrieves a customer by their ID.
func (r *CustomerRepository) GetByID(ctx context.Context, id string) (*model.Customer, error) {
	var row customerRow
	err := r.db.GetContext(ctx, &row, "SELECT * FROM customers WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer by id: %w", err)
	}
	return rowToCustomer(&row), nil
}

// GetByEmail retrieves a customer by their email.
func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*model.Customer, error) {
	var row customerRow
	err := r.db.GetContext(ctx, &row, "SELECT * FROM customers WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer by email: %w", err)
	}
	return rowToCustomer(&row), nil
}

// Update updates an existing customer.
func (r *CustomerRepository) Update(ctx context.Context, customer *model.Customer) error {
	query := `UPDATE customers SET company_name=$1, phone=$2, vat_number=$3, address=$4, city=$5, postal_code=$6, country=$7, updated_at=$8 WHERE id=$9`

	result, err := r.db.ExecContext(ctx, query,
		customer.CompanyName,
		customer.Phone,
		customer.VatNumber,
		customer.Address,
		customer.City,
		customer.PostalCode,
		customer.Country,
		customer.UpdatedAt,
		customer.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}

// List retrieves customers with filtering and pagination.
func (r *CustomerRepository) List(ctx context.Context, page, pageSize int, status string) ([]*model.Customer, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Build WHERE clause
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customers %s", where)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	// Fetch with pagination
	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf("SELECT * FROM customers %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", where, argIdx, argIdx+1)
	fetchArgs := append(args, pageSize, offset)

	var rows []customerRow
	if err := r.db.SelectContext(ctx, &rows, dataQuery, fetchArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}

	customers := make([]*model.Customer, 0, len(rows))
	for i := range rows {
		customers = append(customers, rowToCustomer(&rows[i]))
	}

	return customers, total, nil
}

// UpdateStatus changes a customer's status and optionally sets reject reason.
func (r *CustomerRepository) UpdateStatus(ctx context.Context, id string, status model.CustomerStatus, reason string) error {
	now := time.Now()
	query := "UPDATE customers SET status=$1, updated_at=$2"
	args := []interface{}{string(status), now}
	argIdx := 3

	if reason != "" {
		query += fmt.Sprintf(", reject_reason=$%d", argIdx)
		args = append(args, reason)
		argIdx++
	}

	query += fmt.Sprintf(" WHERE id=$%d", argIdx)
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update customer status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}

// FindByAdminID retrieves a customer by company_admin_id.
func (r *CustomerRepository) FindByAdminID(ctx context.Context, adminUserID string) (*model.Customer, error) {
	var row customerRow
	err := r.db.GetContext(ctx, &row, "SELECT * FROM customers WHERE company_admin_id = $1", adminUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer by admin id: %w", err)
	}
	return rowToCustomer(&row), nil
}

// Ensure CustomerRepository implements repository.CustomerRepository.
var _ repository.CustomerRepository = (*CustomerRepository)(nil)

func rowToCustomer(row *customerRow) *model.Customer {
	return &model.Customer{
		ID:             row.ID,
		CompanyName:    row.CompanyName,
		Email:          row.Email,
		Phone:          row.Phone,
		VatNumber:      row.VatNumber,
		Address:        row.Address,
		City:           row.City,
		PostalCode:     row.PostalCode,
		Country:        row.Country,
		Status:         model.CustomerStatus(row.Status),
		CompanyAdminID: row.CompanyAdminID,
		RejectReason:   row.RejectReason,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
}
