package grpcdapter

import (
	"context"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/application/auth/transformer"
	"github.com/beabys/wms/auth-service/internal/application/auth/usecase"
	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer implements authv1.UserServiceServer.
type UserServer struct {
	authv1.UnimplementedUserServiceServer
	userUC *usecase.UserUseCase
}

// NewUserServer creates a new UserServer.
func NewUserServer(userUC *usecase.UserUseCase) *UserServer {
	return &UserServer{
		userUC: userUC,
	}
}

// CreateUser creates a new user.
func (s *UserServer) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	var customerID *string
	if cid := req.GetCustomerId(); cid != "" {
		customerID = &cid
	}

	result, err := s.userUC.CreateUser(ctx, command.CreateUserCommand{
		Email:      req.GetEmail(),
		Password:   req.GetPassword(),
		Name:       req.GetName(),
		Role:       req.GetRole(),
		CustomerID: customerID,
		CreatedBy:  req.GetCreatedBy(),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authv1.CreateUserResponse{
		User: transformer.UserResultToProto(result),
	}, nil
}

// GetUser retrieves a user by ID.
func (s *UserServer) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	result, err := s.userUC.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &authv1.GetUserResponse{
		User: transformer.UserResultToProto(result),
	}, nil
}

// UpdateUser updates an existing user.
func (s *UserServer) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	cmd := command.UpdateUserCommand{
		UserID: req.GetId(),
	}
	if req.Name != "" {
		n := req.Name
		cmd.Name = &n
	}
	if req.Role != "" {
		r := req.Role
		cmd.Role = &r
	}
	if req.GetActive() != nil {
		active := req.GetActive().GetValue()
		cmd.Active = &active
	}

	result, err := s.userUC.UpdateUser(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authv1.UpdateUserResponse{
		User: transformer.UserResultToProto(result),
	}, nil
}

// ListUsers retrieves users with filtering and pagination.
func (s *UserServer) ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	var customerID *string
	if cid := req.GetCustomerId(); cid != "" {
		customerID = &cid
	}
	var role *string
	if r := req.GetRole(); r != "" {
		role = &r
	}

	result, err := s.userUC.ListUsers(ctx, command.ListUsersQuery{
		CustomerID: customerID,
		Role:       role,
		Page:       int(req.GetPage()),
		PageSize:   int(req.GetPageSize()),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	users := make([]*authv1.User, 0, len(result.Users))
	for _, u := range result.Users {
		uCopy := u
		users = append(users, transformer.UserResultToProto(&uCopy))
	}

	totalPages := 0
	if result.PageSize > 0 {
		totalPages = (result.TotalCount + result.PageSize - 1) / result.PageSize
	}

	return &authv1.ListUsersResponse{
		Users: users,
		Pagination: &commonv1.Pagination{
			Page:       int32(result.Page),
			PageSize:   int32(result.PageSize),
			Total:      int32(result.TotalCount),
			TotalPages: int32(totalPages),
		},
	}, nil
}

// DeleteUser deactivates a user.
func (s *UserServer) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
	if err := s.userUC.DeleteUser(ctx, req.GetId()); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &authv1.DeleteUserResponse{
		Success: true,
	}, nil
}

// AssignRole assigns a role to a user.
func (s *UserServer) AssignRole(ctx context.Context, req *authv1.AssignRoleRequest) (*authv1.AssignRoleResponse, error) {
	if err := s.userUC.AssignRole(ctx, command.AssignRoleCommand{
		UserID: req.GetUserId(),
		Role:   req.GetRole(),
	}); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authv1.AssignRoleResponse{
		Success: true,
	}, nil
}

// ListRoles returns a list of roles (stub).
func (s *UserServer) ListRoles(ctx context.Context, req *authv1.ListRolesRequest) (*authv1.ListRolesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// CreateRole creates a new role (stub).
func (s *UserServer) CreateRole(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}


