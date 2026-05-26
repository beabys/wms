package usecase

import (
	"context"

	"github.com/beabys/wms/inbound-service/internal/application/inbound/command"
	"github.com/beabys/wms/inbound-service/internal/application/inbound/repository"
	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// EventPublisher is the port for publishing domain events.
type EventPublisher interface {
	PublishSubmitted(ctx context.Context, inbound *model.Inbound) error
	PublishApproved(ctx context.Context, inbound *model.Inbound) error
	PublishFlagged(ctx context.Context, inbound *model.Inbound) error
}

// InboundServiceHandler defines the application service contract for inbound operations.
type InboundServiceHandler interface {
	SubmitInbound(ctx context.Context, cmd command.SubmitInboundCommand) (*command.InboundResult, error)
	InspectInbound(ctx context.Context, cmd command.InspectInboundCommand) (*command.InboundResult, error)
	ApproveInbound(ctx context.Context, cmd command.ApproveInboundCommand) (*command.InboundResult, error)
	FlagInbound(ctx context.Context, cmd command.FlagInboundCommand) (*command.InboundResult, error)
	HoldInbound(ctx context.Context, cmd command.HoldInboundCommand) (*command.InboundResult, error)
	ReleaseInbound(ctx context.Context, cmd command.ReleaseInboundCommand) (*command.InboundResult, error)
	GetInbound(ctx context.Context, qry command.GetInboundQuery) (*command.InboundResult, error)
	ListInbounds(ctx context.Context, qry command.ListInboundsQuery) (*command.ListInboundsResult, error)
}

// InboundService implements InboundServiceHandler.
type InboundService struct {
	repo   repository.InboundRepository
	events EventPublisher
}
