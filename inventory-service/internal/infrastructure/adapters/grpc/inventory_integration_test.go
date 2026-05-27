// Copyright (c) 2026 WMS.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

//go:build integration

package inventorygrpc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	inventoryv1 "github.com/beabys/wms/proto/gen/go/inventory/v1"
)

const (
	addr = "localhost:50004"
)

// --- helpers ---

func testToken(t *testing.T) string {
	t.Helper()

	key, err := loadPrivateKey()
	if err != nil {
		t.Logf("warning: could not load private key, generating temp key: %v", err)
		key, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("generate temp key: %v", err)
		}
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": "e2e-test-user",
		"iat": now.Unix(),
		"exp": now.Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}
	return signed
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	candidates := []string{
		"login-service/keys/private.pem",
		"../login-service/keys/private.pem",
		"../../login-service/keys/private.pem",
		"../../../login-service/keys/private.pem",
		"../../../../login-service/keys/private.pem",
		"../../../../../login-service/keys/private.pem",
	}
	var firstErr error
	for _, path := range candidates {
		key, err := readPrivateKey(path)
		if err == nil {
			return key, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	return nil, fmt.Errorf("no private key found (tried %v): %w", candidates, firstErr)
}

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}
	return rsaKey, nil
}

func dial(t *testing.T) *grpc.ClientConn {
	t.Helper()
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial %s: %v", addr, err)
	}
	return conn
}

func authContext(t *testing.T) context.Context {
	t.Helper()
	return metadata.AppendToOutgoingContext(
		context.Background(),
		"authorization", "Bearer "+testToken(t),
	)
}

// --- Product Tests ---

// Scenario 1+2: Create product + Get product by ID
func TestInventoryIntegration_CreateAndGetProduct(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	var productID string
	sku := fmt.Sprintf("E2E-CREATE-%d", time.Now().UnixNano())

	t.Run("CreateProductSuccess", func(t *testing.T) {
		req := &inventoryv1.CreateProductRequest{
			Sku:               sku,
			Name:              "E2E Widget",
			Description:       "E2E test product",
			Category:          "e2e-test",
			Unit:              "pcs",
			WeightKg:          1.5,
			LowStockThreshold: 10,
		}
		resp, err := client.CreateProduct(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateProduct: %v", err)
		}
		if resp.Product == nil {
			t.Fatal("CreateProduct: response product is nil")
		}
		if resp.Product.Id == "" {
			t.Fatal("CreateProduct: product id is empty")
		}
		if resp.Product.Sku != sku {
			t.Errorf("CreateProduct: sku = %q, want %q", resp.Product.Sku, sku)
		}
		if resp.Product.Name != "E2E Widget" {
			t.Errorf("CreateProduct: name = %q, want %q", resp.Product.Name, "E2E Widget")
		}
		if !resp.Product.IsActive {
			t.Error("CreateProduct: product should be active")
		}
		productID = resp.Product.Id
		t.Logf("created product: id=%s sku=%s", productID, sku)
	})

	if productID == "" {
		t.Fatal("productID not set, aborting")
	}

	t.Run("GetProductSuccess", func(t *testing.T) {
		resp, err := client.GetProduct(authContext(t), &inventoryv1.GetProductRequest{Id: productID})
		if err != nil {
			t.Fatalf("GetProduct: %v", err)
		}
		if resp.Product == nil {
			t.Fatal("GetProduct: response product is nil")
		}
		if resp.Product.Id != productID {
			t.Errorf("GetProduct: id = %q, want %q", resp.Product.Id, productID)
		}
		if resp.Product.Sku != sku {
			t.Errorf("GetProduct: sku = %q, want %q", resp.Product.Sku, sku)
		}
		if resp.Product.Name != "E2E Widget" {
			t.Errorf("GetProduct: name = %q, want %q", resp.Product.Name, "E2E Widget")
		}
	})
}

