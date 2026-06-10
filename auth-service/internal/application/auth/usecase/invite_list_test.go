package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	repomocks "github.com/beabys/wms/auth-service/mocks/application/auth/repository"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
)

func setupInviteListUC(t *testing.T) (*InviteUseCase, *repomocks.InviteRepository) {
	t.Helper()
	log := logger.NewSlogLogger(slog.LevelDebug)
	inviteRepo := repomocks.NewInviteRepository(t)
	userRepo := repomocks.NewUserRepository(t)
	return NewInviteUseCase(log, inviteRepo, userRepo), inviteRepo
}

func seedInvites() []*model.InviteToken {
	now := time.Now()
	return []*model.InviteToken{
		{ID: "1", Email: "used@example.com", Token: "token-used", InvitedBy: "admin", Status: model.StatusUsed, ExpiresAt: now.Add(24 * time.Hour), CreatedAt: now.Add(-48 * time.Hour)},
		{ID: "2", Email: "pending@example.com", Token: "token-pending", InvitedBy: "admin", Status: model.StatusPending, ExpiresAt: now.Add(24 * time.Hour), CreatedAt: now.Add(-24 * time.Hour)},
		{ID: "3", Email: "expired@example.com", Token: "token-expired", InvitedBy: "admin", Status: model.StatusPending, ExpiresAt: now.Add(-24 * time.Hour), CreatedAt: now.Add(-72 * time.Hour)},
		{ID: "4", Email: "active@example.com", Token: "token-active", InvitedBy: "admin", Status: model.StatusPending, ExpiresAt: now.Add(48 * time.Hour), CreatedAt: now.Add(-1 * time.Hour)},
		{ID: "5", Email: "cancelled@example.com", Token: "token-cancelled", InvitedBy: "admin", Status: model.StatusCancelled, ExpiresAt: now.Add(24 * time.Hour), CreatedAt: now.Add(-12 * time.Hour)},
	}
}

func TestListInvites_NoFilters(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Page: 1, PageSize: 20}).Return(&repository.InviteListResult{
		Invites:    invites,
		TotalItems: 5,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 5, result.TotalItems)
	assert.Len(t, result.Invites, 5)
}

func TestListInvites_FilterByStatus(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	status := model.StatusUsed
	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Status: &status, Page: 1, PageSize: 20}).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{invites[0]},
		TotalItems: 1,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{Status: &status, Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalItems)
	assert.Len(t, result.Invites, 1)
	assert.Equal(t, "used@example.com", result.Invites[0].Email)
}

func TestListInvites_FilterPending(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	status := model.StatusPending
	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Status: &status, Page: 1, PageSize: 20}).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{invites[1], invites[2], invites[3]},
		TotalItems: 3,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{Status: &status, Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalItems)
	assert.Len(t, result.Invites, 3)
}

func TestListInvites_FilterCancelled(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	status := model.StatusCancelled
	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Status: &status, Page: 1, PageSize: 20}).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{invites[4]},
		TotalItems: 1,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{Status: &status, Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalItems)
	assert.Len(t, result.Invites, 1)
	assert.Equal(t, "cancelled@example.com", result.Invites[0].Email)
}

func TestListInvites_FilterExpired(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	expired := true
	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Expired: &expired, Page: 1, PageSize: 20}).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{invites[2]},
		TotalItems: 1,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{Expired: &expired, Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalItems)
	assert.Len(t, result.Invites, 1)
	assert.Equal(t, "expired@example.com", result.Invites[0].Email)
}

func TestListInvites_FilterNotExpired(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	expired := false
	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Expired: &expired, Page: 1, PageSize: 20}).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{invites[0], invites[1], invites[3], invites[4]},
		TotalItems: 4,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{Expired: &expired, Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 4, result.TotalItems)
	assert.Len(t, result.Invites, 4)
}

func TestListInvites_FilterCreatedAfter(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	after := time.Now().Add(-48 * time.Hour)
	inviteRepo.On("List", mock.Anything, mock.MatchedBy(func(f repository.InviteFilter) bool {
		return f.CreatedAfter != nil && f.Page == 1 && f.PageSize == 20
	})).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{invites[1], invites[3], invites[4]},
		TotalItems: 3,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{CreatedAfter: &after, Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalItems)
	assert.Len(t, result.Invites, 3)
}

