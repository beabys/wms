package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterProduct(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		sku       string
		productName string
		desc      string
		category  string
		unit      string
		weight    float64
		threshold float64
		wantErr   bool
		errMsg    string
	}{
		{
			name:        "valid product",
			id:          "prod_1",
			sku:         "SKU-001",
			productName: "Widget",
			desc:        "A widget",
			category:    "gadgets",
			unit:        "pcs",
			weight:      1.5,
			threshold:   10,
			wantErr:     false,
		},
		{
			name:        "empty id",
			id:          "",
			sku:         "SKU-002",
			productName: "Gadget",
			wantErr:     true,
			errMsg:      "product id is required",
		},
		{
			name:        "empty sku",
			id:          "prod_2",
			sku:         "",
			productName: "Gadget",
			wantErr:     true,
			errMsg:      "sku is required",
		},
		{
			name:        "empty name",
			id:          "prod_3",
			sku:         "SKU-003",
			productName: "",
			wantErr:     true,
			errMsg:      "name is required",
		},
		{
			name:        "negative weight",
			id:          "prod_4",
			sku:         "SKU-004",
			productName: "Heavy",
			weight:      -1,
			wantErr:     true,
			errMsg:      "weight cannot be negative",
		},
		{
			name:        "negative threshold",
			id:          "prod_5",
			sku:         "SKU-005",
			productName: "Critical",
			threshold:   -5,
			wantErr:     true,
			errMsg:      "low stock threshold cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := RegisterProduct(tt.id, tt.sku, tt.productName, tt.desc, tt.category, tt.unit, tt.weight, tt.threshold)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, tt.id, p.ID)
			assert.Equal(t, tt.sku, p.SKU)
			assert.Equal(t, tt.productName, p.Name)
			assert.Equal(t, tt.desc, p.Description)
			assert.Equal(t, tt.category, p.Category)
			assert.Equal(t, tt.unit, p.Unit)
			assert.Equal(t, tt.weight, p.WeightKg)
			assert.Equal(t, tt.threshold, p.LowStockThreshold)
			assert.True(t, p.IsActive)
			assert.False(t, p.CreatedAt.IsZero())
			assert.False(t, p.UpdatedAt.IsZero())
		})
	}
}

func TestProductUpdate(t *testing.T) {
	p, err := RegisterProduct("prod_1", "SKU-001", "Widget", "desc", "cat", "pcs", 1.0, 5)
	require.NoError(t, err)

	err = p.Update("NewName", "new desc", "new cat", "kg", 2.5)
	require.NoError(t, err)
	assert.Equal(t, "NewName", p.Name)
	assert.Equal(t, "kg", p.Unit)
	assert.Equal(t, 2.5, p.WeightKg)

	err = p.Update("", "desc", "cat", "pcs", 1.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	err = p.Update("Valid", "desc", "cat", "pcs", -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "weight cannot be negative")
}

func TestProductArchive(t *testing.T) {
	p, err := RegisterProduct("prod_1", "SKU-001", "Widget", "desc", "cat", "pcs", 1.0, 5)
	require.NoError(t, err)
	assert.True(t, p.IsActive)

	p.Archive()
	assert.False(t, p.IsActive)
}
