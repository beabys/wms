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

// companyRow maps to the companies table in PostgreSQL.
type companyRow struct {
	ID         string    `db:"id"`
	CustomerID string    `db:"customer_id"`
	Name       string    `db:"name"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// CompanyRepository implements repository.CompanyRepository using PostgreSQL.
type CompanyRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewCompanyRepository creates a new CompanyRepository.
func NewCompanyRepository(log logger.Logger, d database.Database) *CompanyRepository {
	pg, ok := d.GetDBImpl().(*database.Postgres)
	if !ok {
		return &CompanyRepository{db: nil, log: log}
	}
	return &CompanyRepository{db: pg.DB, log: log}
}

// Create inserts a new company into the database.
func (r *CompanyRepository) Create(ctx context.Context, company *model.Company) error {
	query := `INSERT INTO companies (id, customer_id, name, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		company.ID,
		company.CustomerID,
		company.Name,
		company.CreatedAt,
		company.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

// GetByID retrieves a company by its ID.
func (r *CompanyRepository) GetByID(ctx context.Context, id string) (*model.Company, error) {
	var row companyRow
	err := r.db.GetContext(ctx, &row, "SELECT * FROM companies WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company by id: %w", err)
	}
	return rowToCompany(&row), nil
}

// GetByCustomerID retrieves a company by customer ID.
func (r *CompanyRepository) GetByCustomerID(ctx context.Context, customerID string) (*model.Company, error) {
	var row companyRow
	err := r.db.GetContext(ctx, &row, "SELECT * FROM companies WHERE customer_id = $1", customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company by customer id: %w", err)
	}
	return rowToCompany(&row), nil
}

// Ensure CompanyRepository implements repository.CompanyRepository.
var _ repository.CompanyRepository = (*CompanyRepository)(nil)

func rowToCompany(row *companyRow) *model.Company {
	return &model.Company{
		ID:         row.ID,
		CustomerID: row.CustomerID,
		Name:       row.Name,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
