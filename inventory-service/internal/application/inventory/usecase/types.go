package usecase

import (
	"context"

	"github.com/beabys/wms/inventory-service/internal/application/inventory/command"
	"github.com/beabys/wms/inventory-service/internal/application/inventory/repository"
	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// EventPublisher is the port for publishing inventory domain events.
type EventPublisher interface {
	PublishStockAdjusted(ctx context.Context, entry *model.StockEntry, reason string) error
	PublishLowStockAlert(ctx context.Context, entry *model.StockEntry, threshold float64) error
}

// InventoryServiceHandler defines the application service contract for inventory operations.
type InventoryServiceHandler interface {
	CreateProduct(ctx context.Context, cmd command.CreateProductCommand) (*command.ProductResult, error)
	GetProduct(ctx context.Context, qry command.GetProductQuery) (*command.ProductResult, error)
	ListProducts(ctx context.Context, qry command.ListProductsQuery) (*command.ListProductsResult, error)
	UpdateProduct(ctx context.Context, cmd command.UpdateProductCommand) (*command.ProductResult, error)
	ArchiveProduct(ctx context.Context, cmd command.ArchiveProductCommand) (*command.ProductResult, error)

	AddStock(ctx context.Context, cmd command.AddStockCommand) (*command.StockResult, error)
	AdjustStock(ctx context.Context, cmd command.AdjustStockCommand) (*command.StockResult, error)
	GetStock(ctx context.Context, qry command.GetStockQuery) (*command.StockResult, error)
	ListStock(ctx context.Context, qry command.ListStockQuery) (*command.ListStockResult, error)
	ReserveStock(ctx context.Context, cmd command.ReserveStockCommand) (*command.StockResult, error)
	ReleaseStock(ctx context.Context, cmd command.ReleaseStockCommand) (*command.StockResult, error)

	CreateBinLocation(ctx context.Context, cmd command.CreateBinLocationCommand) (*command.BinLocationResult, error)
	ListBinLocations(ctx context.Context, qry command.ListBinLocationsQuery) (*command.ListBinLocationsResult, error)
}

// InventoryService implements InventoryServiceHandler.
type InventoryService struct {
	productRepo    repository.ProductRepository
	stockRepo      repository.StockRepository
	binLocationRepo repository.BinLocationRepository
	events         EventPublisher
}
