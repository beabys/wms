package model

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/mail"
	"time"
)

// Status constants for invite tokens.
const (
	StatusPending   = "pending"
	StatusUsed      = "used"
	StatusCancelled = "cancelled"
)

// InviteToken represents an invitation token for user registration.
type InviteToken struct {
	ID        string
	Email     string
	Token     string
	InvitedBy string
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// NewInviteToken creates a new InviteToken with validation.
// ttl defaults to 7 days if <= 0.
func NewInviteToken(email, invitedBy string, ttl time.Duration) (*InviteToken, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("invalid email format: %w", err)
	}
	if invitedBy == "" {
		return nil, fmt.Errorf("invited_by is required")
	}

	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour // 7 days default
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	now := time.Now()
	return &InviteToken{
		ID:        "",
		Email:     email,
		Token:     token,
		InvitedBy: invitedBy,
		Status:    StatusPending,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}, nil
}

// IsExpired checks if the invite token has expired.
func (i *InviteToken) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// IsValid checks if the invite token is still valid (pending and not expired).
func (i *InviteToken) IsValid() bool {
	return i.Status == StatusPending && !i.IsExpired()
}

// CanCancel checks if the invite token can be cancelled.
// Returns (true, "") if cancellable, otherwise (false, reason).
func (i *InviteToken) CanCancel() (bool, string) {
	if i.IsExpired() {
		return false, "cannot cancel expired invite"
	}
	if i.Status != StatusPending {
		return false, fmt.Sprintf("cannot cancel invite with status: %s", i.Status)
	}
	return true, ""
}
