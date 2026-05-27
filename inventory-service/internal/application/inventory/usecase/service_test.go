package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/inventory-service/internal/application/inventory/command"
	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// --- mocks ---

type mockProductRepo struct {
	mock.Mock
}

func (m *mockProductRepo) GetByID(ctx context.Context, id string) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *mockProductRepo) Save(ctx context.Context, p *model.Product) error {
	return m.Called(ctx, p).Error(0)
}

func (m *mockProductRepo) Update(ctx context.Context, p *model.Product) error {
	return m.Called(ctx, p).Error(0)
}

func (m *mockProductRepo) List(ctx context.Context, category, search string, pageSize int32, pageToken string) ([]*model.Product, string, error) {
	args := m.Called(ctx, category, search, pageSize, pageToken)
	return args.Get(0).([]*model.Product), args.String(1), args.Error(2)
}

func (m *mockProductRepo) Archive(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockStockRepo struct {
	mock.Mock
}

func (m *mockStockRepo) GetByID(ctx context.Context, id string) (*model.StockEntry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StockEntry), args.Error(1)
}

func (m *mockStockRepo) Save(ctx context.Context, e *model.StockEntry) error {
	return m.Called(ctx, e).Error(0)
}

func (m *mockStockRepo) UpdateQuantity(ctx context.Context, e *model.StockEntry) error {
	return m.Called(ctx, e).Error(0)
}

func (m *mockStockRepo) List(ctx context.Context, productID, binLocationID, status string, pageSize int32, pageToken string) ([]*model.StockEntry, string, error) {
	args := m.Called(ctx, productID, binLocationID, status, pageSize, pageToken)
	return args.Get(0).([]*model.StockEntry), args.String(1), args.Error(2)
}

type mockBinLocationRepo struct {
	mock.Mock
}

func (m *mockBinLocationRepo) GetByID(ctx context.Context, id string) (*model.BinLocation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BinLocation), args.Error(1)
}

func (m *mockBinLocationRepo) Save(ctx context.Context, l *model.BinLocation) error {
	return m.Called(ctx, l).Error(0)
}

func (m *mockBinLocationRepo) List(ctx context.Context, pageSize int32, pageToken string) ([]*model.BinLocation, string, error) {
	args := m.Called(ctx, pageSize, pageToken)
	return args.Get(0).([]*model.BinLocation), args.String(1), args.Error(2)
}

type mockEventPublisher struct {
	mock.Mock
}

func (m *mockEventPublisher) PublishStockAdjusted(ctx context.Context, entry *model.StockEntry, reason string) error {
	return m.Called(ctx, entry, reason).Error(0)
}

func (m *mockEventPublisher) PublishLowStockAlert(ctx context.Context, entry *model.StockEntry, threshold float64) error {
	return m.Called(ctx, entry, threshold).Error(0)
}

// --- tests ---

func TestCreateProduct(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockStock := new(mockStockRepo)
	mockBin := new(mockBinLocationRepo)
	mockEvt := new(mockEventPublisher)

	svc := NewInventoryService(mockProd, mockStock, mockBin, mockEvt)

	cmd := command.CreateProductCommand{
		SKU:               "SKU-001",
		Name:              "Widget",
		Description:       "A widget",
		Category:          "gadgets",
		Unit:              "pcs",
		WeightKg:          1.5,
		LowStockThreshold: 10,
	}

	mockProd.On("Save", mock.Anything, mock.AnythingOfType("*model.Product")).Return(nil)

	result, err := svc.CreateProduct(context.Background(), cmd)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Widget", result.Product.Name)
	assert.Equal(t, "SKU-001", result.Product.SKU)
	assert.True(t, result.Product.IsActive)
	mockProd.AssertExpectations(t)
}

func TestGetProduct(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockStock := new(mockStockRepo)
	mockBin := new(mockBinLocationRepo)
	mockEvt := new(mockEventPublisher)

	svc := NewInventoryService(mockProd, mockStock, mockBin, mockEvt)

	expected := &model.Product{ID: "prod_1", SKU: "SKU-001", Name: "Widget", IsActive: true}
	mockProd.On("GetByID", mock.Anything, "prod_1").Return(expected, nil)

	result, err := svc.GetProduct(context.Background(), command.GetProductQuery{ProductID: "prod_1"})
	require.NoError(t, err)
	assert.Equal(t, expected, result.Product)
	mockProd.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockStock := new(mockStockRepo)
	mockBin := new(mockBinLocationRepo)
	mockEvt := new(mockEventPublisher)

	svc := NewInventoryService(mockProd, mockStock, mockBin, mockEvt)

	existing := &model.Product{ID: "prod_1", SKU: "SKU-001", Name: "Old", IsActive: true}
	mockProd.On("GetByID", mock.Anything, "prod_1").Return(existing, nil)
	mockProd.On("Update", mock.Anything, mock.AnythingOfType("*model.Product")).Return(nil)

	cmd := command.UpdateProductCommand{
		ProductID: "prod_1",
		Name:      "NewName",
		Category:  "new-cat",
		Unit:      "kg",
		WeightKg:  2.0,
	}
	result, err := svc.UpdateProduct(context.Background(), cmd)
	require.NoError(t, err)
	assert.Equal(t, "NewName", result.Product.Name)
	mockProd.AssertExpectations(t)
}

