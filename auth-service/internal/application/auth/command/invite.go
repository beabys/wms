package command

// InviteUserCommand represents the data needed to invite a user.
type InviteUserCommand struct {
	Email     string
	InvitedBy string
}

// InviteUserResult contains the result of an invite operation.
type InviteUserResult struct {
	Token      string
	InviteLink string
	ExpiresAt  int64
}

// ValidateInviteCommand represents the data needed to validate an invite.
type ValidateInviteCommand struct {
	Token string
}

// ValidateInviteResult contains the result of invite validation.
type ValidateInviteResult struct {
	Email     string
	InvitedBy string
}

// CancelInviteCommand represents the data needed to cancel an invite.
type CancelInviteCommand struct {
	Token string
}
