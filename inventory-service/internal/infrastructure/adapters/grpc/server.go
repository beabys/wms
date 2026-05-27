package inventorygrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryv1 "github.com/beabys/wms/proto/gen/go/inventory/v1"

	"github.com/beabys/wms/inventory-service/internal/application/inventory/command"
	"github.com/beabys/wms/inventory-service/internal/application/inventory/usecase"
	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// InventoryServer implements the inventory.v1.InventoryService gRPC server.
type InventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	svc usecase.InventoryServiceHandler
}

// NewInventoryServer creates a new InventoryServer.
func NewInventoryServer(svc usecase.InventoryServiceHandler) *InventoryServer {
	return &InventoryServer{svc: svc}
}

// CreateProduct handles CreateProduct RPC.
func (s *InventoryServer) CreateProduct(ctx context.Context, req *inventoryv1.CreateProductRequest) (*inventoryv1.CreateProductResponse, error) {
	cmd := command.CreateProductCommand{
		SKU:               req.Sku,
		Name:              req.Name,
		Description:       req.Description,
		Category:          req.Category,
		Unit:              req.Unit,
		WeightKg:          req.WeightKg,
		LowStockThreshold: req.LowStockThreshold,
	}

	result, err := s.svc.CreateProduct(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create product: %v", err)
	}

	return &inventoryv1.CreateProductResponse{
		Product: productDomainToProto(result.Product),
	}, nil
}

// GetProduct handles GetProduct RPC.
func (s *InventoryServer) GetProduct(ctx context.Context, req *inventoryv1.GetProductRequest) (*inventoryv1.GetProductResponse, error) {
	qry := command.GetProductQuery{ProductID: req.Id}

	result, err := s.svc.GetProduct(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get product: %v", err)
	}

	return &inventoryv1.GetProductResponse{
		Product: productDomainToProto(result.Product),
	}, nil
}

// ListProducts handles ListProducts RPC.
func (s *InventoryServer) ListProducts(ctx context.Context, req *inventoryv1.ListProductsRequest) (*inventoryv1.ListProductsResponse, error) {
	qry := command.ListProductsQuery{
		Category: req.Category,
		Search:   req.Search,
	}
	if req.Pagination != nil {
		qry.PageSize = req.Pagination.Limit
		qry.PageToken = fmt.Sprintf("%d", req.Pagination.Page)
	}

	result, err := s.svc.ListProducts(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list products: %v", err)
	}

	protoProducts := make([]*inventoryv1.Product, len(result.Products))
	for i, p := range result.Products {
		protoProducts[i] = productDomainToProto(p)
	}

	return &inventoryv1.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// UpdateProduct handles UpdateProduct RPC.
func (s *InventoryServer) UpdateProduct(ctx context.Context, req *inventoryv1.UpdateProductRequest) (*inventoryv1.UpdateProductResponse, error) {
	cmd := command.UpdateProductCommand{
		ProductID:   req.Id,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Unit:        req.Unit,
		WeightKg:    req.WeightKg,
	}

	result, err := s.svc.UpdateProduct(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update product: %v", err)
	}

	return &inventoryv1.UpdateProductResponse{
		Product: productDomainToProto(result.Product),
	}, nil
}

// ArchiveProduct handles ArchiveProduct RPC.
func (s *InventoryServer) ArchiveProduct(ctx context.Context, req *inventoryv1.ArchiveProductRequest) (*inventoryv1.ArchiveProductResponse, error) {
	cmd := command.ArchiveProductCommand{ProductID: req.Id}

	result, err := s.svc.ArchiveProduct(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "archive product: %v", err)
	}

	return &inventoryv1.ArchiveProductResponse{
		Product: productDomainToProto(result.Product),
	}, nil
}

// AddStock handles AddStock RPC.
func (s *InventoryServer) AddStock(ctx context.Context, req *inventoryv1.AddStockRequest) (*inventoryv1.AddStockResponse, error) {
	var expiry *time.Time
	if req.ExpiryDate != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiryDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid expiry_date: %v", err)
		}
		expiry = &t
	}

	cmd := command.AddStockCommand{
		ProductID:     req.ProductId,
		BinLocationID: req.BinLocationId,
		Quantity:      req.Quantity,
		LotNumber:     req.LotNumber,
		ExpiryDate:    expiry,
	}

	result, err := s.svc.AddStock(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "add stock: %v", err)
	}

	return &inventoryv1.AddStockResponse{
		StockEntry: stockDomainToProto(result.StockEntry),
	}, nil
}

// AdjustStock handles AdjustStock RPC.
func (s *InventoryServer) AdjustStock(ctx context.Context, req *inventoryv1.AdjustStockRequest) (*inventoryv1.AdjustStockResponse, error) {
	cmd := command.AdjustStockCommand{
		StockEntryID: req.Id,
		Delta:        req.Delta,
		Reason:       req.Reason,
		Notes:        req.Notes,
	}

	result, err := s.svc.AdjustStock(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "adjust stock: %v", err)
	}

	return &inventoryv1.AdjustStockResponse{
		StockEntry: stockDomainToProto(result.StockEntry),
	}, nil
}

