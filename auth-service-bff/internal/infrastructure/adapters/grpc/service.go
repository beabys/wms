package grpcdapter

import (
	"context"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
)

// AuthServiceAdapter wraps AuthClient and maps gRPC errors to domain errors.
// It is a thin layer: proto in, proto out, no business logic, no model mapping.
type AuthServiceAdapter struct {
	client *AuthClient
}

// NewAuthServiceAdapter creates a new adapter.
func NewAuthServiceAdapter(client *AuthClient) *AuthServiceAdapter {
	return &AuthServiceAdapter{client: client}
}

// Login authenticates a user and returns proto response.
func (a *AuthServiceAdapter) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	resp, err := a.client.Login(ctx, email, password)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// RefreshToken refreshes an access token.
func (a *AuthServiceAdapter) RefreshToken(ctx context.Context, refreshToken string) (*authv1.RefreshTokenResponse, error) {
	resp, err := a.client.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// ValidateToken validates a JWT and returns user claims.
func (a *AuthServiceAdapter) ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error) {
	resp, err := a.client.ValidateToken(ctx, token)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// Logout revokes a refresh token.
func (a *AuthServiceAdapter) Logout(ctx context.Context, refreshToken string) (*authv1.LogoutResponse, error) {
	resp, err := a.client.Logout(ctx, refreshToken)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// ListInvites retrieves paginated invite tokens with optional filters.
func (a *AuthServiceAdapter) ListInvites(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error) {
	resp, err := a.client.ListInvites(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// InviteUser creates an invite token for a new user.
func (a *AuthServiceAdapter) InviteUser(ctx context.Context, email, invitedBy string) (*authv1.InviteUserResponse, error) {
	resp, err := a.client.InviteUser(ctx, email, invitedBy)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// CancelInvite cancels a pending invite token.
func (a *AuthServiceAdapter) CancelInvite(ctx context.Context, token string) (*authv1.CancelInviteResponse, error) {
	resp, err := a.client.CancelInvite(ctx, token)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// CreateUser creates a new user.
func (a *AuthServiceAdapter) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	resp, err := a.client.CreateUser(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// GetUser retrieves a user by ID.
func (a *AuthServiceAdapter) GetUser(ctx context.Context, id string) (*authv1.GetUserResponse, error) {
	resp, err := a.client.GetUser(ctx, id)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// UpdateUser updates a user's fields.
func (a *AuthServiceAdapter) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	resp, err := a.client.UpdateUser(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// ListUsers lists users with optional filters.
func (a *AuthServiceAdapter) ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	resp, err := a.client.ListUsers(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// DeleteUser deactivates a user.
func (a *AuthServiceAdapter) DeleteUser(ctx context.Context, id string) (*authv1.DeleteUserResponse, error) {
	resp, err := a.client.DeleteUser(ctx, id)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// AssignRole assigns a role to a user.
func (a *AuthServiceAdapter) AssignRole(ctx context.Context, userID, role string) (*authv1.AssignRoleResponse, error) {
	resp, err := a.client.AssignRole(ctx, userID, role)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// ListRoles lists all roles.
func (a *AuthServiceAdapter) ListRoles(ctx context.Context, req *authv1.ListRolesRequest) (*authv1.ListRolesResponse, error) {
	resp, err := a.client.ListRoles(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// CreateRole creates a new role.
func (a *AuthServiceAdapter) CreateRole(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error) {
	resp, err := a.client.CreateRole(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}
