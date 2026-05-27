package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitInbound(t *testing.T) {
	t.Run("valid submission", func(t *testing.T) {
		in, err := SubmitInbound("inb_1", "cust_1", "2026-06-01", "test notes", []InboundItem{
			{SKU: "SKU001", QuantityDeclared: 10},
		})
		require.NoError(t, err)
		assert.Equal(t, "inb_1", in.ID)
		assert.Equal(t, "cust_1", in.CustomerID)
		assert.Equal(t, StatusSubmitted, in.Status)
		assert.Len(t, in.Items, 1)
	})

	t.Run("empty id", func(t *testing.T) {
		_, err := SubmitInbound("", "cust_1", "2026-06-01", "", []InboundItem{
			{SKU: "SKU001", QuantityDeclared: 10},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "id is required")
	})

	t.Run("empty customer", func(t *testing.T) {
		_, err := SubmitInbound("inb_1", "", "2026-06-01", "", []InboundItem{
			{SKU: "SKU001", QuantityDeclared: 10},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "customer id")
	})

	t.Run("empty expected date", func(t *testing.T) {
		_, err := SubmitInbound("inb_1", "cust_1", "", "", []InboundItem{
			{SKU: "SKU001", QuantityDeclared: 10},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected date")
	})

	t.Run("no items", func(t *testing.T) {
		_, err := SubmitInbound("inb_1", "cust_1", "2026-06-01", "", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one item")
	})

	t.Run("invalid item - empty SKU", func(t *testing.T) {
		_, err := SubmitInbound("inb_1", "cust_1", "2026-06-01", "", []InboundItem{
			{SKU: "", QuantityDeclared: 10},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sku is required")
	})

	t.Run("invalid item - zero quantity", func(t *testing.T) {
		_, err := SubmitInbound("inb_1", "cust_1", "2026-06-01", "", []InboundItem{
			{SKU: "SKU001", QuantityDeclared: 0},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be positive")
	})
}

func TestInbound_Inspect(t *testing.T) {
	validInspection := Inspection{
		InspectorID: "inspector_1",
		InspectedAt: time.Now(),
		Passed:      true,
		Notes:       "all good",
	}

	t.Run("inspect submitted inbound", func(t *testing.T) {
		in := mustSubmit(t)
		err := in.Inspect(validInspection)
		require.NoError(t, err)
		assert.Equal(t, StatusInspected, in.Status)
		assert.NotNil(t, in.Inspection)
		assert.True(t, in.Inspection.Passed)
	})

	t.Run("inspect with empty inspector id", func(t *testing.T) {
		in := mustSubmit(t)
		err := in.Inspect(Inspection{InspectedAt: time.Now()})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "inspector_id")
	})

	t.Run("inspect already inspected inbound", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(validInspection)
		err := in.Inspect(validInspection)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status transition")
	})

	t.Run("inspect approved inbound", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(validInspection)
		_ = in.Approve()
		err := in.Inspect(validInspection)
		assert.Error(t, err)
	})
}

func TestInbound_Approve(t *testing.T) {
	t.Run("approve inspected inbound", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		err := in.Approve()
		require.NoError(t, err)
		assert.Equal(t, StatusApproved, in.Status)
	})

	t.Run("approve without inspection", func(t *testing.T) {
		in := mustSubmit(t)
		err := in.Approve()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status transition")
	})

	t.Run("approve already approved", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		_ = in.Approve()
		err := in.Approve()
		assert.Error(t, err)
	})
}

func TestInbound_Flag(t *testing.T) {
	t.Run("flag inspected inbound", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		err := in.Flag("damaged goods")
		require.NoError(t, err)
		assert.Equal(t, StatusFlagged, in.Status)
	})

	t.Run("flag without reason", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		err := in.Flag("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reason is required")
	})

	t.Run("flag submitted inbound (invalid)", func(t *testing.T) {
		in := mustSubmit(t)
		err := in.Flag("reason")
		assert.Error(t, err)
	})
}

func TestInbound_HoldAndRelease(t *testing.T) {
	t.Run("hold inspected inbound", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		err := in.PlaceHold("awaiting docs")
			require.NoError(t, err)
			assert.Equal(t, StatusHeld, in.Status)
			assert.NotNil(t, in.HoldRecord)
			assert.Equal(t, "awaiting docs", in.HoldRecord.Reason)
	})

	t.Run("hold without reason", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		err := in.PlaceHold("")
			assert.Error(t, err)
	})

	t.Run("release held inbound", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		_ = in.PlaceHold("docs missing")
			err := in.Release("ops_user")
			require.NoError(t, err)
			assert.Equal(t, StatusReleased, in.Status)
			assert.NotNil(t, in.HoldRecord.ReleasedAt)
			assert.Equal(t, "ops_user", in.HoldRecord.ReleasedBy)
	})

	t.Run("release without released_by", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		_ = in.PlaceHold("docs missing")
		err := in.Release("")
		assert.Error(t, err)
	})

	t.Run("release without hold", func(t *testing.T) {
		in := mustSubmit(t)
		_ = in.Inspect(Inspection{InspectorID: "insp_1", InspectedAt: time.Now()})
		// Don't hold, try to release
		err := in.Release("ops_user")
		assert.Error(t, err)
	})

	t.Run("release non-held inbound", func(t *testing.T) {
		in := mustSubmit(t)
		err := in.Release("ops_user")
		assert.Error(t, err)
	})
}

func TestInbound_ItemValidation(t *testing.T) {
	t.Run("negative quantity_received", func(t *testing.T) {
		err := InboundItem{SKU: "SKU001", QuantityDeclared: 10, QuantityReceived: -1}.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")
	})

	t.Run("valid item", func(t *testing.T) {
		err := InboundItem{SKU: "SKU001", QuantityDeclared: 10, QuantityReceived: 5}.Validate()
		assert.NoError(t, err)
	})
}

// mustSubmit is a test helper.
func mustSubmit(t *testing.T) *Inbound {
	t.Helper()
	in, err := SubmitInbound("inb_test", "cust_test", "2026-06-01", "", []InboundItem{
		{SKU: "SKU001", QuantityDeclared: 10},
	})
	require.NoError(t, err)
	return in
}