func TestAddStock(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockStock := new(mockStockRepo)
	mockBin := new(mockBinLocationRepo)
	mockEvt := new(mockEventPublisher)

	svc := NewInventoryService(mockProd, mockStock, mockBin, mockEvt)

	mockStock.On("Save", mock.Anything, mock.AnythingOfType("*model.StockEntry")).Return(nil)

	cmd := command.AddStockCommand{
		ProductID: "prod_1",
		Quantity:  100,
		LotNumber: "LOT-001",
	}
	result, err := svc.AddStock(context.Background(), cmd)
	require.NoError(t, err)
	assert.Equal(t, 100.0, result.StockEntry.Quantity)
	assert.Equal(t, "active", result.StockEntry.Status)
	mockStock.AssertExpectations(t)
}

func TestAdjustStock(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockStock := new(mockStockRepo)
	mockBin := new(mockBinLocationRepo)
	mockEvt := new(mockEventPublisher)

	svc := NewInventoryService(mockProd, mockStock, mockBin, mockEvt)

	entry := &model.StockEntry{ID: "stk_1", ProductID: "prod_1", Quantity: 100}
	mockStock.On("GetByID", mock.Anything, "stk_1").Return(entry, nil)
	mockStock.On("UpdateQuantity", mock.Anything, mock.AnythingOfType("*model.StockEntry")).Return(nil)
	mockEvt.On("PublishStockAdjusted", mock.Anything, mock.AnythingOfType("*model.StockEntry"), "MANUAL").Return(nil)

	cmd := command.AdjustStockCommand{
		StockEntryID: "stk_1",
		Delta:        -20,
		Reason:       "MANUAL",
		Notes:        "damaged",
	}
	result, err := svc.AdjustStock(context.Background(), cmd)
	require.NoError(t, err)
	assert.Equal(t, 80.0, result.StockEntry.Quantity)
	mockStock.AssertExpectations(t)
	mockEvt.AssertExpectations(t)
}

func TestCreateBinLocation(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockStock := new(mockStockRepo)
	mockBin := new(mockBinLocationRepo)
	mockEvt := new(mockEventPublisher)

	svc := NewInventoryService(mockProd, mockStock, mockBin, mockEvt)

	mockBin.On("Save", mock.Anything, mock.AnythingOfType("*model.BinLocation")).Return(nil)

	cmd := command.CreateBinLocationCommand{
		WarehouseZone: "ZONE-A",
		Aisle:         "A1",
		Rack:          "R1",
		Shelf:         "S1",
	}
	result, err := svc.CreateBinLocation(context.Background(), cmd)
	require.NoError(t, err)
	assert.Equal(t, "ZONE-A", result.BinLocation.WarehouseZone)
	assert.True(t, result.BinLocation.IsActive)
	mockBin.AssertExpectations(t)
}

func TestReserveStock(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockStock := new(mockStockRepo)
	mockBin := new(mockBinLocationRepo)
	mockEvt := new(mockEventPublisher)

	svc := NewInventoryService(mockProd, mockStock, mockBin, mockEvt)

	entry := &model.StockEntry{ID: "stk_1", ProductID: "prod_1", Quantity: 100}
	mockStock.On("GetByID", mock.Anything, "stk_1").Return(entry, nil)
	mockStock.On("UpdateQuantity", mock.Anything, mock.AnythingOfType("*model.StockEntry")).Return(nil)

	result, err := svc.ReserveStock(context.Background(), command.ReserveStockCommand{StockEntryID: "stk_1", Quantity: 30})
	require.NoError(t, err)
	assert.Equal(t, 30.0, result.StockEntry.ReservedQuantity)
	mockStock.AssertExpectations(t)
}

func TestReleaseStock(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockStock := new(mockStockRepo)
	mockBin := new(mockBinLocationRepo)
	mockEvt := new(mockEventPublisher)

	svc := NewInventoryService(mockProd, mockStock, mockBin, mockEvt)

	entry := &model.StockEntry{ID: "stk_1", ProductID: "prod_1", Quantity: 100, ReservedQuantity: 50}
	mockStock.On("GetByID", mock.Anything, "stk_1").Return(entry, nil)
	mockStock.On("UpdateQuantity", mock.Anything, mock.AnythingOfType("*model.StockEntry")).Return(nil)

	result, err := svc.ReleaseStock(context.Background(), command.ReleaseStockCommand{StockEntryID: "stk_1", Quantity: 20})
	require.NoError(t, err)
	assert.Equal(t, 30.0, result.StockEntry.ReservedQuantity)
	mockStock.AssertExpectations(t)
}
