package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// Topics
const (
	TopicStockAdjusted = "wms.inventory.stock-adjusted"
	TopicLowStockAlert = "wms.inventory.low-stock-alert"
)

// EventPublisher publishes domain events to Kafka.
type EventPublisher struct {
	writer *kafka.Writer
}

// NewEventPublisher creates a new Kafka event publisher.
func NewEventPublisher(brokers []string) *EventPublisher {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}
	return &EventPublisher{writer: w}
}

// Close closes the underlying writer.
func (p *EventPublisher) Close() error {
	return p.writer.Close()
}

// StockAdjustedEvent is the payload for wms.inventory.stock-adjusted.
type StockAdjustedEvent struct {
	StockEntryID string    `json:"stock_entry_id"`
	ProductID    string    `json:"product_id"`
	Quantity     float64   `json:"quantity"`
	Delta        float64   `json:"delta"`
	Reason       string    `json:"reason"`
	Timestamp    time.Time `json:"timestamp"`
}

// LowStockAlertEvent is the payload for wms.inventory.low-stock-alert.
type LowStockAlertEvent struct {
	StockEntryID string    `json:"stock_entry_id"`
	ProductID    string    `json:"product_id"`
	Quantity     float64   `json:"quantity"`
	Threshold    float64   `json:"threshold"`
	Timestamp    time.Time `json:"timestamp"`
}

// PublishStockAdjusted publishes a stock adjusted event.
func (p *EventPublisher) PublishStockAdjusted(ctx context.Context, entry *model.StockEntry, reason string) error {
	evt := StockAdjustedEvent{
		StockEntryID: entry.ID,
		ProductID:    entry.ProductID,
		Quantity:     entry.Quantity,
		Delta:        entry.Quantity, // simplified; real impl tracks delta separately
		Reason:       reason,
		Timestamp:    time.Now(),
	}
	return p.publish(ctx, TopicStockAdjusted, entry.ID, evt)
}

// PublishLowStockAlert publishes a low stock alert event.
func (p *EventPublisher) PublishLowStockAlert(ctx context.Context, entry *model.StockEntry, threshold float64) error {
	evt := LowStockAlertEvent{
		StockEntryID: entry.ID,
		ProductID:    entry.ProductID,
		Quantity:     entry.Quantity,
		Threshold:    threshold,
		Timestamp:    time.Now(),
	}
	return p.publish(ctx, TopicLowStockAlert, entry.ID, evt)
}

func (p *EventPublisher) publish(ctx context.Context, topic, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(topic)},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}

	return nil
}
