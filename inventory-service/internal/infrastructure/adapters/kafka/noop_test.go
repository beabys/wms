package kafka

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/inventory-service/internal/application/inventory/usecase"
	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

func TestNewNoopEventPublisher(t *testing.T) {
	p := NewNoopEventPublisher()
	require.NotNil(t, p)

	// Ensure it satisfies the interface
	var _ usecase.EventPublisher = p
}

func TestNoopPublishStockAdjusted(t *testing.T) {
	p := NewNoopEventPublisher()
	entry := &model.StockEntry{
		ID:        "stk_1",
		ProductID: "prod_1",
		Quantity:  100,
	}
	err := p.PublishStockAdjusted(context.Background(), entry, "MANUAL")
	assert.NoError(t, err)
}

func TestNoopPublishLowStockAlert(t *testing.T) {
	p := NewNoopEventPublisher()
	entry := &model.StockEntry{
		ID:        "stk_1",
		ProductID: "prod_1",
		Quantity:  5,
	}
	err := p.PublishLowStockAlert(context.Background(), entry, 10)
	assert.NoError(t, err)
}
