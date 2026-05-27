package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"

	v1 "github.com/beabys/wms/inventory-service-bff/internal/api/v1"
)

// extractToken extracts the Bearer token from the Authorization header.
func extractToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// ListProducts handles GET /v1/products.
func (hs *HttpServer) ListProducts(w http.ResponseWriter, r *http.Request, params v1.ListProductsParams) {
	pageSize := int32(0)
	if params.PageSize != nil {
		pageSize = int32(*params.PageSize)
	}

	category := ""
	if params.Category != nil {
		category = *params.Category
	}

	search := ""
	if params.Search != nil {
		search = *params.Search
	}

	token := extractToken(r)
	products, err := hs.InventoryClient.ListProducts(r.Context(), category, search, pageSize, token)
	if err != nil {
		hs.Logger.Error("list products", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, products)
}

// CreateProduct handles POST /v1/products.
func (hs *HttpServer) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateProductJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode create product body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	weightKg := 0.0
	if req.WeightKg != nil {
		weightKg = *req.WeightKg
	}

	lowStockThreshold := 0.0
	if req.LowStockThreshold != nil {
		lowStockThreshold = *req.LowStockThreshold
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	category := ""
	if req.Category != nil {
		category = *req.Category
	}

	unit := ""
	if req.Unit != nil {
		unit = *req.Unit
	}

	token := extractToken(r)
	product, err := hs.InventoryClient.CreateProduct(r.Context(), req.SKU, req.Name, description, category, unit, weightKg, lowStockThreshold, token)
	if err != nil {
		hs.Logger.Error("create product", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, product)
}

// GetProduct handles GET /v1/products/{id}.
func (hs *HttpServer) GetProduct(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	product, err := hs.InventoryClient.GetProduct(r.Context(), id, token)
	if err != nil {
		hs.Logger.Error("get product", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusNotFound, err)
		return
	}

	successResponseJSON(w, product)
}

// UpdateProduct handles PUT /v1/products/{id}.
func (hs *HttpServer) UpdateProduct(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.UpdateProductJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode update product body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	sku := ""
	if req.SKU != nil {
		sku = *req.SKU
	}
	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	category := ""
	if req.Category != nil {
		category = *req.Category
	}
	unit := ""
	if req.Unit != nil {
		unit = *req.Unit
	}
	weightKg := 0.0
	if req.WeightKg != nil {
		weightKg = *req.WeightKg
	}
	lowStockThreshold := 0.0
	if req.LowStockThreshold != nil {
		lowStockThreshold = *req.LowStockThreshold
	}

	token := extractToken(r)
	product, err := hs.InventoryClient.UpdateProduct(r.Context(), id, sku, name, description, category, unit, weightKg, lowStockThreshold, token)
	if err != nil {
		hs.Logger.Error("update product", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, product)
}

// ArchiveProduct handles DELETE /v1/products/{id}.
func (hs *HttpServer) ArchiveProduct(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	product, err := hs.InventoryClient.ArchiveProduct(r.Context(), id, token)
	if err != nil {
		hs.Logger.Error("archive product", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, product)
}

// ListStock handles GET /v1/stock.
func (hs *HttpServer) ListStock(w http.ResponseWriter, r *http.Request, params v1.ListStockParams) {
	pageSize := int32(0)
	if params.PageSize != nil {
		pageSize = int32(*params.PageSize)
	}

	productID := ""
	if params.ProductId != nil {
		productID = *params.ProductId
	}

	binLocationID := ""
	if params.BinLocationId != nil {
		binLocationID = *params.BinLocationId
	}

	status := ""
	if params.Status != nil {
		status = *params.Status
	}

	token := extractToken(r)
	stockEntries, err := hs.InventoryClient.ListStock(r.Context(), productID, binLocationID, status, pageSize, token)
	if err != nil {
		hs.Logger.Error("list stock", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, stockEntries)
}

// AddStock handles POST /v1/stock.
func (hs *HttpServer) AddStock(w http.ResponseWriter, r *http.Request) {
	var req v1.AddStockJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode add stock body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	lotNumber := ""
	if req.LotNumber != nil {
		lotNumber = *req.LotNumber
	}

	expiryDate := ""
	if req.ExpiryDate != nil {
		expiryDate = *req.ExpiryDate
	}

	token := extractToken(r)
	stockEntry, err := hs.InventoryClient.AddStock(r.Context(), req.ProductId, req.BinLocationId, lotNumber, expiryDate, req.Quantity, token)
	if err != nil {
		hs.Logger.Error("add stock", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, stockEntry)
}

// GetStock handles GET /v1/stock/{id}.
func (hs *HttpServer) GetStock(w http.ResponseWriter, r *http.Request, id string) {
	token := extractToken(r)
	stockEntry, err := hs.InventoryClient.GetStock(r.Context(), id, token)
	if err != nil {
		hs.Logger.Error("get stock", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusNotFound, err)
		return
	}

	successResponseJSON(w, stockEntry)
}

// AdjustStock handles PUT /v1/stock/{id}.
func (hs *HttpServer) AdjustStock(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.AdjustStockJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode adjust stock body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	token := extractToken(r)
	stockEntry, err := hs.InventoryClient.AdjustStock(r.Context(), id, req.Delta, reason, notes, token)
	if err != nil {
		hs.Logger.Error("adjust stock", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, stockEntry)
}

// ReserveStock handles POST /v1/stock/{id}/reserve.
func (hs *HttpServer) ReserveStock(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.ReserveStockJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode reserve stock body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	token := extractToken(r)
	stockEntry, err := hs.InventoryClient.ReserveStock(r.Context(), id, req.Quantity, token)
	if err != nil {
		hs.Logger.Error("reserve stock", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, stockEntry)
}

// ReleaseStock handles POST /v1/stock/{id}/release.
func (hs *HttpServer) ReleaseStock(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.ReleaseStockJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode release stock body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	token := extractToken(r)
	stockEntry, err := hs.InventoryClient.ReleaseStock(r.Context(), id, req.Quantity, token)
	if err != nil {
		hs.Logger.Error("release stock", zap.Error(err), zap.String("id", id))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, stockEntry)
}

// ListBinLocations handles GET /v1/bin-locations.
func (hs *HttpServer) ListBinLocations(w http.ResponseWriter, r *http.Request, params v1.ListBinLocationsParams) {
	pageSize := int32(0)
	if params.PageSize != nil {
		pageSize = int32(*params.PageSize)
	}

	token := extractToken(r)
	binLocations, err := hs.InventoryClient.ListBinLocations(r.Context(), pageSize, token)
	if err != nil {
		hs.Logger.Error("list bin locations", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, binLocations)
}

// CreateBinLocation handles POST /v1/bin-locations.
func (hs *HttpServer) CreateBinLocation(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateBinLocationJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		hs.Logger.Error("decode create bin location body", zap.Error(err))
		errorResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	token := extractToken(r)
	binLocation, err := hs.InventoryClient.CreateBinLocation(r.Context(), req.WarehouseZone, req.Aisle, req.Rack, req.Shelf, token)
	if err != nil {
		hs.Logger.Error("create bin location", zap.Error(err))
		errorResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	successResponseJSON(w, binLocation)
}
