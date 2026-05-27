package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// ProductRepository is the PostgreSQL implementation of application ProductRepository.
type ProductRepository struct {
	pool *pgxpool.Pool
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

// GetByID retrieves a product by ID.
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	query := `
		SELECT id, sku, name, description, category, unit, weight_kg, is_active, low_stock_threshold, created_at, updated_at
		FROM products WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	p := &model.Product{}
	err := row.Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.Unit,
		&p.WeightKg, &p.IsActive, &p.LowStockThreshold, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("product not found: %s", id)
		}
		return nil, fmt.Errorf("scan product: %w", err)
	}
	return p, nil
}

// Save persists a new product.
func (r *ProductRepository) Save(ctx context.Context, p *model.Product) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO products (id, sku, name, description, category, unit, weight_kg, is_active, low_stock_threshold, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, p.ID, p.SKU, p.Name, p.Description, p.Category, p.Unit, p.WeightKg, p.IsActive, p.LowStockThreshold, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

// Update updates an existing product.
func (r *ProductRepository) Update(ctx context.Context, p *model.Product) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE products SET name = $1, description = $2, category = $3, unit = $4, weight_kg = $5, updated_at = $6
		WHERE id = $7
	`, p.Name, p.Description, p.Category, p.Unit, p.WeightKg, p.UpdatedAt, p.ID)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

// List returns products with optional filters.
func (r *ProductRepository) List(ctx context.Context, category, search string, pageSize int32, pageToken string) ([]*model.Product, string, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	args := make([]any, 0, 4)
	where := "WHERE 1=1"
	argIdx := 1

	if category != "" {
		where += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}
	if search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT id, sku, name, description, category, unit, weight_kg, is_active, low_stock_threshold, created_at, updated_at
		FROM products %s
		ORDER BY created_at DESC
		LIMIT $%d
	`, where, argIdx)
	args = append(args, pageSize+1)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.Unit,
			&p.WeightKg, &p.IsActive, &p.LowStockThreshold, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, "", fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}

	var nextToken string
	if len(products) > int(pageSize) {
		products = products[:pageSize]
		nextToken = products[len(products)-1].ID
	}

	return products, nextToken, nil
}

// Archive soft-deletes a product.
func (r *ProductRepository) Archive(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE products SET is_active = false, updated_at = $1 WHERE id = $2
	`, time.Now(), id)
	if err != nil {
		return fmt.Errorf("archive product: %w", err)
	}
	return nil
}