func TestListInvites_FilterCreatedBefore(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	before := time.Now().Add(-2 * time.Hour)
	inviteRepo.On("List", mock.Anything, mock.MatchedBy(func(f repository.InviteFilter) bool {
		return f.CreatedBefore != nil && f.Page == 1 && f.PageSize == 20
	})).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{invites[0], invites[1], invites[2], invites[4]},
		TotalItems: 4,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{CreatedBefore: &before, Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 4, result.TotalItems)
	assert.Len(t, result.Invites, 4)
}

func TestListInvites_WithPagination(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)

	// Page 1, PageSize 2
	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Page: 1, PageSize: 2}).Return(&repository.InviteListResult{
		Invites: []*model.InviteToken{
			{ID: "a", Email: "page_test@example.com", Token: "token-page-a", InvitedBy: "admin", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
			{ID: "b", Email: "page_test@example.com", Token: "token-page-b", InvitedBy: "admin", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
		},
		TotalItems: 5,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{Page: 1, PageSize: 2})
	require.NoError(t, err)
	assert.Equal(t, 5, result.TotalItems)
	assert.Len(t, result.Invites, 2)

	// Page 2, PageSize 2
	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Page: 2, PageSize: 2}).Return(&repository.InviteListResult{
		Invites: []*model.InviteToken{
			{ID: "c", Email: "page_test@example.com", Token: "token-page-c", InvitedBy: "admin", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
			{ID: "d", Email: "page_test@example.com", Token: "token-page-d", InvitedBy: "admin", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
		},
		TotalItems: 5,
	}, nil)

	result, err = uc.ListInvites(context.Background(), repository.InviteFilter{Page: 2, PageSize: 2})
	require.NoError(t, err)
	assert.Equal(t, 5, result.TotalItems)
	assert.Len(t, result.Invites, 2)

	// Page 3, PageSize 2
	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Page: 3, PageSize: 2}).Return(&repository.InviteListResult{
		Invites: []*model.InviteToken{
			{ID: "e", Email: "page_test@example.com", Token: "token-page-e", InvitedBy: "admin", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
		},
		TotalItems: 5,
	}, nil)

	result, err = uc.ListInvites(context.Background(), repository.InviteFilter{Page: 3, PageSize: 2})
	require.NoError(t, err)
	assert.Equal(t, 5, result.TotalItems)
	assert.Len(t, result.Invites, 1)
}

func TestListInvites_EmptyRepo(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)

	inviteRepo.On("List", mock.Anything, repository.InviteFilter{Page: 1, PageSize: 20}).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{},
		TotalItems: 0,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 0, result.TotalItems)
	assert.Empty(t, result.Invites)
}

func TestListInvites_DefaultsPagination(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	inviteRepo.On("List", mock.Anything, repository.InviteFilter{
		Page:     1,
		PageSize: 20,
	}).Return(&repository.InviteListResult{
		Invites:    invites,
		TotalItems: 5,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{})
	require.NoError(t, err)
	assert.Equal(t, 5, result.TotalItems)
	assert.Len(t, result.Invites, 5)
}

func TestListInvites_CombinedFilters(t *testing.T) {
	uc, inviteRepo := setupInviteListUC(t)
	invites := seedInvites()

	status := model.StatusPending
	expired := false
	after := time.Now().Add(-48 * time.Hour)
	inviteRepo.On("List", mock.Anything, mock.MatchedBy(func(f repository.InviteFilter) bool {
		return f.Status != nil && *f.Status == model.StatusPending &&
			f.Expired != nil && *f.Expired == false &&
			f.CreatedAfter != nil &&
			f.Page == 1 && f.PageSize == 20
	})).Return(&repository.InviteListResult{
		Invites:    []*model.InviteToken{invites[1], invites[3]},
		TotalItems: 2,
	}, nil)

	result, err := uc.ListInvites(context.Background(), repository.InviteFilter{
		Status:       &status,
		Expired:      &expired,
		CreatedAfter: &after,
		Page:         1,
		PageSize:     20,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalItems)
	assert.Len(t, result.Invites, 2)
}
