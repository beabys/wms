package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/beabys/wms/inventory-service/internal/application/inventory/command"
	"github.com/beabys/wms/inventory-service/internal/application/inventory/repository"
	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// NewInventoryService creates a new InventoryService.
func NewInventoryService(
	productRepo repository.ProductRepository,
	stockRepo repository.StockRepository,
	binLocationRepo repository.BinLocationRepository,
	events EventPublisher,
) *InventoryService {
	return &InventoryService{
		productRepo:     productRepo,
		stockRepo:       stockRepo,
		binLocationRepo: binLocationRepo,
		events:          events,
	}
}

// CreateProduct handles CreateProductCommand.
func (s *InventoryService) CreateProduct(ctx context.Context, cmd command.CreateProductCommand) (*command.ProductResult, error) {
	id := newID()

	product, err := model.RegisterProduct(id, cmd.SKU, cmd.Name, cmd.Description, cmd.Category, cmd.Unit, cmd.WeightKg, cmd.LowStockThreshold)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}

	return &command.ProductResult{Product: product}, nil
}

// GetProduct handles GetProductQuery.
func (s *InventoryService) GetProduct(ctx context.Context, qry command.GetProductQuery) (*command.ProductResult, error) {
	product, err := s.productRepo.GetByID(ctx, qry.ProductID)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return &command.ProductResult{Product: product}, nil
}

// ListProducts handles ListProductsQuery.
func (s *InventoryService) ListProducts(ctx context.Context, qry command.ListProductsQuery) (*command.ListProductsResult, error) {
	products, nextToken, err := s.productRepo.List(ctx, qry.Category, qry.Search, qry.PageSize, qry.PageToken)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	return &command.ListProductsResult{
		Products:      products,
		NextPageToken: nextToken,
	}, nil
}

// UpdateProduct handles UpdateProductCommand.
func (s *InventoryService) UpdateProduct(ctx context.Context, cmd command.UpdateProductCommand) (*command.ProductResult, error) {
	product, err := s.productRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("get product for update: %w", err)
	}

	if err := product.Update(cmd.Name, cmd.Description, cmd.Category, cmd.Unit, cmd.WeightKg); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("persist product update: %w", err)
	}

	return &command.ProductResult{Product: product}, nil
}

// ArchiveProduct handles ArchiveProductCommand.
func (s *InventoryService) ArchiveProduct(ctx context.Context, cmd command.ArchiveProductCommand) (*command.ProductResult, error) {
	product, err := s.productRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("get product for archive: %w", err)
	}

	product.Archive()

	if err := s.productRepo.Archive(ctx, product.ID); err != nil {
		return nil, fmt.Errorf("persist product archive: %w", err)
	}

	return &command.ProductResult{Product: product}, nil
}

// AddStock handles AddStockCommand.
func (s *InventoryService) AddStock(ctx context.Context, cmd command.AddStockCommand) (*command.StockResult, error) {
	id := newID()

	entry, err := model.NewStockEntry(id, cmd.ProductID, cmd.BinLocationID, cmd.Quantity, cmd.LotNumber, cmd.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("add stock: %w", err)
	}

	if err := s.stockRepo.Save(ctx, entry); err != nil {
		return nil, fmt.Errorf("save stock entry: %w", err)
	}

	return &command.StockResult{StockEntry: entry}, nil
}

// AdjustStock handles AdjustStockCommand.
func (s *InventoryService) AdjustStock(ctx context.Context, cmd command.AdjustStockCommand) (*command.StockResult, error) {
	entry, err := s.stockRepo.GetByID(ctx, cmd.StockEntryID)
	if err != nil {
		return nil, fmt.Errorf("get stock entry for adjust: %w", err)
	}

	if err := entry.Adjust(cmd.Delta, cmd.Reason); err != nil {
		return nil, fmt.Errorf("adjust stock: %w", err)
	}

	if err := s.stockRepo.UpdateQuantity(ctx, entry); err != nil {
		return nil, fmt.Errorf("persist stock adjustment: %w", err)
	}

	if err := s.events.PublishStockAdjusted(ctx, entry, cmd.Reason); err != nil {
		return nil, fmt.Errorf("publish stock adjusted event: %w", err)
	}

	// Check low stock after adjustment
	if entry.IsLowStock(0) { // if available is negative, alert
		if err := s.events.PublishLowStockAlert(ctx, entry, 0); err != nil {
			return nil, fmt.Errorf("publish low stock alert: %w", err)
		}
	}

	return &command.StockResult{StockEntry: entry}, nil
}

