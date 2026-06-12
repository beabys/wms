package transformer

import (
	"testing"
	"time"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/auth-service-bff/internal/domain/model"
)

func TestUserFromGrpc(t *testing.T) {
	ts := &commonv1.Timestamp{Seconds: time.Now().Unix()}
	customerID := "cust-1"

	t.Run("full user", func(t *testing.T) {
		u := &authv1.User{
			Id:         "user-1",
			Email:      "test@example.com",
			Name:       "Test User",
			Role:       "admin",
			CustomerId: customerID,
			Active:     true,
			CreatedAt:  ts,
		}
		result := UserFromGrpc(u)
		require.NotNil(t, result)
		assert.Equal(t, "user-1", result.ID)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "Test User", result.Name)
		assert.Equal(t, "admin", result.Role)
		assert.True(t, result.Active)
		require.NotNil(t, result.CustomerID)
		assert.Equal(t, "cust-1", *result.CustomerID)
		assert.NotEmpty(t, result.CreatedAt)
	})

	t.Run("nil user returns nil", func(t *testing.T) {
		assert.Nil(t, UserFromGrpc(nil))
	})

	t.Run("no customer id", func(t *testing.T) {
		u := &authv1.User{
			Id:    "user-2",
			Email: "no-cust@example.com",
			Name:  "No Customer",
			Role:  "viewer",
		}
		result := UserFromGrpc(u)
		require.NotNil(t, result)
		assert.Nil(t, result.CustomerID)
	})

	t.Run("no created at", func(t *testing.T) {
		u := &authv1.User{
			Id:    "user-3",
			Email: "no-ts@example.com",
			Name:  "No Timestamp",
			Role:  "viewer",
		}
		result := UserFromGrpc(u)
		require.NotNil(t, result)
		assert.Empty(t, result.CreatedAt)
	})
}

func TestRoleFromGrpc(t *testing.T) {
	t.Run("full role", func(t *testing.T) {
		r := &authv1.Role{
			Id:          "role-1",
			Name:        "admin",
			Description: "Administrator",
			Permissions: []*authv1.Permission{
				{Service: "users", Action: "read", Resource: "*"},
				{Service: "users", Action: "write", Resource: "*"},
			},
			IsSystem: true,
		}
		result := RoleFromGrpc(r)
		require.NotNil(t, result)
		assert.Equal(t, "role-1", result.ID)
		assert.Equal(t, "admin", result.Name)
		assert.Equal(t, "Administrator", result.Description)
		assert.True(t, result.IsSystem)
		assert.Len(t, result.Permissions, 2)
		assert.Equal(t, "users", result.Permissions[0].Service)
		assert.Equal(t, "write", result.Permissions[1].Action)
	})

	t.Run("nil role returns nil", func(t *testing.T) {
		assert.Nil(t, RoleFromGrpc(nil))
	})

	t.Run("no permissions", func(t *testing.T) {
		r := &authv1.Role{
			Id:   "role-2",
			Name: "viewer",
		}
		result := RoleFromGrpc(r)
		require.NotNil(t, result)
		assert.Empty(t, result.Permissions)
	})
}

func TestPermissionsToProto(t *testing.T) {
	t.Run("converts permissions", func(t *testing.T) {
		perms := []model.RolePermission{
			{Service: "users", Action: "read", Resource: "*"},
			{Service: "orders", Action: "write", Resource: "order:*"},
		}
		result := PermissionsToProto(perms)
		require.Len(t, result, 2)
		assert.Equal(t, "users", result[0].GetService())
		assert.Equal(t, "read", result[0].GetAction())
		assert.Equal(t, "*", result[0].GetResource())
		assert.Equal(t, "orders", result[1].GetService())
		assert.Equal(t, "order:*", result[1].GetResource())
	})

	t.Run("nil input returns nil", func(t *testing.T) {
		assert.Nil(t, PermissionsToProto(nil))
	})

	t.Run("empty input returns empty", func(t *testing.T) {
		result := PermissionsToProto([]model.RolePermission{})
		assert.Empty(t, result)
	})
}

func TestInviteFromGrpc(t *testing.T) {
	t.Run("full invite", func(t *testing.T) {
		inv := &authv1.InviteEntry{
			Id:        "inv-1",
			Email:     "invited@example.com",
			Token:     "token-abc",
			InvitedBy: "admin",
			Status:    "pending",
			ExpiresAt: 1900000000,
			CreatedAt: 1800000000,
		}
		result := InviteFromGrpc(inv)
		require.NotNil(t, result)
		assert.Equal(t, "inv-1", result.ID)
		assert.Equal(t, "invited@example.com", result.Email)
		assert.Equal(t, "token-abc", result.Token)
		assert.Equal(t, "admin", result.InvitedBy)
		assert.Equal(t, "pending", result.Status)
		assert.Equal(t, int64(1900000000), result.ExpiresAt)
		assert.Equal(t, int64(1800000000), result.CreatedAt)
	})

	t.Run("nil invite returns nil", func(t *testing.T) {
		assert.Nil(t, InviteFromGrpc(nil))
	})
}

func TestPaginationFromGrpc(t *testing.T) {
	t.Run("full pagination", func(t *testing.T) {
		p := &commonv1.Pagination{
			Page:     2,
			PageSize: 20,
			Total:    50,
		}
		result := PaginationFromGrpc(p)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 20, result.PageSize)
		assert.Equal(t, 50, result.TotalItems)
	})

	t.Run("nil pagination returns zero values", func(t *testing.T) {
		result := PaginationFromGrpc(nil)
		assert.Equal(t, 0, result.Page)
		assert.Equal(t, 0, result.PageSize)
		assert.Equal(t, 0, result.TotalItems)
	})

	t.Run("zero values handled", func(t *testing.T) {
		p := &commonv1.Pagination{}
		result := PaginationFromGrpc(p)
		assert.Equal(t, 0, result.Page)
		assert.Equal(t, 0, result.PageSize)
		assert.Equal(t, 0, result.TotalItems)
	})
}

func TestUserFromGrpc_TimestampFormat(t *testing.T) {
	// Fixed timestamp: 2024-01-15T10:30:00Z
	ts := &commonv1.Timestamp{Seconds: 1705314600, Nanos: 0}
	u := &authv1.User{
		Id:        "u-1",
		Email:     "ts@test.com",
		Name:      "TS Test",
		Role:      "viewer",
		CreatedAt: ts,
	}
	result := UserFromGrpc(u)
	require.NotNil(t, result)
	assert.Equal(t, "2024-01-15T10:30:00Z", result.CreatedAt)
}
