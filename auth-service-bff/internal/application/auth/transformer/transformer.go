package transformer

import (
	"time"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"

	"github.com/beabys/wms/auth-service-bff/internal/domain/model"
)

// UserFromGrpc converts a proto User to a model UserResponse.
func UserFromGrpc(u *authv1.User) *model.UserResponse {
	if u == nil {
		return nil
	}
	resp := &model.UserResponse{
		ID:     u.GetId(),
		Email:  u.GetEmail(),
		Name:   u.GetName(),
		Role:   u.GetRole(),
		Active: u.GetActive(),
	}
	if cid := u.GetCustomerId(); cid != "" {
		resp.CustomerID = &cid
	}
	if ts := u.GetCreatedAt(); ts != nil {
		resp.CreatedAt = time.Unix(ts.GetSeconds(), int64(ts.GetNanos())).UTC().Format(time.RFC3339)
	}
	return resp
}

// RoleFromGrpc converts a proto Role to a model RoleResponse.
func RoleFromGrpc(r *authv1.Role) *model.RoleResponse {
	if r == nil {
		return nil
	}
	perms := make([]model.RolePermission, 0, len(r.GetPermissions()))
	for _, p := range r.GetPermissions() {
		perms = append(perms, model.RolePermission{
			Service:  p.GetService(),
			Action:   p.GetAction(),
			Resource: p.GetResource(),
		})
	}
	return &model.RoleResponse{
		ID:          r.GetId(),
		Name:        r.GetName(),
		Description: r.GetDescription(),
		Permissions: perms,
		IsSystem:    r.GetIsSystem(),
	}
}

// PermissionsToProto converts model permissions to proto permissions.
func PermissionsToProto(perms []model.RolePermission) []*authv1.Permission {
	if perms == nil {
		return nil
	}
	out := make([]*authv1.Permission, 0, len(perms))
	for _, p := range perms {
		out = append(out, &authv1.Permission{
			Service:  p.Service,
			Action:   p.Action,
			Resource: p.Resource,
		})
	}
	return out
}

// InviteFromGrpc converts a proto InviteEntry to a model InviteEntry.
func InviteFromGrpc(invite *authv1.InviteEntry) *model.InviteEntry {
	if invite == nil {
		return nil
	}
	return &model.InviteEntry{
		ID:        invite.GetId(),
		Email:     invite.GetEmail(),
		Token:     invite.GetToken(),
		InvitedBy: invite.GetInvitedBy(),
		Status:    invite.GetStatus(),
		ExpiresAt: invite.GetExpiresAt(),
		CreatedAt: invite.GetCreatedAt(),
	}
}

// PaginationFromGrpc converts a proto Pagination to a model Pagination.
func PaginationFromGrpc(p *commonv1.Pagination) model.Pagination {
	if p == nil {
		return model.Pagination{}
	}
	return model.Pagination{
		Page:       int(p.GetPage()),
		PageSize:   int(p.GetPageSize()),
		TotalItems: int(p.GetTotal()),
	}
}
