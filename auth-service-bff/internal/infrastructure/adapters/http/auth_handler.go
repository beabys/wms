package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/beabys/wms/auth-service-bff/internal/api/v1"
	"github.com/beabys/wms/auth-service-bff/internal/application/auth/validator"
	grpcdapter "github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/grpc"
	"github.com/beabys/wms/auth-service-bff/internal/domain/model"
	"github.com/beabys/wms/pkg/logger"
)

// AuthHandlerService defines what the HTTP handler needs from the use case.
// Interface defined by consumer (handler layer).
type AuthHandlerService interface {
	Login(ctx context.Context, email, password string) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetMe(ctx context.Context) (*model.UserResponse, error)
	InviteUser(ctx context.Context, email, invitedBy string) (*model.InviteUserResponse, error)
	ListInvites(ctx context.Context, status *string, expired *bool, createdAfter, createdBefore *int64, page, pageSize *int) (*model.InviteListResponse, error)
	CancelInvite(ctx context.Context, token string) (*model.CancelInviteResponse, error)
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
	ListUsers(ctx context.Context, role, customerID *string, page, pageSize *int) (*model.UserListResponse, error)
	GetUser(ctx context.Context, id string) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	AssignRole(ctx context.Context, userID, role string) error
	ListRoles(ctx context.Context) (*model.RoleListResponse, error)
	CreateRole(ctx context.Context, req *model.CreateRoleRequest) (*model.RoleResponse, error)
}

// --- Auth handlers ---

// Login authenticates a user.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req v1.LoginJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.ValidateEmail(string(req.Email)); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validator.ValidatePassword(req.Password); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := s.authSvc.Login(r.Context(), string(req.Email), req.Password)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("user logged in",
		logger.LogField{Key: "email", Value: string(req.Email)},
	)
	s.writeSuccess(w, result)
}

// RefreshToken refreshes an access token.
func (s *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req v1.RefreshTokenJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		s.writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	result, err := s.authSvc.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// Logout revokes a refresh token.
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.RefreshToken = ""
	}

	err := s.authSvc.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("user logged out")
	s.writeSuccess(w, map[string]bool{"success": true})
}

// GetMe returns the current user from the JWT.
func (s *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	result, err := s.authSvc.GetMe(r.Context())
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// ListInvites retrieves paginated invite tokens with optional filters.
func (s *Server) ListInvites(w http.ResponseWriter, r *http.Request, params v1.ListInvitesParams) {
	var (
		status                                *string
		expired                               *bool
		createdAfter, createdBefore, page, pageSize *int
	)
	if params.Status != nil && *params.Status != "" {
		status = params.Status
	}
	if params.Expired != nil {
		expired = params.Expired
	}
	if params.CreatedAfter != nil && *params.CreatedAfter > 0 {
		createdAfter = params.CreatedAfter
	}
	if params.CreatedBefore != nil && *params.CreatedBefore > 0 {
		createdBefore = params.CreatedBefore
	}
	if params.Page != nil && *params.Page > 0 {
		page = params.Page
	}
	if params.PageSize != nil && *params.PageSize > 0 {
		pageSize = params.PageSize
	}

	// Convert int -> int64 for the service call
	var ca, cb *int64
	if createdAfter != nil {
		v := int64(*createdAfter)
		ca = &v
	}
	if createdBefore != nil {
		v := int64(*createdBefore)
		cb = &v
	}

	result, err := s.authSvc.ListInvites(r.Context(), status, expired, ca, cb, page, pageSize)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// InviteUser invites a new user (admin only).
func (s *Server) InviteUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.ValidateInviteEmail(req.Email); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := s.authSvc.InviteUser(r.Context(), req.Email, "")
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("user invited",
		logger.LogField{Key: "email", Value: req.Email},
	)
	s.writeSuccess(w, result)
}

// CancelInvite cancels a pending invite token.
func (s *Server) CancelInvite(w http.ResponseWriter, r *http.Request, token string) {
	result, err := s.authSvc.CancelInvite(r.Context(), token)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.logger.Info("invite cancelled",
		logger.LogField{Key: "token", Value: token},
	)
	s.writeSuccess(w, result)
}

// --- User handlers ---

// CreateUser creates a new user.
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.ValidateCreateUser(string(req.Email), req.Password, req.Name, req.Role); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	modelReq := &model.CreateUserRequest{
		Email:    string(req.Email),
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
	}
	if req.CustomerId != nil && *req.CustomerId != "" {
		modelReq.CustomerID = req.CustomerId
	}

	result, err := s.authSvc.CreateUser(r.Context(), modelReq)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// ListUsers lists users with optional filters.
func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request, params v1.ListUsersParams) {
	var (
		role, customerID *string
		page, pageSize   *int
	)
	if params.Role != nil && *params.Role != "" {
		role = params.Role
	}
	if params.CustomerId != nil && *params.CustomerId != "" {
		customerID = params.CustomerId
	}
	if params.Page != nil && *params.Page > 0 {
		page = params.Page
	}
	if params.PageSize != nil && *params.PageSize > 0 {
		pageSize = params.PageSize
	}

	result, err := s.authSvc.ListUsers(r.Context(), role, customerID, page, pageSize)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// GetUser retrieves a user by ID.
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request, id string) {
	result, err := s.authSvc.GetUser(r.Context(), id)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// UpdateUser updates a user's fields.
func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.UpdateUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	modelReq := &model.UpdateUserRequest{
		Name:   req.Name,
		Role:   req.Role,
		Active: req.Active,
	}

	result, err := s.authSvc.UpdateUser(r.Context(), id, modelReq)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// DeleteUser deactivates a user.
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	err := s.authSvc.DeleteUser(r.Context(), id)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, map[string]string{"message": "user deactivated"})
}

// AssignRole assigns a role to a user.
func (s *Server) AssignRole(w http.ResponseWriter, r *http.Request, id string) {
	var req v1.AssignRoleJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role == "" {
		s.writeError(w, http.StatusBadRequest, "role is required")
		return
	}

	err := s.authSvc.AssignRole(r.Context(), id, req.Role)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, map[string]string{"message": "role assigned"})
}

// --- Role handlers ---

// ListRoles lists all roles.
func (s *Server) ListRoles(w http.ResponseWriter, r *http.Request) {
	result, err := s.authSvc.ListRoles(r.Context())
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	s.writeSuccess(w, result)
}

// CreateRole creates a new role.
func (s *Server) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateRoleJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		s.writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	modelReq := &model.CreateRoleRequest{
		Name:        req.Name,
		Description: "",
		Permissions: nil,
	}
	if req.Description != nil {
		modelReq.Description = *req.Description
	}
	if req.Permissions != nil {
		perms := make([]model.RolePermission, 0, len(*req.Permissions))
		for _, p := range *req.Permissions {
			perm := model.RolePermission{}
			if p.Service != nil {
				perm.Service = *p.Service
			}
			if p.Action != nil {
				perm.Action = *p.Action
			}
			if p.Resource != nil {
				perm.Resource = *p.Resource
			}
			perms = append(perms, perm)
		}
		modelReq.Permissions = perms
	}

	result, err := s.authSvc.CreateRole(r.Context(), modelReq)
	if err != nil {
		mapAndWriteError(s, w, err)
		return
	}

	successResponseJSON(w, http.StatusCreated, result)
}

// ---- helpers ----

func mapAndWriteError(s *Server, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	status := grpcdapter.HTTPStatusFromError(err)
	msg := err.Error()
	if errors.Is(err, grpcdapter.ErrInternal) {
		status = http.StatusInternalServerError
		msg = "internal server error"
	}
	s.writeError(w, status, msg)
}


