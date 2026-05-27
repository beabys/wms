package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/beabys/wms/inventory-service/internal/application/inventory/command"
	"github.com/beabys/wms/inventory-service/internal/application/inventory/usecase"
)

// TopicInboundApproved is the Kafka topic for inbound approved events.
const TopicInboundApproved = "wms.inbound.approved"

// InboundApprovedEvent matches the event published by inbound-service.
type InboundApprovedEvent struct {
	InboundID  string         `json:"inbound_id"`
	CustomerID string         `json:"customer_id"`
	Items      []InboundItem  `json:"items"`
	Timestamp  time.Time      `json:"timestamp"`
}

// InboundItem represents a single item in an approved inbound event.
type InboundItem struct {
	SKU              string  `json:"sku"`
	QuantityDeclared int32   `json:"quantity_declared"`
	QuantityReceived int32   `json:"quantity_received"`
	Weight           float64 `json:"weight"`
}

// messageReader is an interface for reading Kafka messages, enabling test mocks.
type messageReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

// InboundEventConsumer consumes wms.inbound.approved events and creates stock.
type InboundEventConsumer struct {
	reader messageReader
	invSvc *usecase.InventoryService
	logger *zap.Logger
}

// NewInboundEventConsumer creates a new consumer for wms.inbound.approved.
func NewInboundEventConsumer(brokers []string, svc *usecase.InventoryService, logger *zap.Logger) *InboundEventConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     TopicInboundApproved,
		GroupID:   "inventory-svc",
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	return &InboundEventConsumer{
		reader: r,
		invSvc: svc,
		logger: logger,
	}
}

// Run starts consuming messages in an errgroup goroutine.
func (c *InboundEventConsumer) Run(ctx context.Context, wg *errgroup.Group) error {
	wg.Go(func() error {
		c.logger.Info("kafka consumer started",
			zap.String("topic", TopicInboundApproved),
			zap.String("group_id", "inventory-svc"),
		)
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("kafka consumer shutting down")
				return nil
			default:
			}

			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				c.logger.Error("read kafka message", zap.Error(err))
				continue
			}
			c.processMessage(ctx, msg)
		}
	})
	return nil
}

// processMessage handles a single Kafka message.
func (c *InboundEventConsumer) processMessage(ctx context.Context, msg kafka.Message) {
	var evt InboundApprovedEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		c.logger.Error("unmarshal inbound approved event",
			zap.String("key", string(msg.Key)),
			zap.Error(err),
		)
		return
	}

	c.logger.Info("processing inbound approved event",
		zap.String("inbound_id", evt.InboundID),
		zap.Int("items", len(evt.Items)),
	)

	for _, item := range evt.Items {
		c.processItem(ctx, evt.InboundID, item)
	}
}

// processItem handles stock creation for a single inbound item.
func (c *InboundEventConsumer) processItem(ctx context.Context, inboundID string, item InboundItem) {
	// Look up the product by SKU (SKU is used as the external product identity).
	prodResult, err := c.invSvc.GetProduct(ctx, command.GetProductQuery{ProductID: item.SKU})
	if err != nil {
		// Product not found — create a placeholder
		c.logger.Info("product not found, creating placeholder",
			zap.String("sku", item.SKU),
			zap.String("inbound_id", inboundID),
		)
		prodResult, err = c.invSvc.CreateProduct(ctx, command.CreateProductCommand{
			SKU:               item.SKU,
			Name:              fmt.Sprintf("Product %s", item.SKU),
			Description:       fmt.Sprintf("Auto-created from inbound event %s", inboundID),
			Category:          "general",
			Unit:              "unit",
			WeightKg:          float64(item.Weight),
			LowStockThreshold: 0,
		})
		if err != nil {
			c.logger.Error("create placeholder product",
				zap.String("sku", item.SKU),
				zap.Error(err),
			)
			return
		}
	}

	// Add stock for this product using the actual product ID
	_, err = c.invSvc.AddStock(ctx, command.AddStockCommand{
		ProductID: prodResult.Product.ID,
		Quantity:  float64(item.QuantityReceived),
	})
	if err != nil {
		c.logger.Error("add stock from inbound event",
			zap.String("sku", item.SKU),
			zap.String("product_id", prodResult.Product.ID),
			zap.Int32("quantity", item.QuantityReceived),
			zap.Error(err),
		)
		return
	}

	c.logger.Info("stock added from inbound event",
		zap.String("sku", item.SKU),
		zap.String("product_id", prodResult.Product.ID),
		zap.Int32("quantity", item.QuantityReceived),
		zap.String("inbound_id", inboundID),
	)
}

// Close closes the underlying Kafka reader.
func (c *InboundEventConsumer) Close() error {
	return c.reader.Close()
}
