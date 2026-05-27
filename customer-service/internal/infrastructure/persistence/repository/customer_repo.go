package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// CustomerRepo implements repository.CustomerRepository using pgx.
type CustomerRepo struct {
	pool *pgxpool.Pool
}

// NewCustomerRepo creates a new PostgreSQL customer repository.
func NewCustomerRepo(pool *pgxpool.Pool) *CustomerRepo {
	return &CustomerRepo{pool: pool}
}

// Save inserts a new customer row.
func (r *CustomerRepo) Save(ctx context.Context, c *model.Customer) error {
	addrJSON, err := json.Marshal(c.Address)
	if err != nil {
		return fmt.Errorf("marshal address: %w", err)
	}

	query := `INSERT INTO customers (id, company_name, vat_number, address, status, rate_card_id, credit_balance, created_at, updated_at)
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = r.pool.Exec(ctx, query,
		c.ID, c.CompanyName, c.VATNumber, addrJSON, c.Status, c.RateCardID, c.CreditBalance, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("pg insert customer: %w", err)
	}
	return nil
}

// GetByID retrieves a customer by primary key.
func (r *CustomerRepo) GetByID(ctx context.Context, id string) (*model.Customer, error) {
	query := `SELECT id, company_name, vat_number, address, status, rate_card_id, credit_balance, created_at, updated_at
	           FROM customers WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	return scanCustomer(row)
}

// GetByEmail retrieves a customer by company email (first match from address json).
// For simplicity, this searches via a JSON field. In production, add a dedicated column.
func (r *CustomerRepo) GetByEmail(ctx context.Context, email string) (*model.Customer, error) {
	return nil, fmt.Errorf("get by email not implemented")
}

// List returns a paginated list of customers, optionally filtered by status.
func (r *CustomerRepo) List(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var countQuery string
	var dataQuery string
	var args []interface{}
	var countArgs []interface{}

	if status != "" {
		countQuery = `SELECT count(*) FROM customers WHERE status = $1`
		dataQuery = `SELECT id, company_name, vat_number, address, status, rate_card_id, credit_balance, created_at, updated_at
		             FROM customers WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, status, pageSize, offset)
		countArgs = append(countArgs, status)
	} else {
		countQuery = `SELECT count(*) FROM customers`
		dataQuery = `SELECT id, company_name, vat_number, address, status, rate_card_id, credit_balance, created_at, updated_at
		             FROM customers ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = append(args, pageSize, offset)
	}

	var total int32
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("pg count customers: %w", err)
	}

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("pg list customers: %w", err)
	}
	defer rows.Close()

	var customers []*model.Customer
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, 0, err
		}
		customers = append(customers, c)
	}
	if customers == nil {
		customers = []*model.Customer{}
	}
	return customers, total, nil
}

// Update modifies an existing customer row.
func (r *CustomerRepo) Update(ctx context.Context, c *model.Customer) error {
	addrJSON, err := json.Marshal(c.Address)
	if err != nil {
		return fmt.Errorf("marshal address: %w", err)
	}

	query := `UPDATE customers SET company_name=$2, vat_number=$3, address=$4, status=$5,
	           rate_card_id=$6, credit_balance=$7, updated_at=$8 WHERE id=$1`

	_, err = r.pool.Exec(ctx, query,
		c.ID, c.CompanyName, c.VATNumber, addrJSON, c.Status, c.RateCardID, c.CreditBalance, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("pg update customer: %w", err)
	}
	return nil
}

// rowScanner is satisfied by pgx.Row and pgx.Rows.
type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanCustomer(row rowScanner) (*model.Customer, error) {
	var (
		c         model.Customer
		addrJSON  []byte
		createdAt time.Time
		updatedAt time.Time
	)

	if err := row.Scan(
		&c.ID, &c.CompanyName, &c.VATNumber, &addrJSON, &c.Status, &c.RateCardID,
		&c.CreditBalance, &createdAt, &updatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("pg scan customer: %w", err)
	}

	if err := json.Unmarshal(addrJSON, &c.Address); err != nil {
		return nil, fmt.Errorf("unmarshal address: %w", err)
	}
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt

	return &c, nil
}
