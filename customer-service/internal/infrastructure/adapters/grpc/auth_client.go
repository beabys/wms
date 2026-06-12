package grpcdapter

import (
	"context"
	"fmt"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/ports"
	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Ensure AuthClient implements ports.AuthClient interface.
var _ ports.AuthClient = (*AuthClient)(nil)

// AuthClient is a gRPC client for auth-service.
type AuthClient struct {
	userClient authv1.UserServiceClient
	authClient authv1.AuthServiceClient
	conn       *grpc.ClientConn
}

// NewAuthClient creates a new AuthClient.
func NewAuthClient(grpcAddr string) (*AuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth-service at %s: %w", grpcAddr, err)
	}

	return &AuthClient{
		userClient: authv1.NewUserServiceClient(conn),
		authClient: authv1.NewAuthServiceClient(conn),
		conn:       conn,
	}, nil
}

// CreateUser creates a user in auth-service.
func (c *AuthClient) CreateUser(ctx context.Context, req command.CreateUserRequest) (*command.CreateUserResponse, error) {
	resp, err := c.userClient.CreateUser(ctx, &authv1.CreateUserRequest{
		Email:      req.Email,
		Password:   req.Password,
		Name:       req.Name,
		Role:       req.Role,
		CustomerId: req.CustomerID,
		CreatedBy:  req.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("auth-service CreateUser failed: %w", err)
	}

	if resp.User == nil {
		return nil, fmt.Errorf("auth-service returned nil user")
	}

	return &command.CreateUserResponse{
		UserID: resp.User.GetId(),
		Email:  resp.User.GetEmail(),
		Name:   resp.User.GetName(),
		Role:   resp.User.GetRole(),
	}, nil
}

// ValidateToken validates a JWT via auth-service.
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*command.ValidateTokenResult, error) {
	resp, err := c.authClient.ValidateToken(ctx, &authv1.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, fmt.Errorf("auth-service ValidateToken failed: %w", err)
	}

	if !resp.GetValid() {
		return nil, fmt.Errorf("token is not valid")
	}

	return &command.ValidateTokenResult{
		UserID:      resp.GetUserId(),
		Role:        resp.GetRole(),
		CustomerID:  resp.GetCustomerId(),
		Permissions: resp.GetPermissions(),
	}, nil
}

// ValidateInvite validates an invite token via auth-service's ValidateInvite RPC.
func (c *AuthClient) ValidateInvite(ctx context.Context, token string) (string, string, error) {
	resp, err := c.authClient.ValidateInvite(ctx, &authv1.ValidateInviteRequest{
		Token: token,
	})
	if err != nil {
		return "", "", fmt.Errorf("auth-service ValidateInvite failed: %w", err)
	}

	if !resp.GetValid() {
		return "", "", fmt.Errorf("invite token is not valid")
	}

	return resp.GetEmail(), resp.GetInvitedBy(), nil
}

// Close closes the underlying gRPC connection.
func (c *AuthClient) Close() error {
	return c.conn.Close()
}
