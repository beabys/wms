package http

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	inventoryv1 "github.com/beabys/wms/proto/gen/go/inventory/v1"
)

// GRPCInventoryClient defines the interface the HTTP handler needs from the gRPC layer.
type GRPCInventoryClient interface {
	CreateProduct(ctx context.Context, sku, name, description, category, unit string, weightKg, lowStockThreshold float64, token string) (*inventoryv1.Product, error)
	GetProduct(ctx context.Context, id string, token string) (*inventoryv1.Product, error)
	ListProducts(ctx context.Context, category, search string, pageSize int32, token string) ([]*inventoryv1.Product, error)
	UpdateProduct(ctx context.Context, id, sku, name, description, category, unit string, weightKg, lowStockThreshold float64, token string) (*inventoryv1.Product, error)
	ArchiveProduct(ctx context.Context, id string, token string) (*inventoryv1.Product, error)
	AddStock(ctx context.Context, productID, binLocationID, lotNumber, expiryDate string, quantity float64, token string) (*inventoryv1.StockEntry, error)
	AdjustStock(ctx context.Context, id string, delta float64, reason, notes string, token string) (*inventoryv1.StockEntry, error)
	GetStock(ctx context.Context, id string, token string) (*inventoryv1.StockEntry, error)
	ListStock(ctx context.Context, productID, binLocationID, status string, pageSize int32, token string) ([]*inventoryv1.StockEntry, error)
	ReserveStock(ctx context.Context, id string, quantity float64, token string) (*inventoryv1.StockEntry, error)
	ReleaseStock(ctx context.Context, id string, quantity float64, token string) (*inventoryv1.StockEntry, error)
	CreateBinLocation(ctx context.Context, warehouseZone, aisle, rack, shelf string, token string) (*inventoryv1.BinLocation, error)
	ListBinLocations(ctx context.Context, pageSize int32, token string) ([]*inventoryv1.BinLocation, error)
	Close() error
}

// HttpServer holds dependencies for the HTTP server.
type HttpServer struct {
	Server           *http.Server
	Config           *Config
	Logger           *zap.Logger
	InventoryClient GRPCInventoryClient
}

// Config holds BFF HTTP configuration.
type Config struct {
	Port int `yaml:"port"`
}

// NewHttpServer creates a new HttpServer.
func NewHttpServer() *HttpServer {
	return &HttpServer{}
}

// SetConfig sets the server config.
func (hs *HttpServer) SetConfig(cfg *Config) *HttpServer {
	hs.Config = cfg
	return hs
}

// SetLogger sets the logger.
func (hs *HttpServer) SetLogger(l *zap.Logger) *HttpServer {
	hs.Logger = l
	return hs
}

// SetInventoryClient sets the gRPC client for inventory operations.
func (hs *HttpServer) SetInventoryClient(client GRPCInventoryClient) *HttpServer {
	hs.InventoryClient = client
	return hs
}

// Run starts the HTTP server and handles graceful shutdown.
func (hs *HttpServer) Run(ctx context.Context, wg *errgroup.Group) {
	wg.Go(func() error {
		hs.Logger.Info("http server started", zap.Int("port", hs.Config.Port))
		if err := hs.Server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return nil
			}
			hs.Logger.Error("http server stopped with error", zap.Error(err))
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		hs.Logger.Info("shutting down gracefully http server")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		if err := hs.Server.Shutdown(ctxTimeout); err != nil {
			hs.Logger.Error("error shutting server down", zap.Error(err))
			return err
		}
		return nil
	})
}
