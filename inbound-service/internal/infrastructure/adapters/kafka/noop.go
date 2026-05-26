package kafka

import (
	"context"
	"fmt"

	"github.com/beabys/wms/inbound-service/internal/application/inbound/usecase"
	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// NoopEventPublisher logs events instead of publishing to Kafka.
// Used in development when Kafka is unavailable.
type NoopEventPublisher struct{}

// NewNoopEventPublisher creates a new NoopEventPublisher.
func NewNoopEventPublisher() *NoopEventPublisher {
	return &NoopEventPublisher{}
}

// PublishSubmitted logs the submitted event.
func (p *NoopEventPublisher) PublishSubmitted(_ context.Context, inbound *model.Inbound) error {
	fmt.Printf("[EVENT] inbound.submitted: id=%s customer=%s items=%d\n",
		inbound.ID, inbound.CustomerID, len(inbound.Items))
	return nil
}

// PublishApproved logs the approved event.
func (p *NoopEventPublisher) PublishApproved(_ context.Context, inbound *model.Inbound) error {
	fmt.Printf("[EVENT] inbound.approved: id=%s customer=%s items=%d\n",
		inbound.ID, inbound.CustomerID, len(inbound.Items))
	return nil
}

// PublishFlagged logs the flagged event.
func (p *NoopEventPublisher) PublishFlagged(_ context.Context, inbound *model.Inbound) error {
	fmt.Printf("[EVENT] inbound.flagged: id=%s customer=%s\n",
		inbound.ID, inbound.CustomerID)
	return nil
}

// Ensure NoopEventPublisher satisfies the EventPublisher interface at compile time.
var _ usecase.EventPublisher = (*NoopEventPublisher)(nil)