// GetStock handles GetStock RPC.
func (s *InventoryServer) GetStock(ctx context.Context, req *inventoryv1.GetStockRequest) (*inventoryv1.GetStockResponse, error) {
	qry := command.GetStockQuery{StockEntryID: req.Id}

	result, err := s.svc.GetStock(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get stock: %v", err)
	}

	return &inventoryv1.GetStockResponse{
		StockEntry: stockDomainToProto(result.StockEntry),
	}, nil
}

// ListStock handles ListStock RPC.
func (s *InventoryServer) ListStock(ctx context.Context, req *inventoryv1.ListStockRequest) (*inventoryv1.ListStockResponse, error) {
	qry := command.ListStockQuery{
		ProductID:     req.ProductId,
		BinLocationID: req.BinLocationId,
		Status:        req.Status,
	}
	if req.Pagination != nil {
		qry.PageSize = req.Pagination.Limit
		qry.PageToken = fmt.Sprintf("%d", req.Pagination.Page)
	}

	result, err := s.svc.ListStock(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list stock: %v", err)
	}

	protoEntries := make([]*inventoryv1.StockEntry, len(result.StockEntries))
	for i, e := range result.StockEntries {
		protoEntries[i] = stockDomainToProto(e)
	}

	return &inventoryv1.ListStockResponse{
		StockEntries: protoEntries,
	}, nil
}

// ReserveStock handles ReserveStock RPC.
func (s *InventoryServer) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	cmd := command.ReserveStockCommand{
		StockEntryID: req.Id,
		Quantity:     req.Quantity,
	}

	result, err := s.svc.ReserveStock(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "reserve stock: %v", err)
	}

	return &inventoryv1.ReserveStockResponse{
		StockEntry: stockDomainToProto(result.StockEntry),
	}, nil
}

// ReleaseStock handles ReleaseStock RPC.
func (s *InventoryServer) ReleaseStock(ctx context.Context, req *inventoryv1.ReleaseStockRequest) (*inventoryv1.ReleaseStockResponse, error) {
	cmd := command.ReleaseStockCommand{
		StockEntryID: req.Id,
		Quantity:     req.Quantity,
	}

	result, err := s.svc.ReleaseStock(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "release stock: %v", err)
	}

	return &inventoryv1.ReleaseStockResponse{
		StockEntry: stockDomainToProto(result.StockEntry),
	}, nil
}

// CreateBinLocation handles CreateBinLocation RPC.
func (s *InventoryServer) CreateBinLocation(ctx context.Context, req *inventoryv1.CreateBinLocationRequest) (*inventoryv1.CreateBinLocationResponse, error) {
	cmd := command.CreateBinLocationCommand{
		WarehouseZone: req.WarehouseZone,
		Aisle:         req.Aisle,
		Rack:          req.Rack,
		Shelf:         req.Shelf,
	}

	result, err := s.svc.CreateBinLocation(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create bin location: %v", err)
	}

	return &inventoryv1.CreateBinLocationResponse{
		BinLocation: binLocationDomainToProto(result.BinLocation),
	}, nil
}

// ListBinLocations handles ListBinLocations RPC.
func (s *InventoryServer) ListBinLocations(ctx context.Context, req *inventoryv1.ListBinLocationsRequest) (*inventoryv1.ListBinLocationsResponse, error) {
	qry := command.ListBinLocationsQuery{}
	if req.Pagination != nil {
		qry.PageSize = req.Pagination.Limit
		qry.PageToken = fmt.Sprintf("%d", req.Pagination.Page)
	}

	result, err := s.svc.ListBinLocations(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list bin locations: %v", err)
	}

	protoLocations := make([]*inventoryv1.BinLocation, len(result.BinLocations))
	for i, l := range result.BinLocations {
		protoLocations[i] = binLocationDomainToProto(l)
	}

	return &inventoryv1.ListBinLocationsResponse{
		BinLocations: protoLocations,
	}, nil
}

// --- mapping helpers ---

func productDomainToProto(p *model.Product) *inventoryv1.Product {
	if p == nil {
		return nil
	}

	return &inventoryv1.Product{
		Id:                p.ID,
		Sku:               p.SKU,
		Name:              p.Name,
		Description:       p.Description,
		Category:          p.Category,
		Unit:              p.Unit,
		WeightKg:          p.WeightKg,
		IsActive:          p.IsActive,
		LowStockThreshold: p.LowStockThreshold,
		CreatedAt:         p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         p.UpdatedAt.Format(time.RFC3339),
	}
}

func stockDomainToProto(s *model.StockEntry) *inventoryv1.StockEntry {
	if s == nil {
		return nil
	}

	pb := &inventoryv1.StockEntry{
		Id:               s.ID,
		ProductId:        s.ProductID,
		BinLocationId:    s.BinLocationID,
		Quantity:         s.Quantity,
		ReservedQuantity: s.ReservedQuantity,
		LotNumber:        s.LotNumber,
		Status:           s.Status,
		LastUpdated:      s.LastUpdated.Format(time.RFC3339),
	}

	if s.ExpiryDate != nil {
		pb.ExpiryDate = s.ExpiryDate.Format(time.RFC3339)
	}

	return pb
}

func binLocationDomainToProto(b *model.BinLocation) *inventoryv1.BinLocation {
	if b == nil {
		return nil
	}

	return &inventoryv1.BinLocation{
		Id:            b.ID,
		WarehouseZone: b.WarehouseZone,
		Aisle:         b.Aisle,
		Rack:          b.Rack,
		Shelf:         b.Shelf,
		IsActive:      b.IsActive,
	}
}
