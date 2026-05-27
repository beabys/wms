package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// --- mocks ---

type mockPool struct {
	beginFn    func(ctx context.Context) (pgx.Tx, error)
	queryRowFn func(ctx context.Context, sql string, args ...any) pgx.Row
	queryFn    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (m *mockPool) Begin(ctx context.Context) (pgx.Tx, error) {
	if m.beginFn != nil {
		return m.beginFn(ctx)
	}
	return &mockTx{}, nil
}

func (m *mockPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.queryRowFn != nil {
		return m.queryRowFn(ctx, sql, args...)
	}
	return &mockRow{}
}

func (m *mockPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.queryFn != nil {
		return m.queryFn(ctx, sql, args...)
	}
	return &mockRows{}, nil
}

type mockTx struct {
	beginFn    func(ctx context.Context) (pgx.Tx, error)
	commitFn   func(ctx context.Context) error
	rollbackFn func(ctx context.Context) error
	execFn     func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	queryRowFn func(ctx context.Context, sql string, args ...any) pgx.Row
	queryFn    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (m *mockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	if m.beginFn != nil {
		return m.beginFn(ctx)
	}
	return &mockTx{}, nil
}

func (m *mockTx) Commit(ctx context.Context) error {
	if m.commitFn != nil {
		return m.commitFn(ctx)
	}
	return nil
}

func (m *mockTx) Rollback(ctx context.Context) error {
	if m.rollbackFn != nil {
		return m.rollbackFn(ctx)
	}
	return nil
}

func (m *mockTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if m.execFn != nil {
		return m.execFn(ctx, sql, args...)
	}
	return pgconn.NewCommandTag("INSERT 0 1"), nil
}

func (m *mockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.queryRowFn != nil {
		return m.queryRowFn(ctx, sql, args...)
	}
	return &mockRow{}
}

func (m *mockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.queryFn != nil {
		return m.queryFn(ctx, sql, args...)
	}
	return &mockRows{}, nil
}

func (m *mockTx) CopyFrom(_ context.Context, _ pgx.Identifier, _ []string, _ pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (m *mockTx) SendBatch(_ context.Context, _ *pgx.Batch) pgx.BatchResults {
	return &mockBatchResults{}
}

func (m *mockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m *mockTx) Prepare(_ context.Context, _, _ string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (m *mockTx) Conn() *pgx.Conn {
	return nil
}

type mockRow struct {
	scanFn func(dest ...any) error
}

func (m *mockRow) Scan(dest ...any) error {
	if m.scanFn != nil {
		return m.scanFn(dest...)
	}
	return pgx.ErrNoRows
}

type mockRows struct {
	nextFn  func() bool
	scanFn  func(dest ...any) error
	closeFn func()
	errFn   func() error
	callNum int // number of times Next has been called
}

func (m *mockRows) Next() bool {
	if m.nextFn != nil {
		return m.nextFn()
	}
	return false
}

func (m *mockRows) Scan(dest ...any) error {
	if m.scanFn != nil {
		return m.scanFn(dest...)
	}
	return nil
}

func (m *mockRows) Close() {
	if m.closeFn != nil {
		m.closeFn()
	}
}

func (m *mockRows) Err() error {
	if m.errFn != nil {
		return m.errFn()
	}
	return nil
}

func (m *mockRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}

func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (m *mockRows) RawValues() [][]byte {
	return nil
}

func (m *mockRows) Values() ([]any, error) {
	return nil, nil
}

func (m *mockRows) Conn() *pgx.Conn {
	return nil
}

type mockBatchResults struct{}

func (m *mockBatchResults) Exec() (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag("INSERT 0 1"), nil
}

func (m *mockBatchResults) Query() (pgx.Rows, error) {
	return &mockRows{}, nil
}

func (m *mockBatchResults) QueryRow() pgx.Row {
	return &mockRow{}
}

func (m *mockBatchResults) Close() error {
	return nil
}

// --- helper to build repo with mock pool ---

func newRepoWithPool(pool PGPool) *InboundRepository {
	return &InboundRepository{pool: pool}
}

// --- tests ---

func TestSave(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var execCalls int
		tx := &mockTx{
			execFn: func(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				execCalls++
				return pgconn.NewCommandTag("INSERT 0 1"), nil
			},
		}
		pool := &mockPool{
			beginFn: func(_ context.Context) (pgx.Tx, error) { return tx, nil },
		}
		repo := newRepoWithPool(pool)

		in := &model.Inbound{
			ID: "inb_1", CustomerID: "cust_1", Status: model.StatusSubmitted,
			ExpectedDate: "2026-06-01", CreatedAt: time.Now(), UpdatedAt: time.Now(),
			Items: []model.InboundItem{
				{ID: "item_1", SKU: "SKU001", QuantityDeclared: 10},
			},
		}

		err := repo.Save(context.Background(), in)
		require.NoError(t, err)
		assert.Equal(t, 2, execCalls) // insert inbound + insert item
	})

	t.Run("begin tx fails", func(t *testing.T) {
		pool := &mockPool{
			beginFn: func(_ context.Context) (pgx.Tx, error) {
				return nil, errors.New("connection refused")
			},
		}
		repo := newRepoWithPool(pool)

		err := repo.Save(context.Background(), &model.Inbound{ID: "inb_1"})
		assert.ErrorContains(t, err, "begin tx")
	})

	t.Run("insert inbound fails", func(t *testing.T) {
		tx := &mockTx{
			execFn: func(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				return pgconn.CommandTag{}, errors.New("duplicate key")
			},
		}
		pool := &mockPool{
			beginFn: func(_ context.Context) (pgx.Tx, error) { return tx, nil },
		}
		repo := newRepoWithPool(pool)

		err := repo.Save(context.Background(), &model.Inbound{ID: "inb_1"})
		assert.ErrorContains(t, err, "insert inbound")
	})
}

func TestGetByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		now := time.Now()
		pool := &mockPool{
			queryRowFn: func(_ context.Context, sql string, args ...any) pgx.Row {
				// inbound main query
				if containsSQL(sql, "FROM inbounds") {
					return &mockRow{
						scanFn: func(dest ...any) error {
							*(dest[0].(*string)) = "inb_exists"
							*(dest[1].(*string)) = "cust_1"
							*(dest[2].(*model.Status)) = model.StatusSubmitted
							*(dest[3].(*string)) = "2026-06-01"
							// packaging_profile_id = nil (dest[4])
							// notes = "" (dest[5])
							*(dest[6].(*time.Time)) = now
							*(dest[7].(*time.Time)) = now
							return nil
						},
					}
				}
				// inspection / hold not found -> nil
				return &mockRow{scanFn: func(dest ...any) error { return pgx.ErrNoRows }}
			},
			queryFn: func(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
				return &mockRows{}, nil // empty items
			},
		}
		repo := newRepoWithPool(pool)

		in, err := repo.GetByID(context.Background(), "inb_exists")
		require.NoError(t, err)
		require.NotNil(t, in)
		assert.Equal(t, "inb_exists", in.ID)
		assert.Equal(t, "cust_1", in.CustomerID)
	})

	t.Run("not found", func(t *testing.T) {
		pool := &mockPool{
			queryRowFn: func(_ context.Context, _ string, _ ...any) pgx.Row {
				return &mockRow{scanFn: func(_ ...any) error { return pgx.ErrNoRows }}
			},
		}
		repo := newRepoWithPool(pool)

		in, err := repo.GetByID(context.Background(), "nonexistent")
		require.Error(t, err)
		assert.Nil(t, in)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("load items error", func(t *testing.T) {
		pool := &mockPool{
			queryRowFn: func(_ context.Context, _ string, _ ...any) pgx.Row {
				return &mockRow{
					scanFn: func(dest ...any) error {
						*(dest[0].(*string)) = "inb_1"
						*(dest[1].(*string)) = "cust_1"
						*(dest[2].(*model.Status)) = model.StatusSubmitted
						*(dest[3].(*string)) = "2026-06-01"
						*(dest[6].(*time.Time)) = time.Now()
						*(dest[7].(*time.Time)) = time.Now()
						return nil
					},
				}
			},
			queryFn: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
				return nil, errors.New("query failed")
			},
		}
		repo := newRepoWithPool(pool)

		in, err := repo.GetByID(context.Background(), "inb_1")
		require.Error(t, err)
		assert.Nil(t, in)
	})
}

