package repository

import (
	"context"

	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// InboundRepository defines the persistence contract for Inbound aggregates.
// Defined on the consumer (application) side per DDD/hexagonal principles.
type InboundRepository interface {
	// GetByID retrieves an inbound aggregate by ID.
	GetByID(ctx context.Context, id string) (*model.Inbound, error)

	// Save persists a new inbound aggregate.
	Save(ctx context.Context, inbound *model.Inbound) error

	// UpdateStatus updates the inbound status and saves related entities.
	UpdateStatus(ctx context.Context, inbound *model.Inbound) error

	// List returns inbounds matching the given filters.
	List(ctx context.Context, customerID, status string, pageSize int32, pageToken string) ([]*model.Inbound, string, error)
}