// Scenario 3: List products with filters
func TestInventoryIntegration_ListProducts(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	sku := fmt.Sprintf("E2E-LIST-%d", time.Now().UnixNano())
	t.Run("CreateProductForList", func(t *testing.T) {
		req := &inventoryv1.CreateProductRequest{
			Sku:               sku,
			Name:              "Listable Product",
			Description:       "For list test",
			Category:          "e2e-list-test",
			Unit:              "pcs",
			WeightKg:          2.0,
			LowStockThreshold: 5,
		}
		resp, err := client.CreateProduct(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateProduct: %v", err)
		}
		if resp.Product == nil || resp.Product.Id == "" {
			t.Fatal("CreateProduct: no id returned")
		}
		t.Logf("created product: id=%s sku=%s", resp.Product.Id, sku)
	})

	t.Run("ListProductsByCategory", func(t *testing.T) {
		req := &inventoryv1.ListProductsRequest{
			Category: "e2e-list-test",
			Pagination: &commonv1.Pagination{
				Limit: 100,
				Page:  1,
			},
		}
		resp, err := client.ListProducts(authContext(t), req)
		if err != nil {
			t.Fatalf("ListProducts: %v", err)
		}
		if len(resp.Products) == 0 {
			t.Fatal("ListProducts: empty results, expected at least 1 product")
		}
		found := false
		for _, p := range resp.Products {
			if p.Sku == sku {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListProducts: product with sku %s not found in results", sku)
		}
	})

	t.Run("ListProductsBySearch", func(t *testing.T) {
		req := &inventoryv1.ListProductsRequest{
			Search: "Listable",
			Pagination: &commonv1.Pagination{
				Limit: 100,
				Page:  1,
			},
		}
		resp, err := client.ListProducts(authContext(t), req)
		if err != nil {
			t.Fatalf("ListProducts by search: %v", err)
		}
		if len(resp.Products) == 0 {
			t.Fatal("ListProducts by search: empty results")
		}
	})
}

// Scenario 4: Update product
func TestInventoryIntegration_UpdateProduct(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	var productID string

	t.Run("CreateProductForUpdate", func(t *testing.T) {
		req := &inventoryv1.CreateProductRequest{
			Sku:  fmt.Sprintf("E2E-UPDATE-SKU-%d", time.Now().UnixNano()),
			Name: "Before Update",
			Unit: "pcs",
		}
		resp, err := client.CreateProduct(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateProduct: %v", err)
		}
		if resp.Product == nil || resp.Product.Id == "" {
			t.Fatal("CreateProduct: no id returned")
		}
		productID = resp.Product.Id
		t.Logf("created product: id=%s", productID)
	})

	if productID == "" {
		t.Fatal("productID not set, aborting")
	}

	t.Run("UpdateProductFields", func(t *testing.T) {
		req := &inventoryv1.UpdateProductRequest{
			Id:          productID,
			Name:        "After Update",
			Description: "Updated description",
			Category:    "updated-cat",
			Unit:        "kg",
			WeightKg:    3.0,
		}
		resp, err := client.UpdateProduct(authContext(t), req)
		if err != nil {
			t.Fatalf("UpdateProduct: %v", err)
		}
		if resp.Product == nil {
			t.Fatal("UpdateProduct: response product is nil")
		}
		if resp.Product.Name != "After Update" {
			t.Errorf("UpdateProduct: name = %q, want %q", resp.Product.Name, "After Update")
		}
		if resp.Product.Description != "Updated description" {
			t.Errorf("UpdateProduct: description = %q, want %q", resp.Product.Description, "Updated description")
		}
		if resp.Product.Category != "updated-cat" {
			t.Errorf("UpdateProduct: category = %q, want %q", resp.Product.Category, "updated-cat")
		}
		if resp.Product.Unit != "kg" {
			t.Errorf("UpdateProduct: unit = %q, want %q", resp.Product.Unit, "kg")
		}
		if resp.Product.WeightKg != 3.0 {
			t.Errorf("UpdateProduct: weight_kg = %f, want %f", resp.Product.WeightKg, 3.0)
		}
	})

	t.Run("VerifyUpdatedProduct", func(t *testing.T) {
		resp, err := client.GetProduct(authContext(t), &inventoryv1.GetProductRequest{Id: productID})
		if err != nil {
			t.Fatalf("GetProduct after update: %v", err)
		}
		if resp.Product.Name != "After Update" {
			t.Errorf("GetProduct after update: name = %q, want %q", resp.Product.Name, "After Update")
		}
	})
}

// Scenario 5: Archive product (soft-delete)
func TestInventoryIntegration_ArchiveProduct(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	var productID string

	t.Run("CreateProductForArchive", func(t *testing.T) {
		req := &inventoryv1.CreateProductRequest{
			Sku:  fmt.Sprintf("E2E-ARCHIVE-SKU-%d", time.Now().UnixNano()),
			Name: "To Be Archived",
			Unit: "pcs",
		}
		resp, err := client.CreateProduct(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateProduct: %v", err)
		}
		if resp.Product == nil || resp.Product.Id == "" {
			t.Fatal("CreateProduct: no id returned")
		}
		if !resp.Product.IsActive {
			t.Error("CreateProduct: new product should be active")
		}
		productID = resp.Product.Id
		t.Logf("created product: id=%s", productID)
	})

	if productID == "" {
		t.Fatal("productID not set, aborting")
	}

	t.Run("ArchiveProduct", func(t *testing.T) {
		resp, err := client.ArchiveProduct(authContext(t), &inventoryv1.ArchiveProductRequest{Id: productID})
		if err != nil {
			t.Fatalf("ArchiveProduct: %v", err)
		}
		if resp.Product == nil {
			t.Fatal("ArchiveProduct: response product is nil")
		}
		if resp.Product.IsActive {
			t.Error("ArchiveProduct: product should be inactive after archive")
		}
	})
}

// Scenario 6: Create with missing SKU → error
func TestInventoryIntegration_CreateProductValidation(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	t.Run("MissingSKU", func(t *testing.T) {
		req := &inventoryv1.CreateProductRequest{
			Sku:  "",
			Name: "No SKU",
			Unit: "pcs",
		}
		_, err := client.CreateProduct(authContext(t), req)
		if err == nil {
			t.Fatal("expected error for missing SKU, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		// The server returns codes.Internal wrapping a domain error
		t.Logf("MissingSKU: code=%v msg=%v", st.Code(), st.Message())
	})
}

// --- Bin Location Tests ---

// Scenario 15+16: Create and list bin locations
func TestInventoryIntegration_BinLocations(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	var binID string
	zone := fmt.Sprintf("ZONE-E2E-%d", time.Now().UnixNano())

	t.Run("CreateBinLocation", func(t *testing.T) {
		req := &inventoryv1.CreateBinLocationRequest{
			WarehouseZone: zone,
			Aisle:         "A1",
			Rack:          "R1",
			Shelf:         "S1",
		}
		resp, err := client.CreateBinLocation(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateBinLocation: %v", err)
		}
		if resp.BinLocation == nil {
			t.Fatal("CreateBinLocation: response bin_location is nil")
		}
		if resp.BinLocation.Id == "" {
			t.Fatal("CreateBinLocation: bin location id is empty")
		}
		if resp.BinLocation.WarehouseZone != zone {
			t.Errorf("CreateBinLocation: warehouse_zone = %q, want %q", resp.BinLocation.WarehouseZone, zone)
		}
		if !resp.BinLocation.IsActive {
			t.Error("CreateBinLocation: bin location should be active")
		}
		binID = resp.BinLocation.Id
		t.Logf("created bin location: id=%s zone=%s", binID, zone)
	})

	t.Run("ListBinLocations", func(t *testing.T) {
		req := &inventoryv1.ListBinLocationsRequest{
			Pagination: &commonv1.Pagination{
				Limit: 100,
				Page:  1,
			},
		}
		resp, err := client.ListBinLocations(authContext(t), req)
		if err != nil {
			t.Fatalf("ListBinLocations: %v", err)
		}
		if len(resp.BinLocations) == 0 {
			t.Fatal("ListBinLocations: empty results, expected at least 1 bin location")
		}
		if binID != "" {
			found := false
			for _, b := range resp.BinLocations {
				if b.Id == binID {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("ListBinLocations: bin location %s not found in results", binID)
			}
		}
	})
}

// --- Stock Entry Tests ---

// Full stock flow: Create product + bin location → Add stock → Get stock → Adjust → Reserve → Release
func TestInventoryIntegration_FullStockFlow(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	var productID string
	var binID string
	var stockID string

	// Setup: create product
	t.Run("SetupProduct", func(t *testing.T) {
		req := &inventoryv1.CreateProductRequest{
			Sku:  fmt.Sprintf("E2E-STOCK-FLOW-SKU-%d", time.Now().UnixNano()),
			Name: "Stock Flow Product",
			Unit: "pcs",
		}
		resp, err := client.CreateProduct(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateProduct: %v", err)
		}
		if resp.Product == nil || resp.Product.Id == "" {
			t.Fatal("CreateProduct: no id returned")
		}
		productID = resp.Product.Id
		t.Logf("created product: id=%s", productID)
	})

	// Setup: create bin location
	t.Run("SetupBinLocation", func(t *testing.T) {
		req := &inventoryv1.CreateBinLocationRequest{
			WarehouseZone: fmt.Sprintf("ZONE-FLOW-%d", time.Now().UnixNano()),
			Aisle:         "A2",
			Rack:          "R2",
			Shelf:         "S2",
		}
		resp, err := client.CreateBinLocation(authContext(t), req)
		if err != nil {
			t.Fatalf("CreateBinLocation: %v", err)
		}
		if resp.BinLocation == nil || resp.BinLocation.Id == "" {
			t.Fatal("CreateBinLocation: no id returned")
		}
		binID = resp.BinLocation.Id
		t.Logf("created bin location: id=%s", binID)
	})

	if productID == "" || binID == "" {
		t.Fatal("setup failed: missing product or bin id")
	}

	// Scenario 7: Add stock
	t.Run("AddStock", func(t *testing.T) {
		req := &inventoryv1.AddStockRequest{
			ProductId:     productID,
			BinLocationId: binID,
			Quantity:      100,
			LotNumber:     "E2E-LOT-001",
		}
		resp, err := client.AddStock(authContext(t), req)
		if err != nil {
			t.Fatalf("AddStock: %v", err)
		}
		if resp.StockEntry == nil {
			t.Fatal("AddStock: response stock_entry is nil")
		}
		if resp.StockEntry.Id == "" {
			t.Fatal("AddStock: stock entry id is empty")
		}
		if resp.StockEntry.Quantity != 100 {
			t.Errorf("AddStock: quantity = %f, want %f", resp.StockEntry.Quantity, 100.0)
		}
		if resp.StockEntry.ProductId != productID {
			t.Errorf("AddStock: product_id = %q, want %q", resp.StockEntry.ProductId, productID)
		}
		if resp.StockEntry.BinLocationId != binID {
			t.Errorf("AddStock: bin_location_id = %q, want %q", resp.StockEntry.BinLocationId, binID)
		}
		stockID = resp.StockEntry.Id
		t.Logf("created stock entry: id=%s", stockID)
	})

	if stockID == "" {
		t.Fatal("stockID not set, aborting")
	}

	// Scenario 8: Get stock entry
	t.Run("GetStock", func(t *testing.T) {
		resp, err := client.GetStock(authContext(t), &inventoryv1.GetStockRequest{Id: stockID})
		if err != nil {
			t.Fatalf("GetStock: %v", err)
		}
		if resp.StockEntry == nil {
			t.Fatal("GetStock: response stock_entry is nil")
		}
		if resp.StockEntry.Quantity != 100 {
			t.Errorf("GetStock: quantity = %f, want %f", resp.StockEntry.Quantity, 100.0)
		}
	})

	// Scenario 9: List stock with filters
	t.Run("ListStockByProduct", func(t *testing.T) {
		req := &inventoryv1.ListStockRequest{
			ProductId: productID,
			Pagination: &commonv1.Pagination{
				Limit: 100,
				Page:  1,
			},
		}
		resp, err := client.ListStock(authContext(t), req)
		if err != nil {
			t.Fatalf("ListStock: %v", err)
		}
		if len(resp.StockEntries) == 0 {
			t.Fatal("ListStock: empty results, expected at least 1 entry")
		}
		found := false
		for _, e := range resp.StockEntries {
			if e.Id == stockID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListStock: stock entry %s not found in results", stockID)
		}
	})

	// Scenario 10: Adjust stock (positive)
	t.Run("AdjustStockPositive", func(t *testing.T) {
		req := &inventoryv1.AdjustStockRequest{
			Id:     stockID,
			Delta:  50,
			Reason: "INBOUND",
			Notes:  "E2E additional stock",
		}
		resp, err := client.AdjustStock(authContext(t), req)
		if err != nil {
			t.Fatalf("AdjustStock positive: %v", err)
		}
		if resp.StockEntry == nil {
			t.Fatal("AdjustStock: response stock_entry is nil")
		}
		if resp.StockEntry.Quantity != 150 {
			t.Errorf("AdjustStock positive: quantity = %f, want %f", resp.StockEntry.Quantity, 150.0)
		}
		t.Logf("adjusted stock (+50): quantity=%f", resp.StockEntry.Quantity)
	})

	// Scenario 11: Adjust stock (negative)
	t.Run("AdjustStockNegative", func(t *testing.T) {
		req := &inventoryv1.AdjustStockRequest{
			Id:     stockID,
			Delta:  -30,
			Reason: "SALE",
			Notes:  "E2E sale adjustment",
		}
		resp, err := client.AdjustStock(authContext(t), req)
		if err != nil {
			t.Fatalf("AdjustStock negative: %v", err)
		}
		if resp.StockEntry.Quantity != 120 {
			t.Errorf("AdjustStock negative: quantity = %f, want %f", resp.StockEntry.Quantity, 120.0)
		}
		t.Logf("adjusted stock (-30): quantity=%f", resp.StockEntry.Quantity)
	})

	// Scenario 12: Reserve stock
	t.Run("ReserveStock", func(t *testing.T) {
		req := &inventoryv1.ReserveStockRequest{
			Id:       stockID,
			Quantity: 40,
		}
		resp, err := client.ReserveStock(authContext(t), req)
		if err != nil {
			t.Fatalf("ReserveStock: %v", err)
		}
		if resp.StockEntry.ReservedQuantity != 40 {
			t.Errorf("ReserveStock: reserved_quantity = %f, want %f", resp.StockEntry.ReservedQuantity, 40.0)
		}
		t.Logf("reserved stock: reserved=%f", resp.StockEntry.ReservedQuantity)
	})

	// Scenario 14: Reserve more than available → error
	t.Run("ReserveMoreThanAvailable", func(t *testing.T) {
		// Available = quantity - reserved = 120 - 40 = 80
		req := &inventoryv1.ReserveStockRequest{
			Id:       stockID,
			Quantity: 200,
		}
		_, err := client.ReserveStock(authContext(t), req)
		if err == nil {
			t.Fatal("expected error for reserve exceeding available, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		t.Logf("ReserveMoreThanAvailable: code=%v msg=%v", st.Code(), st.Message())
	})

	// Scenario 13: Release stock
	t.Run("ReleaseStock", func(t *testing.T) {
		req := &inventoryv1.ReleaseStockRequest{
			Id:       stockID,
			Quantity: 20,
		}
		resp, err := client.ReleaseStock(authContext(t), req)
		if err != nil {
			t.Fatalf("ReleaseStock: %v", err)
		}
		if resp.StockEntry.ReservedQuantity != 20 {
			t.Errorf("ReleaseStock: reserved_quantity = %f, want %f", resp.StockEntry.ReservedQuantity, 20.0)
		}
		t.Logf("released stock: reserved=%f", resp.StockEntry.ReservedQuantity)
	})
}

// Scenario 14 (standalone): Reserve more than available
// Note: This scenario is also covered in TestInventoryIntegration_FullStockFlow,
// but we keep it here as a standalone validation test.
func TestInventoryIntegration_ReserveExceedsAvailable(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	var stockID string

	// Setup
	t.Run("SetupProductAndStock", func(t *testing.T) {
		prodResp, err := client.CreateProduct(authContext(t), &inventoryv1.CreateProductRequest{
			Sku:  "E2E-RESERVE-EXCEED-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Name: "Reserve Exceed Test",
			Unit: "pcs",
		})
		if err != nil {
			t.Fatalf("CreateProduct: %v", err)
		}
		if prodResp.Product == nil || prodResp.Product.Id == "" {
			t.Fatal("CreateProduct: no id")
		}

		binResp, err := client.CreateBinLocation(authContext(t), &inventoryv1.CreateBinLocationRequest{
			WarehouseZone: fmt.Sprintf("ZONE-EXCEED-%d", time.Now().UnixNano()),
			Aisle:         "A1",
			Rack:          "R1",
			Shelf:         "S1",
		})
		if err != nil {
			t.Fatalf("CreateBinLocation: %v", err)
		}
		if binResp.BinLocation == nil || binResp.BinLocation.Id == "" {
			t.Fatal("CreateBinLocation: no id")
		}

		stockResp, err := client.AddStock(authContext(t), &inventoryv1.AddStockRequest{
			ProductId:     prodResp.Product.Id,
			BinLocationId: binResp.BinLocation.Id,
			Quantity:      10,
			LotNumber:     "E2E-EXCEED-LOT",
		})
		if err != nil {
			t.Fatalf("AddStock: %v", err)
		}
		if stockResp.StockEntry == nil || stockResp.StockEntry.Id == "" {
			t.Fatal("AddStock: no id")
		}
		stockID = stockResp.StockEntry.Id
		t.Logf("created stock entry: id=%s qty=%f", stockID, 10.0)
	})

	if stockID == "" {
		t.Fatal("stockID not set, aborting")
	}

	t.Run("ReserveExceedsAvailable", func(t *testing.T) {
		req := &inventoryv1.ReserveStockRequest{
			Id:       stockID,
			Quantity: 100,
		}
		_, err := client.ReserveStock(authContext(t), req)
		if err == nil {
			t.Fatal("expected FailedPrecondition for reserve exceeding available, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.FailedPrecondition {
			t.Errorf("status code = %v, want %v", st.Code(), codes.FailedPrecondition)
		}
		t.Logf("ReserveExceedsAvailable: got expected error: %v", err)
	})
}

// --- Auth required ---
func TestInventoryIntegration_AuthRequired(t *testing.T) {
	conn := dial(t)
	t.Cleanup(func() { conn.Close() })
	client := inventoryv1.NewInventoryServiceClient(conn)

	t.Run("NoAuthMetadata", func(t *testing.T) {
		req := &inventoryv1.CreateProductRequest{
			Sku:  "E2E-NO-AUTH",
			Name: "No Auth",
			Unit: "pcs",
		}
		_, err := client.CreateProduct(context.Background(), req)
		if err == nil {
			t.Fatal("expected Unauthenticated error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("error is not a gRPC status: %v", err)
		}
		if st.Code() != codes.Unauthenticated {
			t.Errorf("status code = %v, want %v", st.Code(), codes.Unauthenticated)
		}
		t.Logf("NoAuthMetadata: got expected error: %v", err)
	})
}
