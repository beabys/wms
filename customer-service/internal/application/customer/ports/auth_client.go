package ports

import (
	"context"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
)

// AuthClient defines the methods the use case needs from auth-service.
type AuthClient interface {
	CreateUser(ctx context.Context, req command.CreateUserRequest) (*command.CreateUserResponse, error)
	ValidateToken(ctx context.Context, token string) (*command.ValidateTokenResult, error)
	ValidateInvite(ctx context.Context, token string) (email string, invitedBy string, err error)
}
