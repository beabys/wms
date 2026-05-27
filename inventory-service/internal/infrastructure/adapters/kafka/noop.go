package kafka

import (
	"context"
	"fmt"

	"github.com/beabys/wms/inventory-service/internal/application/inventory/usecase"
	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// NoopEventPublisher logs events instead of publishing to Kafka.
// Used in development when Kafka is unavailable.
type NoopEventPublisher struct{}

// NewNoopEventPublisher creates a new NoopEventPublisher.
func NewNoopEventPublisher() *NoopEventPublisher {
	return &NoopEventPublisher{}
}

// PublishStockAdjusted logs the stock adjusted event.
func (p *NoopEventPublisher) PublishStockAdjusted(_ context.Context, entry *model.StockEntry, reason string) error {
	fmt.Printf("[EVENT] inventory.stock-adjusted: id=%s product=%s qty=%f reason=%s\n",
		entry.ID, entry.ProductID, entry.Quantity, reason)
	return nil
}

// PublishLowStockAlert logs the low stock alert event.
func (p *NoopEventPublisher) PublishLowStockAlert(_ context.Context, entry *model.StockEntry, threshold float64) error {
	fmt.Printf("[EVENT] inventory.low-stock-alert: id=%s product=%s qty=%f threshold=%f\n",
		entry.ID, entry.ProductID, entry.Quantity, threshold)
	return nil
}

// Ensure NoopEventPublisher satisfies the EventPublisher interface at compile time.
var _ usecase.EventPublisher = (*NoopEventPublisher)(nil)
