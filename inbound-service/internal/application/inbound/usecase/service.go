package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/beabys/wms/inbound-service/internal/application/inbound/command"
	"github.com/beabys/wms/inbound-service/internal/application/inbound/repository"
	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// NewInboundService creates a new InboundService.
func NewInboundService(repo repository.InboundRepository, events EventPublisher) *InboundService {
	return &InboundService{
		repo:   repo,
		events: events,
	}
}

// SubmitInbound handles SubmitInboundCommand.
func (s *InboundService) SubmitInbound(ctx context.Context, cmd command.SubmitInboundCommand) (*command.InboundResult, error) {
	id := newID()

	items := make([]model.InboundItem, len(cmd.Items))
	for i, si := range cmd.Items {
		items[i] = model.InboundItem{
			ID:               newID(),
			SKU:              si.SKU,
			QuantityDeclared: si.QuantityDeclared,
			QuantityReceived: si.QuantityReceived,
			Dimensions:       si.Dimensions,
			Weight:           si.Weight,
		}
	}

	inbound, err := model.SubmitInbound(id, cmd.CustomerID, cmd.ExpectedDate, cmd.Notes, items)
	if err != nil {
		return nil, fmt.Errorf("submit inbound: %w", err)
	}

	if err := s.repo.Save(ctx, inbound); err != nil {
		return nil, fmt.Errorf("save inbound: %w", err)
	}

	if err := s.events.PublishSubmitted(ctx, inbound); err != nil {
		return nil, fmt.Errorf("publish submitted event: %w", err)
	}

	return &command.InboundResult{Inbound: inbound}, nil
}

// InspectInbound handles InspectInboundCommand.
func (s *InboundService) InspectInbound(ctx context.Context, cmd command.InspectInboundCommand) (*command.InboundResult, error) {
	inbound, err := s.repo.GetByID(ctx, cmd.InboundID)
	if err != nil {
		return nil, fmt.Errorf("get inbound: %w", err)
	}

	inspection := model.Inspection{
		ID:          newID(),
		InspectorID: cmd.InspectorID,
		Notes:       cmd.Notes,
		Photos:      cmd.Photos,
		InspectedAt: cmd.InspectedAt,
		Passed:      cmd.Passed,
	}

	if err := inbound.Inspect(inspection); err != nil {
		return nil, fmt.Errorf("inspect inbound: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, inbound); err != nil {
		return nil, fmt.Errorf("update inbound: %w", err)
	}

	return &command.InboundResult{Inbound: inbound}, nil
}

// ApproveInbound handles ApproveInboundCommand.
func (s *InboundService) ApproveInbound(ctx context.Context, cmd command.ApproveInboundCommand) (*command.InboundResult, error) {
	inbound, err := s.repo.GetByID(ctx, cmd.InboundID)
	if err != nil {
		return nil, fmt.Errorf("get inbound: %w", err)
	}

	if err := inbound.Approve(); err != nil {
		return nil, fmt.Errorf("approve inbound: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, inbound); err != nil {
		return nil, fmt.Errorf("update inbound: %w", err)
	}

	if err := s.events.PublishApproved(ctx, inbound); err != nil {
		return nil, fmt.Errorf("publish approved event: %w", err)
	}

	return &command.InboundResult{Inbound: inbound}, nil
}

// FlagInbound handles FlagInboundCommand.
func (s *InboundService) FlagInbound(ctx context.Context, cmd command.FlagInboundCommand) (*command.InboundResult, error) {
	inbound, err := s.repo.GetByID(ctx, cmd.InboundID)
	if err != nil {
		return nil, fmt.Errorf("get inbound: %w", err)
	}

	if err := inbound.Flag(cmd.Reason); err != nil {
		return nil, fmt.Errorf("flag inbound: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, inbound); err != nil {
		return nil, fmt.Errorf("update inbound: %w", err)
	}

	if err := s.events.PublishFlagged(ctx, inbound); err != nil {
		return nil, fmt.Errorf("publish flagged event: %w", err)
	}

	return &command.InboundResult{Inbound: inbound}, nil
}

// HoldInbound handles HoldInboundCommand.
func (s *InboundService) HoldInbound(ctx context.Context, cmd command.HoldInboundCommand) (*command.InboundResult, error) {
	inbound, err := s.repo.GetByID(ctx, cmd.InboundID)
	if err != nil {
		return nil, fmt.Errorf("get inbound: %w", err)
	}

	if err := inbound.PlaceHold(cmd.Reason); err != nil {
		return nil, fmt.Errorf("hold inbound: %w", err)
	}
	inbound.HoldRecord.ID = newID()

	if err := s.repo.UpdateStatus(ctx, inbound); err != nil {
		return nil, fmt.Errorf("update inbound: %w", err)
	}

	return &command.InboundResult{Inbound: inbound}, nil
}

// ReleaseInbound handles ReleaseInboundCommand.
func (s *InboundService) ReleaseInbound(ctx context.Context, cmd command.ReleaseInboundCommand) (*command.InboundResult, error) {
	inbound, err := s.repo.GetByID(ctx, cmd.InboundID)
	if err != nil {
		return nil, fmt.Errorf("get inbound: %w", err)
	}

	if err := inbound.Release(cmd.ReleasedBy); err != nil {
		return nil, fmt.Errorf("release inbound: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, inbound); err != nil {
		return nil, fmt.Errorf("update inbound: %w", err)
	}

	return &command.InboundResult{Inbound: inbound}, nil
}

// GetInbound handles GetInboundQuery.
func (s *InboundService) GetInbound(ctx context.Context, qry command.GetInboundQuery) (*command.InboundResult, error) {
	inbound, err := s.repo.GetByID(ctx, qry.InboundID)
	if err != nil {
		return nil, fmt.Errorf("get inbound: %w", err)
	}

	return &command.InboundResult{Inbound: inbound}, nil
}

// ListInbounds handles ListInboundsQuery.
func (s *InboundService) ListInbounds(ctx context.Context, qry command.ListInboundsQuery) (*command.ListInboundsResult, error) {
	inbounds, nextToken, err := s.repo.List(ctx, qry.CustomerID, qry.Status, qry.PageSize, qry.PageToken)
	if err != nil {
		return nil, fmt.Errorf("list inbounds: %w", err)
	}

	return &command.ListInboundsResult{
		Inbounds:      inbounds,
		NextPageToken: nextToken,
	}, nil
}

// newID generates a unique ID using UUID.
func newID() string {
	return uuid.NewString()
}
