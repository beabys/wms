package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInviteRepositoryTypeCheck(t *testing.T) {
	var r *InviteRepository
	assert.Nil(t, r)
}

func TestRowToInvite_NilInvitedBy(t *testing.T) {
	invite := rowToInvite(&inviteRow{
		ID:        "inv-1",
		Email:     "test@example.com",
		Token:     "token-123",
		InvitedBy: nil,
		Status:    "pending",
	})
	assert.Equal(t, "inv-1", invite.ID)
	assert.Equal(t, "test@example.com", invite.Email)
	assert.Equal(t, "token-123", invite.Token)
	assert.Empty(t, invite.InvitedBy)
	assert.Equal(t, "pending", invite.Status)
}

func TestRowToInvite_WithInvitedBy(t *testing.T) {
	invitedBy := "admin-123"
	invite := rowToInvite(&inviteRow{
		ID:        "inv-2",
		Email:     "user@example.com",
		Token:     "token-456",
		InvitedBy: &invitedBy,
		Status:    "used",
	})
	assert.Equal(t, "inv-2", invite.ID)
	assert.Equal(t, "user@example.com", invite.Email)
	assert.Equal(t, "token-456", invite.Token)
	assert.Equal(t, "admin-123", invite.InvitedBy)
	assert.Equal(t, "used", invite.Status)
}

func TestRowToInvite_CancelledStatus(t *testing.T) {
	invitedBy := "admin-456"
	invite := rowToInvite(&inviteRow{
		ID:        "inv-3",
		Email:     "cancelled@example.com",
		Token:     "token-789",
		InvitedBy: &invitedBy,
		Status:    "cancelled",
	})
	assert.Equal(t, "cancelled", invite.Status)
}
