package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/beabys/wms/inventory-service-bff/internal/api/v1"
	inventoryv1 "github.com/beabys/wms/proto/gen/go/inventory/v1"
)

// mockGRPCClient implements the inventory gRPC client for testing.
type mockGRPCClient struct {
	createProductResp     *inventoryv1.Product
	createProductErr      error
	getProductResp        *inventoryv1.Product
	getProductErr         error
	listProductsResp      []*inventoryv1.Product
	listProductsErr       error
	updateProductResp     *inventoryv1.Product
	updateProductErr      error
	archiveProductResp    *inventoryv1.Product
	archiveProductErr     error
	addStockResp          *inventoryv1.StockEntry
	addStockErr           error
	adjustStockResp       *inventoryv1.StockEntry
	adjustStockErr        error
	getStockResp          *inventoryv1.StockEntry
	getStockErr           error
	listStockResp         []*inventoryv1.StockEntry
	listStockErr          error
	reserveStockResp      *inventoryv1.StockEntry
	reserveStockErr       error
	releaseStockResp      *inventoryv1.StockEntry
	releaseStockErr       error
	createBinLocationResp *inventoryv1.BinLocation
	createBinLocationErr  error
	listBinLocationsResp  []*inventoryv1.BinLocation
	listBinLocationsErr   error
}

func (m *mockGRPCClient) CreateProduct(_ context.Context, sku, name, description, category, unit string, weightKg, lowStockThreshold float64, token string) (*inventoryv1.Product, error) {
	return m.createProductResp, m.createProductErr
}

func (m *mockGRPCClient) GetProduct(_ context.Context, id string, token string) (*inventoryv1.Product, error) {
	return m.getProductResp, m.getProductErr
}

func (m *mockGRPCClient) ListProducts(_ context.Context, category, search string, pageSize int32, token string) ([]*inventoryv1.Product, error) {
	return m.listProductsResp, m.listProductsErr
}

func (m *mockGRPCClient) UpdateProduct(_ context.Context, id, sku, name, description, category, unit string, weightKg, lowStockThreshold float64, token string) (*inventoryv1.Product, error) {
	return m.updateProductResp, m.updateProductErr
}

func (m *mockGRPCClient) ArchiveProduct(_ context.Context, id string, token string) (*inventoryv1.Product, error) {
	return m.archiveProductResp, m.archiveProductErr
}

func (m *mockGRPCClient) AddStock(_ context.Context, productID, binLocationID, lotNumber, expiryDate string, quantity float64, token string) (*inventoryv1.StockEntry, error) {
	return m.addStockResp, m.addStockErr
}

func (m *mockGRPCClient) AdjustStock(_ context.Context, id string, delta float64, reason, notes string, token string) (*inventoryv1.StockEntry, error) {
	return m.adjustStockResp, m.adjustStockErr
}

func (m *mockGRPCClient) GetStock(_ context.Context, id string, token string) (*inventoryv1.StockEntry, error) {
	return m.getStockResp, m.getStockErr
}

func (m *mockGRPCClient) ListStock(_ context.Context, productID, binLocationID, status string, pageSize int32, token string) ([]*inventoryv1.StockEntry, error) {
	return m.listStockResp, m.listStockErr
}

func (m *mockGRPCClient) ReserveStock(_ context.Context, id string, quantity float64, token string) (*inventoryv1.StockEntry, error) {
	return m.reserveStockResp, m.reserveStockErr
}

func (m *mockGRPCClient) ReleaseStock(_ context.Context, id string, quantity float64, token string) (*inventoryv1.StockEntry, error) {
	return m.releaseStockResp, m.releaseStockErr
}

func (m *mockGRPCClient) CreateBinLocation(_ context.Context, warehouseZone, aisle, rack, shelf string, token string) (*inventoryv1.BinLocation, error) {
	return m.createBinLocationResp, m.createBinLocationErr
}

func (m *mockGRPCClient) ListBinLocations(_ context.Context, pageSize int32, token string) ([]*inventoryv1.BinLocation, error) {
	return m.listBinLocationsResp, m.listBinLocationsErr
}

func (m *mockGRPCClient) Close() error {
	return nil
}

// setupTest creates an HttpServer with a mock gRPC client and chi router.
func setupTest(mock GRPCInventoryClient) *chi.Mux {
	hs := NewHttpServer().
		SetLogger(zap.NewNop()).
		SetInventoryClient(mock)

	r := chi.NewRouter()
	v1.HandlerWithOptions(hs, v1.ChiServerOptions{
		BaseRouter: r,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	})
	return r
}

// ---- Products ----

