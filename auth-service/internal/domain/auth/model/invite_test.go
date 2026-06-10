package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInviteToken_Success(t *testing.T) {
	invite, err := NewInviteToken("user@example.com", "admin-123", 24*time.Hour)
	require.NoError(t, err)
	require.NotNil(t, invite)

	assert.Equal(t, "user@example.com", invite.Email)
	assert.Equal(t, "admin-123", invite.InvitedBy)
	assert.NotEmpty(t, invite.Token)
	assert.Equal(t, StatusPending, invite.Status)
	assert.False(t, invite.CreatedAt.IsZero())
	assert.False(t, invite.ExpiresAt.IsZero())
	assert.True(t, invite.ExpiresAt.After(invite.CreatedAt))
	assert.True(t, invite.IsValid())
	ok, reason := invite.CanCancel()
	assert.True(t, ok)
	assert.Empty(t, reason)
}

func TestNewInviteToken_EmptyEmail(t *testing.T) {
	_, err := NewInviteToken("", "admin-123", 24*time.Hour)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is required")
}

func TestNewInviteToken_InvalidEmail(t *testing.T) {
	_, err := NewInviteToken("not-an-email", "admin-123", 24*time.Hour)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email format")
}

func TestNewInviteToken_EmptyInvitedBy(t *testing.T) {
	_, err := NewInviteToken("user@example.com", "", 24*time.Hour)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invited_by is required")
}

func TestNewInviteToken_DefaultTTL(t *testing.T) {
	invite, err := NewInviteToken("user@example.com", "admin-123", 0)
	require.NoError(t, err)
	// Default TTL is 7 days
	assert.True(t, invite.ExpiresAt.After(time.Now().Add(6*24*time.Hour)))
	assert.True(t, invite.ExpiresAt.Before(time.Now().Add(8*24*time.Hour)))
}

func TestInviteToken_IsExpired(t *testing.T) {
	invite := &InviteToken{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Status:    StatusPending,
	}
	assert.True(t, invite.IsExpired())
	assert.False(t, invite.IsValid())
	ok, reason := invite.CanCancel()
	assert.False(t, ok)
	assert.Equal(t, "cannot cancel expired invite", reason)
}

func TestInviteToken_NotExpired(t *testing.T) {
	invite := &InviteToken{
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Status:    StatusPending,
	}
	assert.False(t, invite.IsExpired())
	assert.True(t, invite.IsValid())
	ok, reason := invite.CanCancel()
	assert.True(t, ok)
	assert.Empty(t, reason)
}

func TestInviteToken_UsedIsNotValid(t *testing.T) {
	invite := &InviteToken{
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Status:    StatusUsed,
	}
	assert.False(t, invite.IsValid())
	ok, reason := invite.CanCancel()
	assert.False(t, ok)
	assert.Contains(t, reason, "status: used")
}

func TestInviteToken_CancelledIsNotValid(t *testing.T) {
	invite := &InviteToken{
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Status:    StatusCancelled,
	}
	assert.False(t, invite.IsValid())
	ok, reason := invite.CanCancel()
	assert.False(t, ok)
	assert.Contains(t, reason, "status: cancelled")
}

func TestNewInviteToken_UniqueTokens(t *testing.T) {
	inv1, _ := NewInviteToken("a@example.com", "admin-1", 24*time.Hour)
	inv2, _ := NewInviteToken("b@example.com", "admin-2", 24*time.Hour)
	assert.NotEqual(t, inv1.Token, inv2.Token)
}
