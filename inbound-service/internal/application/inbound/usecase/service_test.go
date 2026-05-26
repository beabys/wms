package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/inbound-service/internal/application/inbound/command"
	"github.com/beabys/wms/inbound-service/internal/application/inbound/repository"
	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// mockRepo implements repository.InboundRepository for testing.
type mockRepo struct {
	inbounds  map[string]*model.Inbound
	saved     []*model.Inbound
	updated   []*model.Inbound
	getErr    error
	saveErr   error
	updateErr error
	listErr   error
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		inbounds: make(map[string]*model.Inbound),
	}
}

func (m *mockRepo) GetByID(_ context.Context, id string) (*model.Inbound, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	in, ok := m.inbounds[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return in, nil
}

func (m *mockRepo) Save(_ context.Context, inbound *model.Inbound) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.inbounds[inbound.ID] = inbound
	m.saved = append(m.saved, inbound)
	return nil
}

func (m *mockRepo) UpdateStatus(_ context.Context, inbound *model.Inbound) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.inbounds[inbound.ID] = inbound
	m.updated = append(m.updated, inbound)
	return nil
}

func (m *mockRepo) List(_ context.Context, customerID, status string, pageSize int32, pageToken string) ([]*model.Inbound, string, error) {
	if m.listErr != nil {
		return nil, "", m.listErr
	}
	var result []*model.Inbound
	for _, in := range m.inbounds {
		if customerID != "" && in.CustomerID != customerID {
			continue
		}
		if status != "" && string(in.Status) != status {
			continue
		}
		result = append(result, in)
	}
	return result, "", nil
}

var _ repository.InboundRepository = (*mockRepo)(nil)

// mockEventPublisher implements EventPublisher for testing.
type mockEventPublisher struct {
	submittedEvents int
	approvedEvents  int
	flaggedEvents   int
}

func (m *mockEventPublisher) PublishSubmitted(_ context.Context, _ *model.Inbound) error {
	m.submittedEvents++
	return nil
}

func (m *mockEventPublisher) PublishApproved(_ context.Context, _ *model.Inbound) error {
	m.approvedEvents++
	return nil
}

func (m *mockEventPublisher) PublishFlagged(_ context.Context, _ *model.Inbound) error {
	m.flaggedEvents++
	return nil
}

var _ EventPublisher = (*mockEventPublisher)(nil)

func setupTest(t *testing.T) (*InboundService, *mockRepo, *mockEventPublisher) {
	t.Helper()
	repo := newMockRepo()
	events := &mockEventPublisher{}
	h := NewInboundService(repo, events)
	return h, repo, events
}

func TestInboundService_SubmitInbound(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, repo, events := setupTest(t)

		result, err := h.SubmitInbound(context.Background(), command.SubmitInboundCommand{
			CustomerID:   "cust_1",
			ExpectedDate: "2026-06-01",
			Notes:        "test",
			Items: []command.SubmitItem{
				{SKU: "SKU001", QuantityDeclared: 10},
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "cust_1", result.Inbound.CustomerID)
		assert.Equal(t, model.StatusSubmitted, result.Inbound.Status)
		require.Len(t, repo.saved, 1)
		assert.Equal(t, 1, events.submittedEvents)
	})

	t.Run("empty items fails", func(t *testing.T) {
		h, _, _ := setupTest(t)

		_, err := h.SubmitInbound(context.Background(), command.SubmitInboundCommand{
			CustomerID:   "cust_1",
			ExpectedDate: "2026-06-01",
		})
		assert.Error(t, err)
	})
}

func TestInboundService_InspectInbound(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, repo, _ := setupTest(t)
		in := mustSaveInbound(t, repo)

		result, err := h.InspectInbound(context.Background(), command.InspectInboundCommand{
			InboundID:   in.ID,
			InspectorID: "insp_1",
			InspectedAt: time.Now(),
			Passed:      true,
		})

		require.NoError(t, err)
		assert.Equal(t, model.StatusInspected, result.Inbound.Status)
		assert.NotNil(t, result.Inbound.Inspection)
		assert.True(t, result.Inbound.Inspection.Passed)
	})

	t.Run("not found", func(t *testing.T) {
		h, _, _ := setupTest(t)
		_, err := h.InspectInbound(context.Background(), command.InspectInboundCommand{
			InboundID: "nonexistent",
		})
		assert.Error(t, err)
	})
}