func TestUpdateStatus(t *testing.T) {
	t.Run("status only", func(t *testing.T) {
		var execCalls int
		tx := &mockTx{
			execFn: func(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				execCalls++
				return pgconn.NewCommandTag("UPDATE 1"), nil
			},
		}
		pool := &mockPool{
			beginFn: func(_ context.Context) (pgx.Tx, error) { return tx, nil },
		}
		repo := newRepoWithPool(pool)

		in := &model.Inbound{
			ID: "inb_1", Status: model.StatusApproved, UpdatedAt: time.Now(),
		}

		err := repo.UpdateStatus(context.Background(), in)
		require.NoError(t, err)
		assert.Equal(t, 1, execCalls) // only update inbound
	})

	t.Run("with inspection and hold", func(t *testing.T) {
		var execCalls int
		tx := &mockTx{
			execFn: func(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				execCalls++
				return pgconn.NewCommandTag("INSERT 0 1"), nil
			},
		}
		pool := &mockPool{
			beginFn: func(_ context.Context) (pgx.Tx, error) { return tx, nil },
		}
		repo := newRepoWithPool(pool)

		now := time.Now()
		in := &model.Inbound{
			ID: "inb_1", Status: model.StatusInspected, UpdatedAt: now,
			Inspection: &model.Inspection{
				ID: "insp_1", InspectorID: "inspector_1", InspectedAt: now, Passed: true,
			},
			HoldRecord: &model.Hold{
				ID: "hold_1", Reason: "docs missing", CreatedAt: now,
			},
		}

		err := repo.UpdateStatus(context.Background(), in)
		require.NoError(t, err)
		assert.Equal(t, 3, execCalls) // update inbound + upsert inspection + upsert hold
	})

	t.Run("inspection ID passed correctly", func(t *testing.T) {
		var capturedArgs []any
		callCount := 0
		tx := &mockTx{
			execFn: func(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				callCount++
				if callCount == 2 { // second call = upsert inspection
					capturedArgs = args
				}
				return pgconn.NewCommandTag("INSERT 0 1"), nil
			},
		}
		pool := &mockPool{
			beginFn: func(_ context.Context) (pgx.Tx, error) { return tx, nil },
		}
		repo := newRepoWithPool(pool)

		now := time.Now()
		in := &model.Inbound{
			ID: "inb_1", Status: model.StatusInspected, UpdatedAt: now,
			Inspection: &model.Inspection{
				ID: "insp_fixed_id", InspectorID: "inspector_1", InspectedAt: now, Passed: true,
			},
		}

		err := repo.UpdateStatus(context.Background(), in)
		require.NoError(t, err)

		// upsert inspection args: inspectionID, inboundID, inspector_id, notes, photos_json, passed, inspected_at
		require.GreaterOrEqual(t, len(capturedArgs), 7)
		assert.Equal(t, "insp_fixed_id", capturedArgs[0])
	})

	t.Run("hold ID passed correctly", func(t *testing.T) {
		var capturedArgs []any
		callCount := 0
		tx := &mockTx{
			execFn: func(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
				callCount++
				if callCount == 3 { // third call = upsert hold
					capturedArgs = args
				}
				return pgconn.NewCommandTag("INSERT 0 1"), nil
			},
		}
		pool := &mockPool{
			beginFn: func(_ context.Context) (pgx.Tx, error) { return tx, nil },
		}
		repo := newRepoWithPool(pool)

		now := time.Now()
		in := &model.Inbound{
			ID: "inb_1", Status: model.StatusHeld, UpdatedAt: now,
			Inspection: &model.Inspection{
				ID: "insp_1", InspectorID: "inspector_1", InspectedAt: now, Passed: true,
			},
			HoldRecord: &model.Hold{
				ID: "hold_fixed_id", Reason: "docs", CreatedAt: now,
			},
		}

		err := repo.UpdateStatus(context.Background(), in)
		require.NoError(t, err)

		// upsert hold args: holdID, inboundID, reason, created_at, released_at, released_by
		require.GreaterOrEqual(t, len(capturedArgs), 6)
		assert.Equal(t, "hold_fixed_id", capturedArgs[0])
	})

	t.Run("begin tx fails", func(t *testing.T) {
		pool := &mockPool{
			beginFn: func(_ context.Context) (pgx.Tx, error) {
				return nil, errors.New("timeout")
			},
		}
		repo := newRepoWithPool(pool)

		err := repo.UpdateStatus(context.Background(), &model.Inbound{ID: "inb_1"})
		assert.ErrorContains(t, err, "begin tx")
	})
}

