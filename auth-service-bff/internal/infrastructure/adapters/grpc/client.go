package grpcdapter

import (
	"context"
	"fmt"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/http/context"
)

// AuthClient is a gRPC client for the auth-service.
// It wraps AuthService and UserService RPCs.
type AuthClient struct {
	authService authv1.AuthServiceClient
	userService authv1.UserServiceClient
	conn        *grpc.ClientConn
}

// NewAuthClient creates a new gRPC connection and client.
func NewAuthClient(target string) (*AuthClient, error) {
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}
	return &AuthClient{
		authService: authv1.NewAuthServiceClient(conn),
		userService: authv1.NewUserServiceClient(conn),
		conn:        conn,
	}, nil
}

// Close closes the gRPC connection.
func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// forwardJWTCtx extracts a JWT token from the context and attaches it
// as gRPC outgoing metadata if present.
func (c *AuthClient) forwardJWTCtx(ctx context.Context) context.Context {
	token, _ := ctx.Value(httpctx.ContextKeyJWT).(string)
	if token != "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	}
	return ctx
}

// AuthService RPCs

// Login authenticates a user and returns tokens.
func (c *AuthClient) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.authService.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RefreshToken refreshes an access token using a refresh token.
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*authv1.RefreshTokenResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.authService.RefreshToken(ctx, &authv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ValidateToken validates a JWT and returns user claims.
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.authService.ValidateToken(ctx, &authv1.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Logout revokes a refresh token.
func (c *AuthClient) Logout(ctx context.Context, refreshToken string) (*authv1.LogoutResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.authService.Logout(ctx, &authv1.LogoutRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListInvites retrieves paginated invite tokens with optional filters.
func (c *AuthClient) ListInvites(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.authService.ListInvites(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// InviteUser creates an invite token for a new user.
func (c *AuthClient) InviteUser(ctx context.Context, email, invitedBy string) (*authv1.InviteUserResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.authService.InviteUser(ctx, &authv1.InviteUserRequest{
		Email:     email,
		InvitedBy: invitedBy,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CancelInvite cancels a pending invite token.
func (c *AuthClient) CancelInvite(ctx context.Context, token string) (*authv1.CancelInviteResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.authService.CancelInvite(ctx, &authv1.CancelInviteRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UserService RPCs

// CreateUser creates a new user.
func (c *AuthClient) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.userService.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUser retrieves a user by ID.
func (c *AuthClient) GetUser(ctx context.Context, id string) (*authv1.GetUserResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.userService.GetUser(ctx, &authv1.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateUser updates a user's fields.
func (c *AuthClient) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.userService.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListUsers lists users with optional filters and pagination.
func (c *AuthClient) ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.userService.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteUser deactivates a user.
func (c *AuthClient) DeleteUser(ctx context.Context, id string) (*authv1.DeleteUserResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.userService.DeleteUser(ctx, &authv1.DeleteUserRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AssignRole assigns a role to a user.
func (c *AuthClient) AssignRole(ctx context.Context, userID, role string) (*authv1.AssignRoleResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.userService.AssignRole(ctx, &authv1.AssignRoleRequest{
		UserId: userID,
		Role:   role,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListRoles lists all roles.
func (c *AuthClient) ListRoles(ctx context.Context, req *authv1.ListRolesRequest) (*authv1.ListRolesResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.userService.ListRoles(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateRole creates a new role.
func (c *AuthClient) CreateRole(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.userService.CreateRole(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
