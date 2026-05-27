package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

func TestNewAddress_Valid(t *testing.T) {
	t.Parallel()

	addr, err := model.NewAddress("Main St 1", "", "Berlin", "10115", "DE")
	require.NoError(t, err)
	assert.Equal(t, "Main St 1", addr.Line1)
	assert.Equal(t, "Berlin", addr.City)
	assert.Equal(t, "DE", addr.Country)
}

func TestNewAddress_MissingLine1(t *testing.T) {
	t.Parallel()

	_, err := model.NewAddress("", "", "Berlin", "10115", "DE")
	assert.ErrorContains(t, err, "line1 is required")
}

func TestNewAddress_MissingCity(t *testing.T) {
	t.Parallel()

	_, err := model.NewAddress("Main St 1", "", "", "10115", "DE")
	assert.ErrorContains(t, err, "city is required")
}

func TestNewAddress_MissingCountry(t *testing.T) {
	t.Parallel()

	_, err := model.NewAddress("Main St 1", "", "Berlin", "10115", "")
	assert.ErrorContains(t, err, "country is required")
}

func TestVATNumber_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
	}{
		{"DE123456789"},
		{"DE000000000"},
		{"DE999999999"},
	}

	for _, tt := range tests {
		vat, err := model.NewVATNumber(tt.input)
		require.NoError(t, err, tt.input)
		assert.Equal(t, tt.input, vat.String())
	}
}

func TestVATNumber_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		desc  string
	}{
		{"", "empty"},
		{"DE12345678", "too short"},
		{"DE1234567890", "too long"},
		{"123456789", "missing DE prefix"},
		{"FR123456789", "wrong prefix"},
		{"DEABCDEFGHI", "non-digits"},
		{"De123456789", "lowercase d"},
	}

	for _, tt := range tests {
		_, err := model.NewVATNumber(tt.input)
		assert.Error(t, err, tt.desc)
	}
}

func TestCustomerStatusTransitions(t *testing.T) {
	t.Parallel()

	t.Run("approve from pending succeeds", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusPending}
		err := c.Approve()
		require.NoError(t, err)
		assert.Equal(t, model.CustomerStatusActive, c.Status)
	})

	t.Run("approve from active fails", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusActive}
		err := c.Approve()
		assert.ErrorContains(t, err, "cannot approve")
	})

	t.Run("approve from suspended fails", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusSuspended}
		err := c.Approve()
		assert.ErrorContains(t, err, "cannot approve")
	})

	t.Run("suspend from active succeeds", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusActive}
		err := c.Suspend()
		require.NoError(t, err)
		assert.Equal(t, model.CustomerStatusSuspended, c.Status)
	})

	t.Run("suspend from pending fails", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusPending}
		err := c.Suspend()
		assert.ErrorContains(t, err, "cannot suspend")
	})

	t.Run("suspend from suspended fails", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusSuspended}
		err := c.Suspend()
		assert.ErrorContains(t, err, "cannot suspend")
	})

	t.Run("reactivate from suspended succeeds", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusSuspended}
		err := c.Reactivate()
		require.NoError(t, err)
		assert.Equal(t, model.CustomerStatusActive, c.Status)
	})

	t.Run("reactivate from pending fails", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusPending}
		err := c.Reactivate()
		assert.ErrorContains(t, err, "cannot reactivate")
	})

	t.Run("reactivate from active fails", func(t *testing.T) {
		c := &model.Customer{Status: model.CustomerStatusActive}
		err := c.Reactivate()
		assert.ErrorContains(t, err, "cannot reactivate")
	})
}

func TestInviteLink_MarkUsed(t *testing.T) {
	t.Parallel()

	t.Run("mark unused invite succeeds", func(t *testing.T) {
		invite := &model.InviteLink{
			ID:        "1",
			Email:     "test@example.com",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		err := invite.MarkUsed()
		require.NoError(t, err)
		assert.True(t, invite.Used)
	})

	t.Run("mark already used fails", func(t *testing.T) {
		invite := &model.InviteLink{
			ID:    "2",
			Used:  true,
			Email: "test@example.com",
		}
		err := invite.MarkUsed()
		assert.ErrorContains(t, err, "already used")
	})

	t.Run("mark expired invite fails", func(t *testing.T) {
		invite := &model.InviteLink{
			ID:        "3",
			Email:     "test@example.com",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		err := invite.MarkUsed()
		assert.ErrorContains(t, err, "has expired")
	})
}

func TestInviteLink_IsExpired(t *testing.T) {
	t.Parallel()

	t.Run("expired returns true", func(t *testing.T) {
		invite := &model.InviteLink{
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		assert.True(t, invite.IsExpired())
	})

	t.Run("not expired returns false", func(t *testing.T) {
		invite := &model.InviteLink{
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		assert.False(t, invite.IsExpired())
	})
}

func TestNewID_GeneratesUniqueIDs(t *testing.T) {
	t.Parallel()

	id1 := model.NewID()
	id2 := model.NewID()
	assert.NotEqual(t, id1, id2)
	assert.Len(t, id1, 36) // UUID format
}