func TestList(t *testing.T) {
	t.Run("no filters", func(t *testing.T) {
		callNum := 0
		now := time.Now()
		rows := &mockRows{
			nextFn: func() bool {
				callNum++
				return callNum <= 2
			},
			scanFn: func(dest ...any) error {
				*(dest[0].(*string)) = "inb_1"
				*(dest[1].(*string)) = "cust_1"
				*(dest[2].(*model.Status)) = model.StatusSubmitted
				*(dest[3].(*string)) = "2026-06-01"
				*(dest[5].(*string)) = ""
				*(dest[6].(*time.Time)) = now
				*(dest[7].(*time.Time)) = now
				return nil
			},
		}
		pool := &mockPool{
			queryFn: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
				return rows, nil
			},
		}
		repo := newRepoWithPool(pool)

		inbounds, token, err := repo.List(context.Background(), "", "", 20, "")
		require.NoError(t, err)
		assert.Len(t, inbounds, 2)
		assert.Empty(t, token)
	})

	t.Run("with customer filter", func(t *testing.T) {
		var capturedArgs []any
		callNum := 0
		rows := &mockRows{
			nextFn: func() bool {
				callNum++
				return callNum <= 1
			},
			scanFn: func(dest ...any) error {
				*(dest[0].(*string)) = "inb_1"
				*(dest[1].(*string)) = "cust_1"
				*(dest[2].(*model.Status)) = model.StatusSubmitted
				*(dest[3].(*string)) = "2026-06-01"
				*(dest[5].(*string)) = ""
				*(dest[6].(*time.Time)) = time.Now()
				*(dest[7].(*time.Time)) = time.Now()
				return nil
			},
		}
		pool := &mockPool{
			queryFn: func(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
				capturedArgs = args
				return rows, nil
			},
		}
		repo := newRepoWithPool(pool)

		inbounds, _, err := repo.List(context.Background(), "cust_1", "", 20, "")
		require.NoError(t, err)
		require.Len(t, inbounds, 1)
		require.GreaterOrEqual(t, len(capturedArgs), 1)
		assert.Equal(t, "cust_1", capturedArgs[0])
	})

	t.Run("pagination", func(t *testing.T) {
		callNum := 0
		now := time.Now()
		rows := &mockRows{
			nextFn: func() bool {
				callNum++
				return callNum <= 2
			},
			scanFn: func(dest ...any) error {
				if callNum == 1 {
					*(dest[0].(*string)) = "inb_1"
				} else {
					*(dest[0].(*string)) = "inb_2"
				}
				*(dest[1].(*string)) = "cust_1"
				*(dest[2].(*model.Status)) = model.StatusSubmitted
				*(dest[3].(*string)) = "2026-06-01"
				*(dest[5].(*string)) = ""
				*(dest[6].(*time.Time)) = now
				*(dest[7].(*time.Time)) = now
				return nil
			},
		}
		pool := &mockPool{
			queryFn: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
				return rows, nil
			},
		}
		repo := newRepoWithPool(pool)

		inbounds, token, err := repo.List(context.Background(), "", "", 1, "")
		require.NoError(t, err)
		require.Len(t, inbounds, 1)
		assert.Equal(t, "inb_1", token)
	})

	t.Run("query error", func(t *testing.T) {
		pool := &mockPool{
			queryFn: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
				return nil, errors.New("permission denied")
			},
		}
		repo := newRepoWithPool(pool)

		_, _, err := repo.List(context.Background(), "", "", 20, "")
		assert.ErrorContains(t, err, "list inbounds")
	})
}