func TestCreateProduct_Success(t *testing.T) {
	mock := &mockGRPCClient{
		createProductResp: &inventoryv1.Product{
			Id:   "prod-1",
			Sku:  "SKU-001",
			Name: "Test Product",
		},
	}
	router := setupTest(mock)

	body := `{"sku":"SKU-001","name":"Test Product"}`
	req := httptest.NewRequest("POST", "/v1/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
	require.NotNil(t, apiResp["data"])
}

func TestCreateProduct_InvalidBody(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/products", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateProduct_GRPCError(t *testing.T) {
	mock := &mockGRPCClient{
		createProductErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	body := `{"sku":"SKU-001","name":"Test Product"}`
	req := httptest.NewRequest("POST", "/v1/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestGetProduct_Success(t *testing.T) {
	mock := &mockGRPCClient{
		getProductResp: &inventoryv1.Product{
			Id:   "prod-1",
			Sku:  "SKU-001",
			Name: "Test Product",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/products/prod-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestGetProduct_NotFound(t *testing.T) {
	mock := &mockGRPCClient{
		getProductErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/products/nonexistent", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestListProducts_Success(t *testing.T) {
	mock := &mockGRPCClient{
		listProductsResp: []*inventoryv1.Product{
			{Id: "prod-1", Sku: "SKU-001", Name: "Product 1"},
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/products?page_size=10&category=electronics", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestUpdateProduct_Success(t *testing.T) {
	mock := &mockGRPCClient{
		updateProductResp: &inventoryv1.Product{
			Id:   "prod-1",
			Sku:  "SKU-001",
			Name: "Updated Product",
		},
	}
	router := setupTest(mock)

	body := `{"name":"Updated Product"}`
	req := httptest.NewRequest("PUT", "/v1/products/prod-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestUpdateProduct_InvalidBody(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("PUT", "/v1/products/prod-1", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestArchiveProduct_Success(t *testing.T) {
	mock := &mockGRPCClient{
		archiveProductResp: &inventoryv1.Product{
			Id:   "prod-1",
			Sku:  "SKU-001",
			Name: "Archived Product",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("DELETE", "/v1/products/prod-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

// ---- Stock ----

func TestAddStock_Success(t *testing.T) {
	mock := &mockGRPCClient{
		addStockResp: &inventoryv1.StockEntry{
			Id:        "stock-1",
			ProductId: "prod-1",
			Quantity:  100,
		},
	}
	router := setupTest(mock)

	body := `{"product_id":"prod-1","bin_location_id":"bin-1","quantity":100}`
	req := httptest.NewRequest("POST", "/v1/stock", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestAddStock_InvalidBody(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/stock", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetStock_Success(t *testing.T) {
	mock := &mockGRPCClient{
		getStockResp: &inventoryv1.StockEntry{
			Id:        "stock-1",
			ProductId: "prod-1",
			Quantity:  100,
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/stock/stock-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestAdjustStock_Success(t *testing.T) {
	mock := &mockGRPCClient{
		adjustStockResp: &inventoryv1.StockEntry{
			Id:        "stock-1",
			ProductId: "prod-1",
			Quantity:  90,
		},
	}
	router := setupTest(mock)

	body := `{"delta":-10,"reason":"damaged","notes":"Broken item"}`
	req := httptest.NewRequest("PUT", "/v1/stock/stock-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestListStock_Success(t *testing.T) {
	mock := &mockGRPCClient{
		listStockResp: []*inventoryv1.StockEntry{
			{Id: "stock-1", ProductId: "prod-1", Quantity: 100},
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/stock?product_id=prod-1&status=available", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}

func TestReserveStock_Success(t *testing.T) {
	mock := &mockGRPCClient{
		reserveStockResp: &inventoryv1.StockEntry{
			Id:               "stock-1",
			ProductId:        "prod-1",
			Quantity:         100,
			ReservedQuantity: 10,
		},
	}
	router := setupTest(mock)

	body := `{"quantity":10}`
	req := httptest.NewRequest("POST", "/v1/stock/stock-1/reserve", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}

func TestReleaseStock_Success(t *testing.T) {
	mock := &mockGRPCClient{
		releaseStockResp: &inventoryv1.StockEntry{
			Id:               "stock-1",
			ProductId:        "prod-1",
			Quantity:         100,
			ReservedQuantity: 0,
		},
	}
	router := setupTest(mock)

	body := `{"quantity":10}`
	req := httptest.NewRequest("POST", "/v1/stock/stock-1/release", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}

// ---- Bin Locations ----

func TestCreateBinLocation_Success(t *testing.T) {
	mock := &mockGRPCClient{
		createBinLocationResp: &inventoryv1.BinLocation{
			Id:            "bin-1",
			WarehouseZone: "WH-A",
			Aisle:         "A1",
			Rack:          "R1",
			Shelf:         "S1",
		},
	}
	router := setupTest(mock)

	body := `{"warehouse_zone":"WH-A","aisle":"A1","rack":"R1","shelf":"S1"}`
	req := httptest.NewRequest("POST", "/v1/bin-locations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestCreateBinLocation_InvalidBody(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/bin-locations", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestListBinLocations_Success(t *testing.T) {
	mock := &mockGRPCClient{
		listBinLocationsResp: []*inventoryv1.BinLocation{
			{Id: "bin-1", WarehouseZone: "WH-A", Aisle: "A1", Rack: "R1", Shelf: "S1"},
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/bin-locations?page_size=10", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}

func TestNotFound(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/nonexistent", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}
