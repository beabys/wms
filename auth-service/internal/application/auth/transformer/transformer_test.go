package transformer_test

import (
	"testing"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/application/auth/transformer"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	"github.com/stretchr/testify/assert"
)

func TestUserToProto(t *testing.T) {
	now := time.Now()
	customerID := "cust-123"

	t.Run("converts full user", func(t *testing.T) {
		u := &model.User{
			ID:         "user-1",
			Email:      "test@example.com",
			Name:       "Test User",
			Role:       "admin",
			CustomerID: &customerID,
			Active:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		p := transformer.UserToProto(u)
		assert.NotNil(t, p)
		assert.Equal(t, "user-1", p.GetId())
		assert.Equal(t, "test@example.com", p.GetEmail())
		assert.Equal(t, "Test User", p.GetName())
		assert.Equal(t, "admin", p.GetRole())
		assert.Equal(t, "cust-123", p.GetCustomerId())
		assert.True(t, p.GetActive())
		assert.NotNil(t, p.GetCreatedAt())
		assert.Equal(t, now.Unix(), p.GetCreatedAt().GetSeconds())
	})

	t.Run("handles nil customer ID", func(t *testing.T) {
		u := &model.User{
			ID:    "user-2",
			Email: "nil@example.com",
			Name:  "Nil CID",
			Role:  "customer",
		}

		p := transformer.UserToProto(u)
		assert.NotNil(t, p)
		assert.Empty(t, p.GetCustomerId())
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		p := transformer.UserToProto(nil)
		assert.Nil(t, p)
	})
}

func TestUserFromProto(t *testing.T) {
	now := time.Now()

	t.Run("converts full proto user", func(t *testing.T) {
		p := &authv1.User{
			Id:         "user-1",
			Email:      "test@example.com",
			Name:       "Test User",
			Role:       "admin",
			CustomerId: "cust-123",
			Active:     true,
			CreatedAt: &commonv1.Timestamp{Seconds: now.Unix(), Nanos: int32(now.Nanosecond())},
			UpdatedAt: &commonv1.Timestamp{Seconds: now.Unix(), Nanos: int32(now.Nanosecond())},
		}

		u := transformer.UserFromProto(p)
		assert.NotNil(t, u)
		assert.Equal(t, "user-1", u.ID)
		assert.Equal(t, "test@example.com", u.Email)
		assert.Equal(t, "Test User", u.Name)
		assert.Equal(t, "admin", u.Role)
		assert.NotNil(t, u.CustomerID)
		assert.Equal(t, "cust-123", *u.CustomerID)
		assert.True(t, u.Active)
	})

	t.Run("handles empty customer ID", func(t *testing.T) {
		p := &authv1.User{
			Id:    "user-2",
			Email: "nil@example.com",
		}

		u := transformer.UserFromProto(p)
		assert.NotNil(t, u)
		assert.Nil(t, u.CustomerID)
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		u := transformer.UserFromProto(nil)
		assert.Nil(t, u)
	})
}

func TestUserResultToProto(t *testing.T) {
	customerID := "cust-123"

	t.Run("converts full UserResult", func(t *testing.T) {
		u := &command.UserResult{
			ID:         "user-1",
			Email:      "test@example.com",
			Name:       "Test User",
			Role:       "admin",
			CustomerID: &customerID,
			Active:     true,
			CreatedAt:  "2024-01-01T00:00:00Z",
			UpdatedAt:  "2024-01-01T00:00:00Z",
		}

		p := transformer.UserResultToProto(u)
		assert.NotNil(t, p)
		assert.Equal(t, "user-1", p.GetId())
		assert.Equal(t, "test@example.com", p.GetEmail())
		assert.Equal(t, "admin", p.GetRole())
		assert.Equal(t, "cust-123", p.GetCustomerId())
		assert.True(t, p.GetActive())
	})

	t.Run("handles nil customer ID", func(t *testing.T) {
		u := &command.UserResult{
			ID:    "user-2",
			Email: "test@example.com",
		}

		p := transformer.UserResultToProto(u)
		assert.NotNil(t, p)
		assert.Empty(t, p.GetCustomerId())
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		p := transformer.UserResultToProto(nil)
		assert.Nil(t, p)
	})

	t.Run("handles invalid timestamp gracefully", func(t *testing.T) {
		u := &command.UserResult{
			ID:        "user-3",
			Email:     "test@example.com",
			CreatedAt: "invalid-date",
		}

		p := transformer.UserResultToProto(u)
		assert.NotNil(t, p)
		assert.Nil(t, p.GetCreatedAt())
	})
}

func TestInviteToProto(t *testing.T) {
	now := time.Now()

	t.Run("converts full invite", func(t *testing.T) {
		i := &model.InviteToken{
			ID:        "inv-1",
			Email:     "invited@example.com",
			Token:     "token-123",
			InvitedBy: "admin-1",
			Status:    "pending",
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
		}

		p := transformer.InviteToProto(i)
		assert.NotNil(t, p)
		assert.Equal(t, "inv-1", p.GetId())
		assert.Equal(t, "invited@example.com", p.GetEmail())
		assert.Equal(t, "token-123", p.GetToken())
		assert.Equal(t, "admin-1", p.GetInvitedBy())
		assert.Equal(t, "pending", p.GetStatus())
		assert.Equal(t, now.Add(24*time.Hour).Unix(), p.GetExpiresAt())
		assert.Equal(t, now.Unix(), p.GetCreatedAt())
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		p := transformer.InviteToProto(nil)
		assert.Nil(t, p)
	})
}

func TestInviteFromProto(t *testing.T) {
	now := time.Now()

	t.Run("converts full proto invite", func(t *testing.T) {
		p := &authv1.InviteEntry{
			Id:        "inv-1",
			Email:     "invited@example.com",
			Token:     "token-123",
			InvitedBy: "admin-1",
			Status:    "pending",
			ExpiresAt: now.Add(24 * time.Hour).Unix(),
			CreatedAt: now.Unix(),
		}

		i := transformer.InviteFromProto(p)
		assert.NotNil(t, i)
		assert.Equal(t, "inv-1", i.ID)
		assert.Equal(t, "invited@example.com", i.Email)
		assert.Equal(t, "token-123", i.Token)
		assert.Equal(t, "admin-1", i.InvitedBy)
		assert.Equal(t, "pending", i.Status)
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		i := transformer.InviteFromProto(nil)
		assert.Nil(t, i)
	})
}

func TestRoleToProto(t *testing.T) {
	t.Run("converts full role with permissions", func(t *testing.T) {
		r := &model.Role{
			ID:          "role-1",
			Name:        "admin",
			Description: "Administrator role",
			Permissions: []model.Permission{
				{Service: model.ServiceAuth, Action: model.ActionRead, Resource: model.ResourceUsers},
				{Service: model.ServiceAuth, Action: model.ActionWrite, Resource: model.ResourceRoles},
			},
			IsSystem:  true,
			CreatedAt: time.Now(),
		}

		p := transformer.RoleToProto(r)
		assert.NotNil(t, p)
		assert.Equal(t, "role-1", p.GetId())
		assert.Equal(t, "admin", p.GetName())
		assert.Equal(t, "Administrator role", p.GetDescription())
		assert.True(t, p.GetIsSystem())
		assert.Len(t, p.GetPermissions(), 2)
		assert.Equal(t, "auth", p.GetPermissions()[0].GetService())
		assert.Equal(t, "read", p.GetPermissions()[0].GetAction())
	})

	t.Run("handles role with no permissions", func(t *testing.T) {
		r := &model.Role{
			ID:   "role-2",
			Name: "viewer",
		}

		p := transformer.RoleToProto(r)
		assert.NotNil(t, p)
		assert.Empty(t, p.GetPermissions())
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		p := transformer.RoleToProto(nil)
		assert.Nil(t, p)
	})
}

func TestPermissionToProto(t *testing.T) {
	t.Run("converts permission with typed enums to proto strings", func(t *testing.T) {
		p := &model.Permission{
			Service:  model.ServiceAuth,
			Action:   model.ActionRead,
			Resource: model.ResourceUsers,
		}

		pp := transformer.PermissionToProto(p)
		assert.NotNil(t, pp)
		assert.Equal(t, "auth", pp.GetService())
		assert.Equal(t, "read", pp.GetAction())
		assert.Equal(t, "users", pp.GetResource())
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		pp := transformer.PermissionToProto(nil)
		assert.Nil(t, pp)
	})
}

func TestPermissionFromProto(t *testing.T) {
	t.Run("converts proto permission to domain with typed enums", func(t *testing.T) {
		pp := &authv1.Permission{
			Service:  "auth",
			Action:   "read",
			Resource: "users",
		}

		p := transformer.PermissionFromProto(pp)
		assert.NotNil(t, p)
		assert.Equal(t, model.ServiceAuth, p.Service)
		assert.Equal(t, model.ActionRead, p.Action)
		assert.Equal(t, model.ResourceUsers, p.Resource)
	})

	t.Run("preserves unknown string values", func(t *testing.T) {
		pp := &authv1.Permission{
			Service:  "unknown-svc",
			Action:   "unknown-act",
			Resource: "unknown-res",
		}

		p := transformer.PermissionFromProto(pp)
		assert.NotNil(t, p)
		assert.Equal(t, model.Service("unknown-svc"), p.Service)
		assert.Equal(t, model.Action("unknown-act"), p.Action)
		assert.Equal(t, model.Resource("unknown-res"), p.Resource)
	})

	t.Run("handles empty strings gracefully", func(t *testing.T) {
		pp := &authv1.Permission{}

		p := transformer.PermissionFromProto(pp)
		assert.NotNil(t, p)
		assert.Empty(t, string(p.Service))
		assert.Empty(t, string(p.Action))
		assert.Empty(t, string(p.Resource))
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		p := transformer.PermissionFromProto(nil)
		assert.Nil(t, p)
	})
}
