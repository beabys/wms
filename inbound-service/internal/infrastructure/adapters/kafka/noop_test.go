package kafka

import (
	"context"
	"testing"

	"github.com/beabys/wms/inbound-service/internal/application/inbound/usecase"
	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

func TestNoopEventPublisher_ImplementsInterface(t *testing.T) {
	var _ usecase.EventPublisher = (*NoopEventPublisher)(nil)
}

func TestNoopEventPublisher_ReturnsNil(t *testing.T) {
	ctx := context.Background()
	p := NewNoopEventPublisher()

	t.Run("PublishSubmitted", func(t *testing.T) {
		inbound := &model.Inbound{
			ID:         "inb_1",
			CustomerID: "cust_1",
			Items:      []model.InboundItem{{SKU: "SKU001", QuantityDeclared: 10}},
		}
		if err := p.PublishSubmitted(ctx, inbound); err != nil {
			t.Errorf("PublishSubmitted returned error: %v", err)
		}
	})

	t.Run("PublishApproved", func(t *testing.T) {
		inbound := &model.Inbound{
			ID:         "inb_2",
			CustomerID: "cust_2",
			Items:      []model.InboundItem{{SKU: "SKU002", QuantityDeclared: 5, QuantityReceived: 5, Weight: 2.5}},
		}
		if err := p.PublishApproved(ctx, inbound); err != nil {
			t.Errorf("PublishApproved returned error: %v", err)
		}
	})

	t.Run("PublishFlagged", func(t *testing.T) {
		inbound := &model.Inbound{
			ID:         "inb_3",
			CustomerID: "cust_3",
		}
		if err := p.PublishFlagged(ctx, inbound); err != nil {
			t.Errorf("PublishFlagged returned error: %v", err)
		}
	})
}

func TestNoopEventPublisher_NilItems(t *testing.T) {
	ctx := context.Background()
	p := NewNoopEventPublisher()

	t.Run("PublishSubmitted with nil items", func(t *testing.T) {
		inbound := &model.Inbound{
			ID:         "inb_4",
			CustomerID: "cust_4",
		}
		if err := p.PublishSubmitted(ctx, inbound); err != nil {
			t.Errorf("PublishSubmitted with nil items returned error: %v", err)
		}
	})

	t.Run("PublishApproved with nil items", func(t *testing.T) {
		inbound := &model.Inbound{
			ID:         "inb_5",
			CustomerID: "cust_5",
		}
		if err := p.PublishApproved(ctx, inbound); err != nil {
			t.Errorf("PublishApproved with nil items returned error: %v", err)
		}
	})
}
