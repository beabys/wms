package grpcdapter

import (
	"context"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/application/auth/transformer"
	"github.com/beabys/wms/auth-service/internal/application/auth/usecase"
	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer implements authv1.AuthServiceServer.
type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	loginUC  *usecase.LoginUseCase
	inviteUC *usecase.InviteUseCase
}

// NewAuthServer creates a new AuthServer.
func NewAuthServer(loginUC *usecase.LoginUseCase, inviteUC *usecase.InviteUseCase) *AuthServer {
	return &AuthServer{
		loginUC:  loginUC,
		inviteUC: inviteUC,
	}
}

// Login authenticates a user and returns JWT tokens.
func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	result, err := s.loginUC.Login(ctx, command.LoginCommand{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	return &authv1.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// RefreshToken validates a refresh token and issues new tokens.
func (s *AuthServer) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	result, err := s.loginUC.RefreshToken(ctx, command.RefreshTokenCommand{
		RefreshToken: req.GetRefreshToken(),
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	return &authv1.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// ValidateToken validates a JWT and returns claims.
func (s *AuthServer) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	result, err := s.loginUC.ValidateToken(ctx, command.ValidateTokenCommand{
		Token: req.GetToken(),
	})
	if err != nil {
		return &authv1.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &authv1.ValidateTokenResponse{
		Valid:       true,
		UserId:      result.UserID,
		Role:        result.Role,
		CustomerId:  result.CustomerID,
		Permissions: result.Permissions,
	}, nil
}

// InviteUser creates an invite token for a new user.
func (s *AuthServer) InviteUser(ctx context.Context, req *authv1.InviteUserRequest) (*authv1.InviteUserResponse, error) {
	// Extract invited_by from JWT claims (set by auth interceptor)
	invitedBy := req.GetInvitedBy()
	if claims := ClaimsFromContext(ctx); claims != nil && claims.UserID != "" {
		invitedBy = claims.UserID
	}

	result, err := s.inviteUC.InviteUser(ctx, command.InviteUserCommand{
		Email:     req.GetEmail(),
		InvitedBy: invitedBy,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authv1.InviteUserResponse{
		Token:     result.Token,
		InviteLink: result.InviteLink,
		ExpiresAt: result.ExpiresAt,
	}, nil
}

// Logout revokes a refresh token.
func (s *AuthServer) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if err := s.loginUC.Logout(ctx, command.LogoutCommand{
		RefreshToken: req.GetRefreshToken(),
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.LogoutResponse{Success: true}, nil
}

// ValidateInvite validates an invite token via the use case.
func (s *AuthServer) ValidateInvite(ctx context.Context, req *authv1.ValidateInviteRequest) (*authv1.ValidateInviteResponse, error) {
	result, err := s.inviteUC.ValidateInvite(ctx, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &authv1.ValidateInviteResponse{
		Valid:     true,
		Email:     result.Email,
		InvitedBy: result.InvitedBy,
	}, nil
}

// ListInvites retrieves paginated invite tokens with optional filters.
func (s *AuthServer) ListInvites(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error) {
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	filter := repository.InviteFilter{
		Page:     int(req.GetPage()),
		PageSize: pageSize,
	}

	if req.Status != nil {
		filter.Status = req.Status
	}
	if req.Expired != nil {
		filter.Expired = req.Expired
	}
	if req.CreatedAfter != nil {
		t := time.Unix(*req.CreatedAfter, 0)
		filter.CreatedAfter = &t
	}
	if req.CreatedBefore != nil {
		t := time.Unix(*req.CreatedBefore, 0)
		filter.CreatedBefore = &t
	}

	result, err := s.inviteUC.ListInvites(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	invites := make([]*authv1.InviteEntry, 0, len(result.Invites))
	for _, inv := range result.Invites {
		invites = append(invites, transformer.InviteToProto(inv))
	}

	totalPages := 0
	if filter.PageSize > 0 {
		totalPages = (result.TotalItems + filter.PageSize - 1) / filter.PageSize
	}

	return &authv1.ListInvitesResponse{
		Invites: invites,
		Pagination: &commonv1.Pagination{
			Page:       int32(filter.Page),
			PageSize:   int32(filter.PageSize),
			Total:      int32(result.TotalItems),
			TotalPages: int32(totalPages),
		},
	}, nil
}

// CancelInvite cancels a pending invite token.
func (s *AuthServer) CancelInvite(ctx context.Context, req *authv1.CancelInviteRequest) (*authv1.CancelInviteResponse, error) {
	if err := s.inviteUC.CancelInvite(ctx, command.CancelInviteCommand{
		Token: req.GetToken(),
	}); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authv1.CancelInviteResponse{
		Success: true,
	}, nil
}

// GetPublicKey returns the RSA public key PEM.
func (s *AuthServer) GetPublicKey(ctx context.Context, req *authv1.GetPublicKeyRequest) (*authv1.GetPublicKeyResponse, error) {
	result, err := s.loginUC.GetPublicKey(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get public key")
	}

	return &authv1.GetPublicKeyResponse{
		PublicKeyPem: result.PublicKeyPEM,
	}, nil
}