// GetStock handles GetStockQuery.
func (s *InventoryService) GetStock(ctx context.Context, qry command.GetStockQuery) (*command.StockResult, error) {
	entry, err := s.stockRepo.GetByID(ctx, qry.StockEntryID)
	if err != nil {
		return nil, fmt.Errorf("get stock: %w", err)
	}
	return &command.StockResult{StockEntry: entry}, nil
}

// ListStock handles ListStockQuery.
func (s *InventoryService) ListStock(ctx context.Context, qry command.ListStockQuery) (*command.ListStockResult, error) {
	entries, nextToken, err := s.stockRepo.List(ctx, qry.ProductID, qry.BinLocationID, qry.Status, qry.PageSize, qry.PageToken)
	if err != nil {
		return nil, fmt.Errorf("list stock: %w", err)
	}
	return &command.ListStockResult{
		StockEntries:  entries,
		NextPageToken: nextToken,
	}, nil
}

// ReserveStock handles ReserveStockCommand.
func (s *InventoryService) ReserveStock(ctx context.Context, cmd command.ReserveStockCommand) (*command.StockResult, error) {
	entry, err := s.stockRepo.GetByID(ctx, cmd.StockEntryID)
	if err != nil {
		return nil, fmt.Errorf("get stock entry for reserve: %w", err)
	}

	if err := entry.Reserve(cmd.Quantity); err != nil {
		return nil, fmt.Errorf("reserve stock: %w", err)
	}

	if err := s.stockRepo.UpdateQuantity(ctx, entry); err != nil {
		return nil, fmt.Errorf("persist stock reserve: %w", err)
	}

	return &command.StockResult{StockEntry: entry}, nil
}

// ReleaseStock handles ReleaseStockCommand.
func (s *InventoryService) ReleaseStock(ctx context.Context, cmd command.ReleaseStockCommand) (*command.StockResult, error) {
	entry, err := s.stockRepo.GetByID(ctx, cmd.StockEntryID)
	if err != nil {
		return nil, fmt.Errorf("get stock entry for release: %w", err)
	}

	if err := entry.Release(cmd.Quantity); err != nil {
		return nil, fmt.Errorf("release stock: %w", err)
	}

	if err := s.stockRepo.UpdateQuantity(ctx, entry); err != nil {
		return nil, fmt.Errorf("persist stock release: %w", err)
	}

	return &command.StockResult{StockEntry: entry}, nil
}

// CreateBinLocation handles CreateBinLocationCommand.
func (s *InventoryService) CreateBinLocation(ctx context.Context, cmd command.CreateBinLocationCommand) (*command.BinLocationResult, error) {
	id := newID()

	location, err := model.NewBinLocation(id, cmd.WarehouseZone, cmd.Aisle, cmd.Rack, cmd.Shelf)
	if err != nil {
		return nil, fmt.Errorf("create bin location: %w", err)
	}

	if err := s.binLocationRepo.Save(ctx, location); err != nil {
		return nil, fmt.Errorf("save bin location: %w", err)
	}

	return &command.BinLocationResult{BinLocation: location}, nil
}

// ListBinLocations handles ListBinLocationsQuery.
func (s *InventoryService) ListBinLocations(ctx context.Context, qry command.ListBinLocationsQuery) (*command.ListBinLocationsResult, error) {
	locations, nextToken, err := s.binLocationRepo.List(ctx, qry.PageSize, qry.PageToken)
	if err != nil {
		return nil, fmt.Errorf("list bin locations: %w", err)
	}
	return &command.ListBinLocationsResult{
		BinLocations:  locations,
		NextPageToken: nextToken,
	}, nil
}

// newID generates a unique ID using UUID.
func newID() string {
	return uuid.NewString()
}

// Ensure InventoryService satisfies the interface.
var _ InventoryServiceHandler = (*InventoryService)(nil)