func TestLoadInspection(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		now := time.Now()
		photos := []string{"photo1.jpg"}
		photosJSON, _ := json.Marshal(photos)

		pool := &mockPool{
			queryRowFn: func(_ context.Context, _ string, _ ...any) pgx.Row {
				return &mockRow{
					scanFn: func(dest ...any) error {
						*(dest[0].(*string)) = "insp_1"
						*(dest[1].(*string)) = "inspector_1"
						*(dest[2].(*string)) = "all good"
						*(dest[3].(*string)) = string(photosJSON)
						*(dest[4].(*bool)) = true
						*(dest[5].(*time.Time)) = now
						return nil
					},
				}
			},
		}
		repo := newRepoWithPool(pool)

		insp, err := repo.loadInspection(context.Background(), "inb_1")
		require.NoError(t, err)
		require.NotNil(t, insp)
		assert.Equal(t, "insp_1", insp.ID)
		assert.Equal(t, "inspector_1", insp.InspectorID)
		assert.True(t, insp.Passed)
	})

	t.Run("not found returns nil", func(t *testing.T) {
		pool := &mockPool{
			queryRowFn: func(_ context.Context, _ string, _ ...any) pgx.Row {
				return &mockRow{scanFn: func(_ ...any) error { return pgx.ErrNoRows }}
			},
		}
		repo := newRepoWithPool(pool)

		insp, err := repo.loadInspection(context.Background(), "inb_none")
		require.NoError(t, err)
		assert.Nil(t, insp)
	})
}

