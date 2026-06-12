package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	repomocks "github.com/beabys/wms/auth-service/mocks/application/auth/repository"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Tests ---

func TestInviteUser_Success(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)
	userRepo := repomocks.NewUserRepository(t)

	userRepo.On("GetByEmail", mock.Anything, "new@customer.com").Return(nil, nil)
	inviteRepo.On("ExistsPendingByEmail", mock.Anything, "new@customer.com").Return(false, nil)
	inviteRepo.On("Create", mock.Anything, mock.MatchedBy(func(inv *model.InviteToken) bool {
		return inv.Email == "new@customer.com" && inv.InvitedBy == "admin-123" && inv.Status == model.StatusPending
	})).Return(nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, userRepo)

	result, err := uc.InviteUser(context.Background(), command.InviteUserCommand{
		Email:     "new@customer.com",
		InvitedBy: "admin-123",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Token)
	assert.NotEmpty(t, result.InviteLink)
	assert.Greater(t, result.ExpiresAt, int64(0))
	assert.Contains(t, result.InviteLink, result.Token)
}

func TestInviteUser_EmptyEmail(t *testing.T) {
	uc := NewInviteUseCase(&testLogger{}, repomocks.NewInviteRepository(t), repomocks.NewUserRepository(t))
	_, err := uc.InviteUser(context.Background(), command.InviteUserCommand{
		Email:     "",
		InvitedBy: "admin-123",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is required")
}

func TestInviteUser_EmptyInvitedBy(t *testing.T) {
	uc := NewInviteUseCase(&testLogger{}, repomocks.NewInviteRepository(t), repomocks.NewUserRepository(t))
	_, err := uc.InviteUser(context.Background(), command.InviteUserCommand{
		Email:     "user@example.com",
		InvitedBy: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invited_by is required")
}

func TestInviteUser_DuplicateEmail(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)
	userRepo := repomocks.NewUserRepository(t)

	userRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(&model.User{
		ID:    "existing-user",
		Email: "existing@example.com",
	}, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, userRepo)

	_, err := uc.InviteUser(context.Background(), command.InviteUserCommand{
		Email:     "existing@example.com",
		InvitedBy: "admin-123",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already registered")
}

func TestInviteUser_DuplicatePendingInvite(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)
	userRepo := repomocks.NewUserRepository(t)

	userRepo.On("GetByEmail", mock.Anything, "pending@example.com").Return(nil, nil)
	inviteRepo.On("ExistsPendingByEmail", mock.Anything, "pending@example.com").Return(true, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, userRepo)

	_, err := uc.InviteUser(context.Background(), command.InviteUserCommand{
		Email:     "pending@example.com",
		InvitedBy: "admin-123",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already has a pending invite")
}

func TestInviteUser_InvalidEmail(t *testing.T) {
	uc := NewInviteUseCase(&testLogger{}, repomocks.NewInviteRepository(t), repomocks.NewUserRepository(t))
	_, err := uc.InviteUser(context.Background(), command.InviteUserCommand{
		Email:     "not-an-email",
		InvitedBy: "admin-123",
	})
	assert.Error(t, err)
}

// --- ValidateInvite Tests ---

func TestValidateInvite_Success(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	invite := &model.InviteToken{
		ID:        "inv-1",
		Email:     "user@example.com",
		Token:     "valid-token-123",
		InvitedBy: "admin-123",
		Status:    model.StatusPending,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	inviteRepo.On("GetByToken", mock.Anything, "valid-token-123").Return(invite, nil)
	inviteRepo.On("MarkUsed", mock.Anything, "valid-token-123").Return(nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))

	result, err := uc.ValidateInvite(context.Background(), "valid-token-123")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "user@example.com", result.Email)
	assert.Equal(t, "admin-123", result.InvitedBy)

	inviteRepo.AssertCalled(t, "MarkUsed", mock.Anything, "valid-token-123")
}

func TestValidateInvite_EmptyToken(t *testing.T) {
	uc := NewInviteUseCase(&testLogger{}, repomocks.NewInviteRepository(t), repomocks.NewUserRepository(t))
	_, err := uc.ValidateInvite(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invite token is required")
}

func TestValidateInvite_TokenNotFound(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	inviteRepo.On("GetByToken", mock.Anything, "nonexistent-token").Return(nil, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	_, err := uc.ValidateInvite(context.Background(), "nonexistent-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invite token not found")
}

func TestValidateInvite_AlreadyUsed(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	invite := &model.InviteToken{
		ID:        "inv-2",
		Email:     "used@example.com",
		Token:     "used-token",
		InvitedBy: "admin-123",
		Status:    model.StatusUsed,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	inviteRepo.On("GetByToken", mock.Anything, "used-token").Return(invite, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	_, err := uc.ValidateInvite(context.Background(), "used-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already used")
}

func TestValidateInvite_Cancelled(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	invite := &model.InviteToken{
		ID:        "inv-4",
		Email:     "cancelled@example.com",
		Token:     "cancelled-token",
		InvitedBy: "admin-123",
		Status:    model.StatusCancelled,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	inviteRepo.On("GetByToken", mock.Anything, "cancelled-token").Return(invite, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	_, err := uc.ValidateInvite(context.Background(), "cancelled-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already used")
}

func TestValidateInvite_Expired(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	invite := &model.InviteToken{
		ID:        "inv-3",
		Email:     "expired@example.com",
		Token:     "expired-token",
		InvitedBy: "admin-123",
		Status:    model.StatusPending,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		CreatedAt: time.Now().Add(-8 * 24 * time.Hour),
	}

	inviteRepo.On("GetByToken", mock.Anything, "expired-token").Return(invite, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	_, err := uc.ValidateInvite(context.Background(), "expired-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has expired")
}

// --- CancelInvite Tests ---

func TestCancelInvite_Success(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	invite := &model.InviteToken{
		ID:        "inv-1",
		Email:     "cancel@example.com",
		Token:     "cancel-token",
		InvitedBy: "admin-123",
		Status:    model.StatusPending,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	inviteRepo.On("GetByToken", mock.Anything, "cancel-token").Return(invite, nil)
	inviteRepo.On("Cancel", mock.Anything, "cancel-token").Return(nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	err := uc.CancelInvite(context.Background(), command.CancelInviteCommand{Token: "cancel-token"})
	assert.NoError(t, err)

	inviteRepo.AssertCalled(t, "Cancel", mock.Anything, "cancel-token")
}

func TestCancelInvite_EmptyToken(t *testing.T) {
	uc := NewInviteUseCase(&testLogger{}, repomocks.NewInviteRepository(t), repomocks.NewUserRepository(t))
	err := uc.CancelInvite(context.Background(), command.CancelInviteCommand{Token: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invite token is required")
}

func TestCancelInvite_NotFound(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	inviteRepo.On("GetByToken", mock.Anything, "nonexistent").Return(nil, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	err := uc.CancelInvite(context.Background(), command.CancelInviteCommand{Token: "nonexistent"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invite token not found")
}

func TestCancelInvite_AlreadyUsed(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	invite := &model.InviteToken{
		ID:        "inv-2",
		Email:     "used@example.com",
		Token:     "used-token",
		InvitedBy: "admin-123",
		Status:    model.StatusUsed,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	inviteRepo.On("GetByToken", mock.Anything, "used-token").Return(invite, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	err := uc.CancelInvite(context.Background(), command.CancelInviteCommand{Token: "used-token"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be cancelled")
}

func TestCancelInvite_AlreadyCancelled(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	invite := &model.InviteToken{
		ID:        "inv-3",
		Email:     "cancelled@example.com",
		Token:     "cancelled-token",
		InvitedBy: "admin-123",
		Status:    model.StatusCancelled,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	inviteRepo.On("GetByToken", mock.Anything, "cancelled-token").Return(invite, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	err := uc.CancelInvite(context.Background(), command.CancelInviteCommand{Token: "cancelled-token"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be cancelled")
}

func TestCancelInvite_Expired(t *testing.T) {
	inviteRepo := repomocks.NewInviteRepository(t)

	invite := &model.InviteToken{
		ID:        "inv-4",
		Email:     "expired@example.com",
		Token:     "expired-token",
		InvitedBy: "admin-123",
		Status:    model.StatusPending,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		CreatedAt: time.Now().Add(-8 * 24 * time.Hour),
	}

	inviteRepo.On("GetByToken", mock.Anything, "expired-token").Return(invite, nil)

	uc := NewInviteUseCase(&testLogger{}, inviteRepo, repomocks.NewUserRepository(t))
	err := uc.CancelInvite(context.Background(), command.CancelInviteCommand{Token: "expired-token"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot cancel expired invite")
}

// testLogger implements logger.Logger for tests.
type testLogger struct{}

func (l *testLogger) GetLogger() any                                         { return nil }
func (l *testLogger) Debug(msg string, fields ...logger.LogField)            {}
func (l *testLogger) Info(msg string, fields ...logger.LogField)             {}
func (l *testLogger) Warn(msg string, fields ...logger.LogField)             {}
func (l *testLogger) Error(msg string, err error, fields ...logger.LogField) {}
func (l *testLogger) Fatal(msg string, fields ...logger.LogField)            {}
