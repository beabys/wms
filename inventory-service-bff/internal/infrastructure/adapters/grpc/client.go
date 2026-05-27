package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	inventoryv1 "github.com/beabys/wms/proto/gen/go/inventory/v1"
)

// withToken creates a context with JWT in gRPC metadata.
func withToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

// Client is the gRPC client adapter for the inventory service.
type Client struct {
	conn   *grpc.ClientConn
	client inventoryv1.InventoryServiceClient
}

// NewClient creates a new gRPC client for the inventory service.
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial grpc: %w", err)
	}

	return &Client{
		conn:   conn,
		client: inventoryv1.NewInventoryServiceClient(conn),
	}, nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateProduct creates a new product.
func (c *Client) CreateProduct(ctx context.Context, sku, name, description, category, unit string, weightKg, lowStockThreshold float64, token string) (*inventoryv1.Product, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.CreateProduct(ctx, &inventoryv1.CreateProductRequest{
		Sku:               sku,
		Name:              name,
		Description:       description,
		Category:          category,
		Unit:              unit,
		WeightKg:          weightKg,
		LowStockThreshold: lowStockThreshold,
	})
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	return resp.Product, nil
}

// GetProduct retrieves a product by ID.
func (c *Client) GetProduct(ctx context.Context, id string, token string) (*inventoryv1.Product, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.GetProduct(ctx, &inventoryv1.GetProductRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return resp.Product, nil
}

// ListProducts lists products with filters.
func (c *Client) ListProducts(ctx context.Context, category, search string, pageSize int32, token string) ([]*inventoryv1.Product, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ListProducts(ctx, &inventoryv1.ListProductsRequest{
		Category: category,
		Search:   search,
		Pagination: &commonv1.Pagination{
			Limit: pageSize,
			Page:  1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	return resp.Products, nil
}

// UpdateProduct updates a product.
func (c *Client) UpdateProduct(ctx context.Context, id, sku, name, description, category, unit string, weightKg, lowStockThreshold float64, token string) (*inventoryv1.Product, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.UpdateProduct(ctx, &inventoryv1.UpdateProductRequest{
		Id:                id,
		Sku:               sku,
		Name:              name,
		Description:       description,
		Category:          category,
		Unit:              unit,
		WeightKg:          weightKg,
		LowStockThreshold: lowStockThreshold,
	})
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}
	return resp.Product, nil
}

// ArchiveProduct archives (soft-deletes) a product.
func (c *Client) ArchiveProduct(ctx context.Context, id string, token string) (*inventoryv1.Product, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ArchiveProduct(ctx, &inventoryv1.ArchiveProductRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("archive product: %w", err)
	}
	return resp.Product, nil
}

// AddStock adds stock to a bin location.
func (c *Client) AddStock(ctx context.Context, productID, binLocationID, lotNumber, expiryDate string, quantity float64, token string) (*inventoryv1.StockEntry, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.AddStock(ctx, &inventoryv1.AddStockRequest{
		ProductId:     productID,
		BinLocationId: binLocationID,
		Quantity:      quantity,
		LotNumber:     lotNumber,
		ExpiryDate:    expiryDate,
	})
	if err != nil {
		return nil, fmt.Errorf("add stock: %w", err)
	}
	return resp.StockEntry, nil
}

// AdjustStock adjusts stock quantity.
func (c *Client) AdjustStock(ctx context.Context, id string, delta float64, reason, notes string, token string) (*inventoryv1.StockEntry, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.AdjustStock(ctx, &inventoryv1.AdjustStockRequest{
		Id:     id,
		Delta:  delta,
		Reason: reason,
		Notes:  notes,
	})
	if err != nil {
		return nil, fmt.Errorf("adjust stock: %w", err)
	}
	return resp.StockEntry, nil
}

// GetStock retrieves a stock entry by ID.
func (c *Client) GetStock(ctx context.Context, id string, token string) (*inventoryv1.StockEntry, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.GetStock(ctx, &inventoryv1.GetStockRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("get stock: %w", err)
	}
	return resp.StockEntry, nil
}

// ListStock lists stock entries with filters.
func (c *Client) ListStock(ctx context.Context, productID, binLocationID, status string, pageSize int32, token string) ([]*inventoryv1.StockEntry, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ListStock(ctx, &inventoryv1.ListStockRequest{
		ProductId:     productID,
		BinLocationId: binLocationID,
		Status:        status,
		Pagination: &commonv1.Pagination{
			Limit: pageSize,
			Page:  1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("list stock: %w", err)
	}
	return resp.StockEntries, nil
}

// ReserveStock reserves stock quantity.
func (c *Client) ReserveStock(ctx context.Context, id string, quantity float64, token string) (*inventoryv1.StockEntry, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ReserveStock(ctx, &inventoryv1.ReserveStockRequest{
		Id:       id,
		Quantity: quantity,
	})
	if err != nil {
		return nil, fmt.Errorf("reserve stock: %w", err)
	}
	return resp.StockEntry, nil
}

// ReleaseStock releases reserved stock quantity.
func (c *Client) ReleaseStock(ctx context.Context, id string, quantity float64, token string) (*inventoryv1.StockEntry, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ReleaseStock(ctx, &inventoryv1.ReleaseStockRequest{
		Id:       id,
		Quantity: quantity,
	})
	if err != nil {
		return nil, fmt.Errorf("release stock: %w", err)
	}
	return resp.StockEntry, nil
}

// CreateBinLocation creates a new bin location.
func (c *Client) CreateBinLocation(ctx context.Context, warehouseZone, aisle, rack, shelf string, token string) (*inventoryv1.BinLocation, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.CreateBinLocation(ctx, &inventoryv1.CreateBinLocationRequest{
		WarehouseZone: warehouseZone,
		Aisle:         aisle,
		Rack:          rack,
		Shelf:         shelf,
	})
	if err != nil {
		return nil, fmt.Errorf("create bin location: %w", err)
	}
	return resp.BinLocation, nil
}

// ListBinLocations lists bin locations.
func (c *Client) ListBinLocations(ctx context.Context, pageSize int32, token string) ([]*inventoryv1.BinLocation, error) {
	ctx = withToken(ctx, token)
	resp, err := c.client.ListBinLocations(ctx, &inventoryv1.ListBinLocationsRequest{
		Pagination: &commonv1.Pagination{
			Limit: pageSize,
			Page:  1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("list bin locations: %w", err)
	}
	return resp.BinLocations, nil
}
