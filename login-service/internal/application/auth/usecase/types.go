package usecase

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/beabys/wms/login-service/internal/application/auth/command"
	"github.com/beabys/wms/login-service/internal/application/auth/repository"
)

// AuthServiceHandler defines the usecase interface for auth operations.
type AuthServiceHandler interface {
	Login(ctx context.Context, cmd command.LoginCommand) (*command.LoginResponse, error)
	Register(ctx context.Context, cmd command.RegisterCommand) (*command.RegisterResponse, error)
	RefreshToken(ctx context.Context, cmd command.RefreshTokenCommand) (*command.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, tokenStr string) (*ValidateResult, error)
	CreateUser(ctx context.Context, cmd command.CreateUserCommand) (*command.CreateUserResponse, error)
	GetUser(ctx context.Context, query command.GetUserQuery) (*command.GetUserResponse, error)
	ListUsers(ctx context.Context, query command.ListUsersQuery) (*command.ListUsersResponse, error)
	UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) (*command.UpdateUserResponse, error)
}

// AuthService implements AuthServiceHandler.
type AuthService struct {
	repo          repository.UserRepository
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// ValidateResult contains the result of token validation.
type ValidateResult struct {
	UserID      string
	Role        string
	Permissions []string
}
