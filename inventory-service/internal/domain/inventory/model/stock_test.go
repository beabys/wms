package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStockEntry(t *testing.T) {
	expiry := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		id            string
		productID     string
		binLocationID string
		quantity      float64
		lotNumber     string
		expiryDate    *time.Time
		wantErr       bool
		errMsg        string
	}{
		{
			name:      "valid entry",
			id:        "stk_1",
			productID: "prod_1",
			quantity:  100,
			lotNumber: "LOT-001",
			wantErr:   false,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
			errMsg:  "stock entry id is required",
		},
		{
			name:      "empty product id",
			id:        "stk_2",
			productID: "",
			wantErr:   true,
			errMsg:    "product id is required",
		},
		{
			name:      "negative quantity",
			id:        "stk_3",
			productID: "prod_1",
			quantity:  -10,
			wantErr:   true,
			errMsg:    "quantity cannot be negative",
		},
		{
			name:          "with bin location and expiry",
			id:            "stk_4",
			productID:     "prod_1",
			binLocationID: "bin_1",
			quantity:      200,
			lotNumber:     "LOT-002",
			expiryDate:    &expiry,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewStockEntry(tt.id, tt.productID, tt.binLocationID, tt.quantity, tt.lotNumber, tt.expiryDate)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, s)
			assert.Equal(t, tt.id, s.ID)
			assert.Equal(t, tt.productID, s.ProductID)
			assert.Equal(t, tt.binLocationID, s.BinLocationID)
			assert.Equal(t, tt.quantity, s.Quantity)
			assert.Equal(t, "active", s.Status)
			if tt.expiryDate != nil {
				assert.True(t, s.ExpiryDate.Equal(*tt.expiryDate))
			}
		})
	}
}

func TestStockEntryAdjust(t *testing.T) {
	s, err := NewStockEntry("stk_1", "prod_1", "", 100, "LOT-001", nil)
	require.NoError(t, err)

	err = s.Adjust(50, "inbound")
	require.NoError(t, err)
	assert.Equal(t, 150.0, s.Quantity)

	err = s.Adjust(-30, "sale")
	require.NoError(t, err)
	assert.Equal(t, 120.0, s.Quantity)

	err = s.Adjust(-200, "disposal")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")

	err = s.Adjust(10, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "adjust reason is required")
}

func TestStockEntryReserveRelease(t *testing.T) {
	s, err := NewStockEntry("stk_1", "prod_1", "", 100, "LOT-001", nil)
	require.NoError(t, err)

	err = s.Reserve(30)
	require.NoError(t, err)
	assert.Equal(t, 30.0, s.ReservedQuantity)

	err = s.Reserve(20)
	require.NoError(t, err)
	assert.Equal(t, 50.0, s.ReservedQuantity)

	err = s.Reserve(60)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient available stock")

	err = s.Release(10)
	require.NoError(t, err)
	assert.Equal(t, 40.0, s.ReservedQuantity)

	err = s.Release(50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only 40")

	err = s.Reserve(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")

	err = s.Release(-5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")
}

func TestStockEntryIsLowStock(t *testing.T) {
	s, err := NewStockEntry("stk_1", "prod_1", "", 100, "LOT-001", nil)
	require.NoError(t, err)

	// available=100, threshold=10 -> 100 >= 10, not low
	assert.False(t, s.IsLowStock(10))

	s.Reserve(95)
	// available=5, threshold=10 -> 5 < 10, low
	assert.True(t, s.IsLowStock(10))

	// available=5, threshold=5 -> 5 < 5, false (not strictly less)
	assert.False(t, s.IsLowStock(5))
	// available=5, threshold=6 -> 5 < 6, true
	assert.True(t, s.IsLowStock(6))
}
