package transformer

import (
	"testing"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/customer-service-bff/internal/domain/model"
)

func TestCustomerFromGrpc(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		assert.Nil(t, CustomerFromGrpc(nil))
	})

	t.Run("full conversion", func(t *testing.T) {
		proto := &customerv1.Customer{
			Id:             "cust-1",
			CompanyName:    "ACME Corp",
			Email:          "admin@acme.com",
			Phone:          "555-0100",
			VatNumber:      "VAT123",
			Address:        "123 Main St",
			City:           "New York",
			PostalCode:     "10001",
			Country:        "US",
			Status:         "active",
			CompanyAdminId: "admin-1",
			CreatedAt: &commonv1.Timestamp{
				Seconds: 1700000000,
				Nanos:   0,
			},
		}

		result := CustomerFromGrpc(proto)
		require.NotNil(t, result)
		assert.Equal(t, "cust-1", result.ID)
		assert.Equal(t, "ACME Corp", result.CompanyName)
		assert.Equal(t, "admin@acme.com", result.Email)
		assert.Equal(t, "555-0100", result.Phone)
		assert.Equal(t, "VAT123", result.VatNumber)
		assert.Equal(t, "123 Main St", result.Address)
		assert.Equal(t, "New York", result.City)
		assert.Equal(t, "10001", result.PostalCode)
		assert.Equal(t, "US", result.Country)
		assert.Equal(t, "active", result.Status)
		assert.Equal(t, "admin-1", result.CompanyAdminID)
		assert.Equal(t, "2023-11-14T22:13:20Z", result.CreatedAt)
	})

	t.Run("nil timestamp", func(t *testing.T) {
		proto := &customerv1.Customer{
			Id: "cust-1",
		}
		result := CustomerFromGrpc(proto)
		require.NotNil(t, result)
		assert.Equal(t, "", result.CreatedAt)
	})
}

func TestCustomerListFromGrpc(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		assert.Nil(t, CustomerListFromGrpc(nil))
	})

	t.Run("empty slice", func(t *testing.T) {
		result := CustomerListFromGrpc([]*customerv1.Customer{})
		assert.Empty(t, result)
	})

	t.Run("multiple customers", func(t *testing.T) {
		proto := []*customerv1.Customer{
			{Id: "c1", CompanyName: "Co 1", Status: "active"},
			{Id: "c2", CompanyName: "Co 2", Status: "pending"},
		}
		result := CustomerListFromGrpc(proto)
		assert.Len(t, result, 2)
		assert.Equal(t, "c1", result[0].ID)
		assert.Equal(t, "Co 2", result[1].CompanyName)
	})

	t.Run("nil element in slice", func(t *testing.T) {
		proto := []*customerv1.Customer{
			{Id: "c1"},
			nil,
		}
		result := CustomerListFromGrpc(proto)
		assert.Len(t, result, 1)
	})
}

func TestAuditEntryFromGrpc(t *testing.T) {
	t.Run("nil input returns empty", func(t *testing.T) {
		result := AuditEntryFromGrpc(nil)
		assert.Equal(t, model.AuditEntry{}, result)
	})

	t.Run("full conversion", func(t *testing.T) {
		proto := &customerv1.AuditEntry{
			Id:          "audit-1",
			CustomerId:  "cust-1",
			Action:      "approved",
			PerformedBy: "admin-1",
			Details:     "details here",
			CreatedAt:   1700000000,
		}
		result := AuditEntryFromGrpc(proto)
		assert.Equal(t, "audit-1", result.ID)
		assert.Equal(t, "cust-1", result.CustomerID)
		assert.Equal(t, "approved", result.Action)
		assert.Equal(t, "admin-1", result.PerformedBy)
		assert.Equal(t, "details here", result.Details)
		assert.Equal(t, "1700000000", result.CreatedAt)
	})
}

func TestPaginationFromGrpc(t *testing.T) {
	t.Run("nil input returns empty", func(t *testing.T) {
		result := PaginationFromGrpc(nil)
		assert.Equal(t, model.Pagination{}, result)
	})

	t.Run("full conversion", func(t *testing.T) {
		proto := &commonv1.Pagination{
			Page:     2,
			PageSize: 20,
			Total:    50,
		}
		result := PaginationFromGrpc(proto)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 20, result.PageSize)
		assert.Equal(t, 50, result.TotalItems)
	})
}

func TestCompanyRoleFromGrpc(t *testing.T) {
	t.Run("nil input returns empty", func(t *testing.T) {
		result := CompanyRoleFromGrpc(nil)
		assert.Equal(t, model.CompanyRoleResponse{}, result)
	})

	t.Run("full conversion", func(t *testing.T) {
		proto := &customerv1.CompanyRole{
			Id:          "role-1",
			CustomerId:  "cust-1",
			Name:        "admin",
			Permissions: []string{"read", "write"},
			IsDefault:   false,
		}
		result := CompanyRoleFromGrpc(proto)
		assert.Equal(t, "role-1", result.ID)
		assert.Equal(t, "cust-1", result.CustomerID)
		assert.Equal(t, "admin", result.Name)
		assert.Equal(t, []string{"read", "write"}, result.Permissions)
		assert.Equal(t, false, result.IsDefault)
	})
}
