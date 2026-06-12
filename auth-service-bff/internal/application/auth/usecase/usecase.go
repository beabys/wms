package usecase

import (
	"context"
	"fmt"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/beabys/wms/auth-service-bff/internal/application/auth/transformer"
	"github.com/beabys/wms/auth-service-bff/internal/application/auth/validator"
	grpcdapter "github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/grpc"
	"github.com/beabys/wms/auth-service-bff/internal/domain/model"
	"github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/http/context"
	"github.com/beabys/wms/pkg/logger"
)

// GrpcClient defines what the use case needs from the gRPC layer.
// Interface defined by consumer (usecase).
type GrpcClient interface {
	Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*authv1.RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) (*authv1.LogoutResponse, error)
	ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error)
	InviteUser(ctx context.Context, email, invitedBy string) (*authv1.InviteUserResponse, error)
	ListInvites(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error)
	CancelInvite(ctx context.Context, token string) (*authv1.CancelInviteResponse, error)
	CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error)
	GetUser(ctx context.Context, id string) (*authv1.GetUserResponse, error)
	UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error)
	ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error)
	DeleteUser(ctx context.Context, id string) (*authv1.DeleteUserResponse, error)
	AssignRole(ctx context.Context, userID, role string) (*authv1.AssignRoleResponse, error)
	ListRoles(ctx context.Context, req *authv1.ListRolesRequest) (*authv1.ListRolesResponse, error)
	CreateRole(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error)
}

// AuthUseCase implements business logic for auth operations.
type AuthUseCase struct {
	grpc GrpcClient
	log  logger.Logger
}

// New creates a new AuthUseCase.
func New(log logger.Logger, grpc GrpcClient) *AuthUseCase {
	return &AuthUseCase{grpc: grpc, log: log}
}

// Login authenticates a user and returns tokens.
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*model.AuthResponse, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return nil, fmt.Errorf("%w: %w", grpcdapter.ErrInvalidArgument, err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("%w: %w", grpcdapter.ErrInvalidArgument, err)
	}
	resp, err := uc.grpc.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return &model.AuthResponse{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		ExpiresIn:    resp.GetExpiresIn(),
	}, nil
}

// RefreshToken refreshes an access token.
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("%w: refresh_token is required", grpcdapter.ErrInvalidArgument)
	}
	resp, err := uc.grpc.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return &model.AuthResponse{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		ExpiresIn:    resp.GetExpiresIn(),
	}, nil
}

// Logout revokes a refresh token.
func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	_, err := uc.grpc.Logout(ctx, refreshToken)
	return err
}

// GetMe returns the current user from the JWT in context.
func (uc *AuthUseCase) GetMe(ctx context.Context) (*model.UserResponse, error) {
	token, ok := ctx.Value(httpctx.ContextKeyJWT).(string)
	if !ok || token == "" {
		return nil, grpcdapter.ErrUnauthenticated
	}
	validateResp, err := uc.grpc.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}
	userResp, err := uc.grpc.GetUser(ctx, validateResp.GetUserId())
	if err != nil {
		return nil, err
	}
	return transformer.UserFromGrpc(userResp.GetUser()), nil
}

// InviteUser creates an invite token for a new user.
func (uc *AuthUseCase) InviteUser(ctx context.Context, email, invitedBy string) (*model.InviteUserResponse, error) {
	if err := validator.ValidateInviteEmail(email); err != nil {
		return nil, fmt.Errorf("%w: %w", grpcdapter.ErrInvalidArgument, err)
	}
	resp, err := uc.grpc.InviteUser(ctx, email, invitedBy)
	if err != nil {
		return nil, err
	}
	return &model.InviteUserResponse{
		Token:      resp.GetToken(),
		InviteLink: resp.GetInviteLink(),
		ExpiresAt:  resp.GetExpiresAt(),
	}, nil
}

