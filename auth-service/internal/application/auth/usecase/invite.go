package usecase

import (
	"context"
	"fmt"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/application/auth/validator"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	"github.com/beabys/wms/pkg/logger"
)

// InviteUseCase handles invite token operations.
type InviteUseCase struct {
	logger     logger.Logger
	inviteRepo repository.InviteRepository
	userRepo   repository.UserRepository
}

// NewInviteUseCase creates a new InviteUseCase.
func NewInviteUseCase(log logger.Logger, inviteRepo repository.InviteRepository, userRepo repository.UserRepository) *InviteUseCase {
	return &InviteUseCase{
		logger:     log,
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
	}
}

// InviteUser generates an invite token for a new user.
func (uc *InviteUseCase) InviteUser(ctx context.Context, cmd command.InviteUserCommand) (*command.InviteUserResult, error) {
	if err := validator.Email(cmd.Email); err != nil {
		return nil, err
	}
	if err := validator.NotEmpty(cmd.InvitedBy, "invited_by"); err != nil {
		return nil, err
	}

	// 1. Check email is not already registered
	existing, err := uc.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		uc.logger.Error("failed to check email uniqueness", err)
		return nil, fmt.Errorf("failed to create invite")
	}
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// 2. Check if email already has a pending invite
	hasPending, err := uc.inviteRepo.ExistsPendingByEmail(ctx, cmd.Email)
	if err != nil {
		uc.logger.Error("failed to check pending invite", err)
		return nil, fmt.Errorf("failed to create invite")
	}
	if hasPending {
		return nil, fmt.Errorf("email already has a pending invite")
	}

	// 3. Create InviteToken entity (default TTL = 7 days)
	invite, err := model.NewInviteToken(cmd.Email, cmd.InvitedBy, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid invite data: %w", err)
	}

	// 4. Save via inviteRepo
	if err := uc.inviteRepo.Create(ctx, invite); err != nil {
		uc.logger.Error("failed to save invite token", err)
		return nil, fmt.Errorf("failed to create invite")
	}

	uc.logger.Info("invite token created",
		logger.LogField{Key: "email", Value: cmd.Email},
		logger.LogField{Key: "invited_by", Value: cmd.InvitedBy},
	)

	// 5. Build result
	inviteLink := fmt.Sprintf("https://portal.example.com/register?token=%s", invite.Token)

	return &command.InviteUserResult{
		Token:      invite.Token,
		InviteLink: inviteLink,
		ExpiresAt:  invite.ExpiresAt.Unix(),
	}, nil
}

// ListInvites retrieves paginated invite tokens with optional filters.
func (uc *InviteUseCase) ListInvites(ctx context.Context, filter repository.InviteFilter) (*repository.InviteListResult, error) {
	uc.logger.Info("listing invites",
		logger.LogField{Key: "page", Value: filter.Page},
		logger.LogField{Key: "page_size", Value: filter.PageSize},
	)

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	result, err := uc.inviteRepo.List(ctx, filter)
	if err != nil {
		uc.logger.Error("failed to list invites", err)
		return nil, fmt.Errorf("failed to list invites: %w", err)
	}

	return result, nil
}

// ValidateInvite validates an invite token and marks it as used.
func (uc *InviteUseCase) ValidateInvite(ctx context.Context, token string) (*command.ValidateInviteResult, error) {
	if err := validator.InviteToken(token); err != nil {
		return nil, err
	}

	invite, err := uc.inviteRepo.GetByToken(ctx, token)
	if err != nil {
		uc.logger.Error("failed to get invite token", err)
		return nil, fmt.Errorf("invalid invite token")
	}
	if invite == nil {
		return nil, fmt.Errorf("invite token not found")
	}
	if invite.Status != model.StatusPending {
		return nil, fmt.Errorf("invite token already used")
	}
	if invite.IsExpired() {
		return nil, fmt.Errorf("invite token has expired")
	}

	// Mark as used
	if err := uc.inviteRepo.MarkUsed(ctx, token); err != nil {
		uc.logger.Error("failed to mark invite as used", err)
		return nil, fmt.Errorf("failed to validate invite")
	}

	uc.logger.Info("invite validated",
		logger.LogField{Key: "email", Value: invite.Email},
		logger.LogField{Key: "invited_by", Value: invite.InvitedBy},
	)

	return &command.ValidateInviteResult{
		Email:     invite.Email,
		InvitedBy: invite.InvitedBy,
	}, nil
}

// CancelInvite cancels a pending invite token.
func (uc *InviteUseCase) CancelInvite(ctx context.Context, cmd command.CancelInviteCommand) error {
	if err := validator.InviteToken(cmd.Token); err != nil {
		return err
	}

	invite, err := uc.inviteRepo.GetByToken(ctx, cmd.Token)
	if err != nil {
		uc.logger.Error("failed to get invite token", err)
		return fmt.Errorf("invalid invite token")
	}
	if invite == nil {
		return fmt.Errorf("invite token not found")
	}
	if ok, reason := invite.CanCancel(); !ok {
		return fmt.Errorf("invite token cannot be cancelled: %s", reason)
	}

	if err := uc.inviteRepo.Cancel(ctx, cmd.Token); err != nil {
		uc.logger.Error("failed to cancel invite token", err)
		return fmt.Errorf("failed to cancel invite")
	}

	uc.logger.Info("invite cancelled",
		logger.LogField{Key: "email", Value: invite.Email},
		logger.LogField{Key: "token", Value: cmd.Token},
	)

	return nil
}
