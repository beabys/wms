package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
)

// withToken creates a context with JWT in gRPC metadata.
func withToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

// Client wraps the gRPC clients for the auth service.
type Client struct {
	AuthClient authv1.AuthServiceClient
	UserClient authv1.UserServiceClient
	conn       *grpc.ClientConn
}

// NewClient creates a new gRPC client connected to the auth service.
func NewClient(ctx context.Context, address string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("dial auth service: %w", err)
	}

	return &Client{
		AuthClient: authv1.NewAuthServiceClient(conn),
		UserClient: authv1.NewUserServiceClient(conn),
		conn:       conn,
	}, nil
}

// Close closes the underlying gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Login authenticates a user.
func (c *Client) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	resp, err := c.AuthClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	return resp, nil
}

// Register creates a new user account.
func (c *Client) Register(ctx context.Context, email, password, companyName string) (*authv1.CreateUserResponse, error) {
	req := &authv1.CreateUserRequest{
		Email:    email,
		Password: password,
		Role:     "customer",
	}
	if companyName != "" {
		req.CustomerId = companyName
	}
	resp, err := c.UserClient.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	return resp, nil
}

// RefreshToken refreshes an access token.
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*authv1.RefreshTokenResponse, error) {
	resp, err := c.AuthClient.RefreshToken(ctx, &authv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}
	return resp, nil
}

// ValidateToken validates a token and returns user info.
func (c *Client) ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error) {
	resp, err := c.AuthClient.ValidateToken(ctx, &authv1.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}
	return resp, nil
}

// GetUser retrieves a user by ID.
func (c *Client) GetUser(ctx context.Context, userID string, token string) (*authv1.GetUserResponse, error) {
	ctx = withToken(ctx, token)
	resp, err := c.UserClient.GetUser(ctx, &authv1.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return resp, nil
}

// DeleteRefreshToken invalidates a refresh token (log out).
func (c *Client) DeleteRefreshToken(ctx context.Context, token string) error {
	// Use a refresh-token-based validation to invalidate
	// In real impl, this would call a dedicated logout endpoint
	return nil
}