// ListInvites retrieves paginated invite tokens with optional filters.
func (uc *AuthUseCase) ListInvites(ctx context.Context, status *string, expired *bool, createdAfter, createdBefore *int64, page, pageSize *int) (*model.InviteListResponse, error) {
	protoReq := &authv1.ListInvitesRequest{
		Status:        status,
		Expired:       expired,
		CreatedAfter:  createdAfter,
		CreatedBefore: createdBefore,
	}
	if page != nil {
		protoReq.Page = int32(*page)
	}
	if pageSize != nil {
		protoReq.PageSize = int32(*pageSize)
	}
	resp, err := uc.grpc.ListInvites(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	invites := make([]model.InviteEntry, 0, len(resp.GetInvites()))
	for _, inv := range resp.GetInvites() {
		invites = append(invites, *transformer.InviteFromGrpc(inv))
	}
	return &model.InviteListResponse{
		Invites:    invites,
		Pagination: transformer.PaginationFromGrpc(resp.GetPagination()),
	}, nil
}

// CancelInvite cancels a pending invite token.
func (uc *AuthUseCase) CancelInvite(ctx context.Context, token string) (*model.CancelInviteResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("%w: token is required", grpcdapter.ErrInvalidArgument)
	}
	resp, err := uc.grpc.CancelInvite(ctx, token)
	if err != nil {
		return nil, err
	}
	return &model.CancelInviteResponse{
		Success: resp.GetSuccess(),
	}, nil
}

// CreateUser creates a new user.
func (uc *AuthUseCase) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	if err := validator.ValidateCreateUser(req.Email, req.Password, req.Name, req.Role); err != nil {
		return nil, fmt.Errorf("%w: %w", grpcdapter.ErrInvalidArgument, err)
	}
	protoReq := &authv1.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
	}
	if req.CustomerID != nil {
		protoReq.CustomerId = *req.CustomerID
	}
	resp, err := uc.grpc.CreateUser(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return transformer.UserFromGrpc(resp.GetUser()), nil
}

// ListUsers lists users with optional filters.
func (uc *AuthUseCase) ListUsers(ctx context.Context, role, customerID *string, page, pageSize *int) (*model.UserListResponse, error) {
	protoReq := &authv1.ListUsersRequest{}
	if role != nil {
		protoReq.Role = *role
	}
	if customerID != nil {
		protoReq.CustomerId = *customerID
	}
	if page != nil {
		protoReq.Page = int32(*page)
	}
	if pageSize != nil {
		protoReq.PageSize = int32(*pageSize)
	}
	resp, err := uc.grpc.ListUsers(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	users := make([]model.UserResponse, 0, len(resp.GetUsers()))
	for _, u := range resp.GetUsers() {
		users = append(users, *transformer.UserFromGrpc(u))
	}
	return &model.UserListResponse{
		Users:      users,
		Pagination: transformer.PaginationFromGrpc(resp.GetPagination()),
	}, nil
}

// GetUser retrieves a user by ID.
func (uc *AuthUseCase) GetUser(ctx context.Context, id string) (*model.UserResponse, error) {
	resp, err := uc.grpc.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return transformer.UserFromGrpc(resp.GetUser()), nil
}

// UpdateUser updates a user's fields.
func (uc *AuthUseCase) UpdateUser(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	protoReq := &authv1.UpdateUserRequest{Id: id}
	if req.Name != nil {
		protoReq.Name = *req.Name
	}
	if req.Role != nil {
		protoReq.Role = *req.Role
	}
	if req.Active != nil {
		protoReq.Active = wrapperspb.Bool(*req.Active)
	}
	resp, err := uc.grpc.UpdateUser(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return transformer.UserFromGrpc(resp.GetUser()), nil
}

// DeleteUser deactivates a user.
func (uc *AuthUseCase) DeleteUser(ctx context.Context, id string) error {
	_, err := uc.grpc.DeleteUser(ctx, id)
	return err
}

// AssignRole assigns a role to a user.
func (uc *AuthUseCase) AssignRole(ctx context.Context, userID, role string) error {
	if role == "" {
		return fmt.Errorf("%w: role is required", grpcdapter.ErrInvalidArgument)
	}
	_, err := uc.grpc.AssignRole(ctx, userID, role)
	return err
}

// ListRoles lists all roles.
func (uc *AuthUseCase) ListRoles(ctx context.Context) (*model.RoleListResponse, error) {
	resp, err := uc.grpc.ListRoles(ctx, &authv1.ListRolesRequest{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return nil, err
	}
	roles := make([]model.RoleResponse, 0, len(resp.GetRoles()))
	for _, r := range resp.GetRoles() {
		roles = append(roles, *transformer.RoleFromGrpc(r))
	}
	return &model.RoleListResponse{Roles: roles}, nil
}

// CreateRole creates a new role.
func (uc *AuthUseCase) CreateRole(ctx context.Context, req *model.CreateRoleRequest) (*model.RoleResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("%w: name is required", grpcdapter.ErrInvalidArgument)
	}
	protoReq := &authv1.CreateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
		Permissions: transformer.PermissionsToProto(req.Permissions),
	}
	resp, err := uc.grpc.CreateRole(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return transformer.RoleFromGrpc(resp.GetRole()), nil
}
