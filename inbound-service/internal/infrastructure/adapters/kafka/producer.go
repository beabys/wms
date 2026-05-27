package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// Topics
const (
	TopicInboundSubmitted = "wms.inbound.submitted"
	TopicInboundApproved  = "wms.inbound.approved"
	TopicInboundFlagged   = "wms.inbound.flagged"
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

// SubmittedEvent is the payload for wms.inbound.submitted.
type SubmittedEvent struct {
	InboundID    string    `json:"inbound_id"`
	CustomerID   string    `json:"customer_id"`
	ItemsCount   int       `json:"items_count"`
	ExpectedDate string    `json:"expected_date"`
	Timestamp    time.Time `json:"timestamp"`
}

// ApprovedEvent is the payload for wms.inbound.approved.
type ApprovedEvent struct {
	InboundID  string           `json:"inbound_id"`
	CustomerID string           `json:"customer_id"`
	Items      []ApprovedItem   `json:"items"`
	Timestamp  time.Time        `json:"timestamp"`
}

// ApprovedItem represents an item in the approved event.
type ApprovedItem struct {
	SKU              string  `json:"sku"`
	QuantityDeclared int32   `json:"quantity_declared"`
	QuantityReceived int32   `json:"quantity_received"`
	Weight           float64 `json:"weight"`
}

// FlaggedEvent is the payload for wms.inbound.flagged.
type FlaggedEvent struct {
	InboundID string    `json:"inbound_id"`
	CustomerID string   `json:"customer_id"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// PublishSubmitted publishes a submitted event.
func (p *EventPublisher) PublishSubmitted(ctx context.Context, inbound *model.Inbound) error {
	evt := SubmittedEvent{
		InboundID:    inbound.ID,
		CustomerID:   inbound.CustomerID,
		ItemsCount:   len(inbound.Items),
		ExpectedDate: inbound.ExpectedDate,
		Timestamp:    time.Now(),
	}
	return p.publish(ctx, TopicInboundSubmitted, inbound.ID, evt)
}

// PublishApproved publishes an approved event.
func (p *EventPublisher) PublishApproved(ctx context.Context, inbound *model.Inbound) error {
	items := make([]ApprovedItem, len(inbound.Items))
	for i, item := range inbound.Items {
		items[i] = ApprovedItem{
			SKU:              item.SKU,
			QuantityDeclared: item.QuantityDeclared,
			QuantityReceived: item.QuantityReceived,
			Weight:           item.Weight,
		}
	}

	evt := ApprovedEvent{
		InboundID:  inbound.ID,
		CustomerID: inbound.CustomerID,
		Items:      items,
		Timestamp:  time.Now(),
	}
	return p.publish(ctx, TopicInboundApproved, inbound.ID, evt)
}

// PublishFlagged publishes a flagged event.
func (p *EventPublisher) PublishFlagged(ctx context.Context, inbound *model.Inbound) error {
	evt := FlaggedEvent{
		InboundID:  inbound.ID,
		CustomerID: inbound.CustomerID,
		Reason:     "flagged", // generic reason; could be enhanced
		Timestamp:  time.Now(),
	}
	return p.publish(ctx, TopicInboundFlagged, inbound.ID, evt)
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