func TestInboundService_ApproveInbound(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, repo, events := setupTest(t)
		in := mustSaveInbound(t, repo)
		// First inspect
		_ = in.Inspect(model.Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		_ = repo.UpdateStatus(context.Background(), in)

		result, err := h.ApproveInbound(context.Background(), command.ApproveInboundCommand{
			InboundID: in.ID,
		})

		require.NoError(t, err)
		assert.Equal(t, model.StatusApproved, result.Inbound.Status)
		assert.Equal(t, 1, events.approvedEvents)
	})

	t.Run("not inspected yet fails", func(t *testing.T) {
		h, repo, _ := setupTest(t)
		in := mustSaveInbound(t, repo)

		_, err := h.ApproveInbound(context.Background(), command.ApproveInboundCommand{
			InboundID: in.ID,
		})
		assert.Error(t, err)
	})
}

func TestInboundService_FlagInbound(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, repo, events := setupTest(t)
		in := mustSaveInbound(t, repo)
		_ = in.Inspect(model.Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		_ = repo.UpdateStatus(context.Background(), in)

		result, err := h.FlagInbound(context.Background(), command.FlagInboundCommand{
			InboundID: in.ID,
			Reason:    "damaged",
		})

		require.NoError(t, err)
		assert.Equal(t, model.StatusFlagged, result.Inbound.Status)
		assert.Equal(t, 1, events.flaggedEvents)
	})

	t.Run("submitted cannot be flagged", func(t *testing.T) {
		h, repo, _ := setupTest(t)
		in := mustSaveInbound(t, repo)

		_, err := h.FlagInbound(context.Background(), command.FlagInboundCommand{
			InboundID: in.ID,
			Reason:    "reason",
		})
		assert.Error(t, err)
	})
}

func TestInboundService_HoldAndRelease(t *testing.T) {
	t.Run("hold and release", func(t *testing.T) {
		h, repo, _ := setupTest(t)
		in := mustSaveInbound(t, repo)
		_ = in.Inspect(model.Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		_ = repo.UpdateStatus(context.Background(), in)

		// Hold
		holdResult, err := h.HoldInbound(context.Background(), command.HoldInboundCommand{
			InboundID: in.ID,
			Reason:    "docs missing",
		})
		require.NoError(t, err)
		assert.Equal(t, model.StatusHeld, holdResult.Inbound.Status)

		// Release
		relResult, err := h.ReleaseInbound(context.Background(), command.ReleaseInboundCommand{
			InboundID:  in.ID,
			ReleasedBy: "ops_user",
		})
		require.NoError(t, err)
		assert.Equal(t, model.StatusReleased, relResult.Inbound.Status)
	})
}

func TestInboundService_GetInbound(t *testing.T) {
	h, repo, _ := setupTest(t)
	in := mustSaveInbound(t, repo)

	result, err := h.GetInbound(context.Background(), command.GetInboundQuery{
		InboundID: in.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, in.ID, result.Inbound.ID)
}

func TestInboundService_ListInbounds(t *testing.T) {
	h, repo, _ := setupTest(t)
	mustSaveInbound(t, repo)
	mustSaveInboundCust(t, repo, "cust_2")

	result, err := h.ListInbounds(context.Background(), command.ListInboundsQuery{
		CustomerID: "cust_2",
	})
	require.NoError(t, err)
	assert.Len(t, result.Inbounds, 1)
}

// --- helpers ---

func mustSaveInbound(t *testing.T, repo *mockRepo) *model.Inbound {
	t.Helper()
	in, err := model.SubmitInbound(
		"inb_"+time.Now().Format("150405.000"),
		"cust_1", "2026-06-01", "", []model.InboundItem{
			{SKU: "SKU001", QuantityDeclared: 10},
		},
	)
	require.NoError(t, err)
	_ = repo.Save(context.Background(), in)
	return in
}

func mustSaveInboundCust(t *testing.T, repo *mockRepo, customerID string) *model.Inbound {
	t.Helper()
	in, err := model.SubmitInbound(
		"inb_"+time.Now().Format("150405.000"),
		customerID, "2026-06-01", "", []model.InboundItem{
			{SKU: "SKU001", QuantityDeclared: 10},
		},
	)
	require.NoError(t, err)
	_ = repo.Save(context.Background(), in)
	return in
}