func TestLoadHold(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		now := time.Now()

		pool := &mockPool{
			queryRowFn: func(_ context.Context, _ string, _ ...any) pgx.Row {
				return &mockRow{
					scanFn: func(dest ...any) error {
						*(dest[0].(*string)) = "hold_1"
						*(dest[1].(*string)) = "docs missing"
						*(dest[2].(*time.Time)) = now
						*(dest[3].(**time.Time)) = nil
						*(dest[4].(*string)) = ""
						return nil
					},
				}
			},
		}
		repo := newRepoWithPool(pool)

		hold, err := repo.loadHold(context.Background(), "inb_1")
		require.NoError(t, err)
		require.NotNil(t, hold)
		assert.Equal(t, "hold_1", hold.ID)
		assert.Equal(t, "docs missing", hold.Reason)
	})

	t.Run("not found returns nil", func(t *testing.T) {
		pool := &mockPool{
			queryRowFn: func(_ context.Context, _ string, _ ...any) pgx.Row {
				return &mockRow{scanFn: func(_ ...any) error { return pgx.ErrNoRows }}
			},
		}
		repo := newRepoWithPool(pool)

		hold, err := repo.loadHold(context.Background(), "inb_none")
		require.NoError(t, err)
		assert.Nil(t, hold)
	})
}

func TestNullString(t *testing.T) {
	assert.Nil(t, nullString(""))
	s := nullString("value")
	require.NotNil(t, s)
	assert.Equal(t, "value", *s)
}

// containsSQL is a test helper to check if sql contains substr.
func containsSQL(sql, substr string) bool {
	return strings.Contains(sql, substr)
}
