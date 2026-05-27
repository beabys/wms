package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"

	"github.com/beabys/wms/login-service/internal/application/auth/command"
	"github.com/beabys/wms/login-service/internal/application/auth/usecase"
)

// AuthServer implements the auth.v1.AuthServiceServer gRPC interface.
type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	authService  *usecase.AuthService
	logger       *zap.Logger
	publicKeyPem string
}

// NewAuthServer creates a new AuthServer.
func NewAuthServer(authService *usecase.AuthService, logger *zap.Logger, publicKeyPem string) *AuthServer {
	return &AuthServer{
		authService:  authService,
		logger:       logger,
		publicKeyPem: publicKeyPem,
	}
}

// Login handles login requests.
func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	result, err := s.authService.Login(ctx, command.LoginCommand{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		s.logger.Warn("login failed", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}

	return &authv1.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &authv1.User{
			Id:         result.UserID,
			Email:      result.UserEmail,
			Role:       result.UserRole,
			CustomerId: result.CustomerID,
		},
	}, nil
}

// RefreshToken handles token refresh requests.
func (s *AuthServer) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	result, err := s.authService.RefreshToken(ctx, command.RefreshTokenCommand{
		RefreshToken: req.GetRefreshToken(),
	})
	if err != nil {
		s.logger.Warn("refresh token failed", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "refresh failed: %v", err)
	}

	return &authv1.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

// ValidateToken handles token validation requests.
func (s *AuthServer) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	result, err := s.authService.ValidateToken(ctx, req.GetToken())
	if err != nil {
		s.logger.Warn("validate token failed", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "validate failed: %v", err)
	}

	return &authv1.ValidateTokenResponse{
		Valid:       true,
		UserId:      result.UserID,
		Role:        result.Role,
		Permissions: result.Permissions,
	}, nil
}

// GetPublicKey returns the public key in PEM format.
func (s *AuthServer) GetPublicKey(_ context.Context, _ *authv1.GetPublicKeyRequest) (*authv1.GetPublicKeyResponse, error) {
	return &authv1.GetPublicKeyResponse{
		PublicKeyPem: s.publicKeyPem,
		KeyId:        "rsa-1",
	}, nil
}

// UserServer implements the auth.v1.UserServiceServer gRPC interface.
type UserServer struct {
	authv1.UnimplementedUserServiceServer
	authService *usecase.AuthService
	logger      *zap.Logger
}

// NewUserServer creates a new UserServer.
func NewUserServer(authService *usecase.AuthService, logger *zap.Logger) *UserServer {
	return &UserServer{
		authService: authService,
		logger:      logger,
	}
}

// CreateUser handles user creation.
func (s *UserServer) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	result, err := s.authService.CreateUser(ctx, command.CreateUserCommand{
		Email:      req.GetEmail(),
		Password:   req.GetPassword(),
		Role:       req.GetRole(),
		CustomerID: req.GetCustomerId(),
	})
	if err != nil {
		s.logger.Warn("create user failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "create user failed: %v", err)
	}

	return &authv1.CreateUserResponse{
		User: &authv1.User{
			Id:    result.UserID,
			Email: result.Email,
			Role:  result.Role,
		},
	}, nil
}

// GetUser handles user retrieval.
func (s *UserServer) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	result, err := s.authService.GetUser(ctx, command.GetUserQuery{UserID: req.GetId()})
	if err != nil {
		s.logger.Warn("get user failed", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &authv1.GetUserResponse{
		User: &authv1.User{
			Id:         result.UserID,
			Email:      result.Email,
			Role:       result.Role,
			CustomerId: result.CustomerID,
			CreatedAt:  result.CreatedAt,
		},
	}, nil
}

// ListUsers handles user listing.
func (s *UserServer) ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	page := int32(1)
	pageSize := int32(20)
	if req.GetPagination() != nil {
		page = req.GetPagination().Page
		if req.GetPagination().Limit > 0 {
			pageSize = req.GetPagination().Limit
		}
	}

	result, err := s.authService.ListUsers(ctx, command.ListUsersQuery{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		s.logger.Warn("list users failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "list users failed: %v", err)
	}

	var protoUsers []*authv1.User
	for _, u := range result.Users {
		protoUsers = append(protoUsers, &authv1.User{
			Id:         u.UserID,
			Email:      u.Email,
			Role:       u.Role,
			CustomerId: u.CustomerID,
			CreatedAt:  u.CreatedAt,
		})
	}

	return &authv1.ListUsersResponse{
		Users: protoUsers,
	}, nil
}

// UpdateUser handles user updates.
func (s *UserServer) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	result, err := s.authService.UpdateUser(ctx, command.UpdateUserCommand{
		UserID:     req.GetId(),
		Email:      req.GetEmail(),
		Role:       req.GetRole(),
		CustomerID: req.GetCustomerId(),
	})
	if err != nil {
		s.logger.Warn("update user failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "update user failed: %v", err)
	}

	return &authv1.UpdateUserResponse{
		User: &authv1.User{
			Id:         result.UserID,
			Email:      result.Email,
			Role:       result.Role,
			CustomerId: result.CustomerID,
		},
	}, nil
}

// Ensure compile-time interface satisfaction.
var _ authv1.AuthServiceServer = (*AuthServer)(nil)
var _ authv1.UserServiceServer = (*UserServer)(nil)
