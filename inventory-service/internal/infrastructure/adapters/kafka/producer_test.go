package kafka

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

func TestNewEventPublisher(t *testing.T) {
	p := NewEventPublisher([]string{"localhost:9092"})
	require.NotNil(t, p)
	assert.NotNil(t, p.writer)
	defer p.Close()
}

func TestEventPublisherPublishStockAdjusted(t *testing.T) {
	// Using nil writer would panic, so we test with a properly initialized publisher.
	// This is a compile-time / interface conformance test.
	p := NewEventPublisher([]string{"localhost:9092"})
	defer p.Close()

	entry := &model.StockEntry{
		ID:        "stk_1",
		ProductID: "prod_1",
		Quantity:  100,
	}

	// We expect an error because no Kafka broker is available
	err := p.PublishStockAdjusted(context.Background(), entry, "MANUAL")
	assert.Error(t, err)
}

func TestEventPublisherPublishLowStockAlert(t *testing.T) {
	p := NewEventPublisher([]string{"localhost:9092"})
	defer p.Close()

	entry := &model.StockEntry{
		ID:        "stk_1",
		ProductID: "prod_1",
		Quantity:  5,
	}

	err := p.PublishLowStockAlert(context.Background(), entry, 10)
	assert.Error(t, err)
}
