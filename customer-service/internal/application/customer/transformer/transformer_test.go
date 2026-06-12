package transformer

import (
	"testing"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomerToProto(t *testing.T) {
	now := time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC)
	c := &model.Customer{
		ID:             "cust-1",
		CompanyName:    "Test Corp",
		Email:          "admin@testcorp.com",
		Phone:          "1234567890",
		VatNumber:      "VAT123",
		Address:        "123 Main St",
		City:           "New York",
		PostalCode:     "10001",
		Country:        "US",
		Status:         model.CustomerStatusActive,
		CompanyAdminID: "admin-1",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	p := CustomerToProto(c)
	require.NotNil(t, p)
	assert.Equal(t, "cust-1", p.GetId())
	assert.Equal(t, "Test Corp", p.GetCompanyName())
	assert.Equal(t, "admin@testcorp.com", p.GetEmail())
	assert.Equal(t, "1234567890", p.GetPhone())
	assert.Equal(t, "VAT123", p.GetVatNumber())
	assert.Equal(t, "123 Main St", p.GetAddress())
	assert.Equal(t, "New York", p.GetCity())
	assert.Equal(t, "10001", p.GetPostalCode())
	assert.Equal(t, "US", p.GetCountry())
	assert.Equal(t, "active", p.GetStatus())
	assert.Equal(t, "admin-1", p.GetCompanyAdminId())
	require.NotNil(t, p.GetCreatedAt())
	assert.Equal(t, now.Unix(), p.GetCreatedAt().GetSeconds())
}

func TestCustomerToProto_Nil(t *testing.T) {
	assert.Nil(t, CustomerToProto(nil))
}

func TestCustomerToProto_ZeroTime(t *testing.T) {
	c := &model.Customer{
		ID:    "cust-1",
		Email: "admin@testcorp.com",
	}
	p := CustomerToProto(c)
	require.NotNil(t, p)
	assert.Nil(t, p.GetCreatedAt())
	assert.Nil(t, p.GetUpdatedAt())
}

func TestCustomerFromProto(t *testing.T) {
	now := time.Date(2024, 6, 10, 12, 0, 0, 500000000, time.UTC)
	p := &customerv1.Customer{
		Id:             "cust-1",
		CompanyName:    "Test Corp",
		Email:          "admin@testcorp.com",
		Phone:          "1234567890",
		VatNumber:      "VAT123",
		Address:        "123 Main St",
		City:           "New York",
		PostalCode:     "10001",
		Country:        "US",
		Status:         "active",
		CompanyAdminId: "admin-1",
		CreatedAt: &commonv1.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		UpdatedAt: &commonv1.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	}

	c := CustomerFromProto(p)
	require.NotNil(t, c)
	assert.Equal(t, "cust-1", c.ID)
	assert.Equal(t, "Test Corp", c.CompanyName)
	assert.Equal(t, "admin@testcorp.com", c.Email)
	assert.Equal(t, "1234567890", c.Phone)
	assert.Equal(t, "VAT123", c.VatNumber)
	assert.Equal(t, "123 Main St", c.Address)
	assert.Equal(t, "New York", c.City)
	assert.Equal(t, "10001", c.PostalCode)
	assert.Equal(t, "US", c.Country)
	assert.Equal(t, model.CustomerStatusActive, c.Status)
	assert.Equal(t, "admin-1", c.CompanyAdminID)
	assert.Equal(t, now.Unix(), c.CreatedAt.Unix())
	assert.Equal(t, 500000000, c.CreatedAt.Nanosecond())
}

func TestCustomerFromProto_Nil(t *testing.T) {
	assert.Nil(t, CustomerFromProto(nil))
}

func TestCustomerFromProto_NilTimestamps(t *testing.T) {
	p := &customerv1.Customer{
		Id:    "cust-1",
		Email: "admin@testcorp.com",
	}
	c := CustomerFromProto(p)
	require.NotNil(t, c)
	assert.True(t, c.CreatedAt.IsZero())
	assert.True(t, c.UpdatedAt.IsZero())
}

func TestAuditEntryToProto(t *testing.T) {
	now := time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC)
	e := &model.AuditLog{
		ID:          "audit-1",
		CustomerID:  "cust-1",
		Action:      model.AuditActionApproved,
		PerformedBy: "admin-1",
		Details:     "Customer approved",
		CreatedAt:   now,
	}

	p := AuditEntryToProto(e)
	require.NotNil(t, p)
	assert.Equal(t, "audit-1", p.GetId())
	assert.Equal(t, "cust-1", p.GetCustomerId())
	assert.Equal(t, "approved", p.GetAction())
	assert.Equal(t, "admin-1", p.GetPerformedBy())
	assert.Equal(t, "Customer approved", p.GetDetails())
	assert.Equal(t, now.Unix(), p.GetCreatedAt())
}

func TestAuditEntryToProto_Nil(t *testing.T) {
	assert.Nil(t, AuditEntryToProto(nil))
}

func TestAuditEntryFromProto(t *testing.T) {
	now := time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC)
	p := &customerv1.AuditEntry{
		Id:          "audit-1",
		CustomerId:  "cust-1",
		Action:      "approved",
		PerformedBy: "admin-1",
		Details:     "Customer approved",
		CreatedAt:   now.Unix(),
	}

	e := AuditEntryFromProto(p)
	require.NotNil(t, e)
	assert.Equal(t, "audit-1", e.ID)
	assert.Equal(t, "cust-1", e.CustomerID)
	assert.Equal(t, model.AuditActionApproved, e.Action)
	assert.Equal(t, "admin-1", e.PerformedBy)
	assert.Equal(t, "Customer approved", e.Details)
	assert.Equal(t, now.Unix(), e.CreatedAt.Unix())
}

func TestAuditEntryFromProto_Nil(t *testing.T) {
	assert.Nil(t, AuditEntryFromProto(nil))
}

func TestCustomerResultToProto(t *testing.T) {
	r := &command.CustomerResult{
		ID:             "cust-1",
		CompanyName:    "Test Corp",
		Email:          "admin@testcorp.com",
		Phone:          "1234567890",
		VatNumber:      "VAT123",
		Address:        "123 Main St",
		City:           "New York",
		PostalCode:     "10001",
		Country:        "US",
		Status:         "active",
		CompanyAdminID: "admin-1",
		RejectReason:   "",
		CreatedAt:      "2024-06-10T12:00:00Z",
		UpdatedAt:      "2024-06-10T12:00:00Z",
	}

	p := CustomerResultToProto(r)
	require.NotNil(t, p)
	assert.Equal(t, "cust-1", p.GetId())
	assert.Equal(t, "Test Corp", p.GetCompanyName())
	assert.Equal(t, "admin@testcorp.com", p.GetEmail())
	assert.Equal(t, "1234567890", p.GetPhone())
	assert.Equal(t, "VAT123", p.GetVatNumber())
	assert.Equal(t, "123 Main St", p.GetAddress())
	assert.Equal(t, "New York", p.GetCity())
	assert.Equal(t, "10001", p.GetPostalCode())
	assert.Equal(t, "US", p.GetCountry())
	assert.Equal(t, "active", p.GetStatus())
	assert.Equal(t, "admin-1", p.GetCompanyAdminId())
	require.NotNil(t, p.GetCreatedAt())
	assert.Equal(t, int64(1718020800), p.GetCreatedAt().GetSeconds())
}

func TestCustomerResultToProto_Nil(t *testing.T) {
	assert.Nil(t, CustomerResultToProto(nil))
}

func TestCustomerResultToProto_NoTimestamps(t *testing.T) {
	r := &command.CustomerResult{
		ID:     "cust-2",
		Status: "pending",
	}
	p := CustomerResultToProto(r)
	require.NotNil(t, p)
	assert.Nil(t, p.GetCreatedAt())
}

func TestAuditLogResultToProto(t *testing.T) {
	r := &command.AuditLogResult{
		ID:          "audit-1",
		CustomerID:  "cust-1",
		Action:      "approved",
		PerformedBy: "admin-1",
		Details:     "Customer approved",
		CreatedAt:   "2024-06-10T12:00:00Z",
	}

	p := AuditLogResultToProto(r)
	require.NotNil(t, p)
	assert.Equal(t, "audit-1", p.GetId())
	assert.Equal(t, "cust-1", p.GetCustomerId())
	assert.Equal(t, "approved", p.GetAction())
	assert.Equal(t, "admin-1", p.GetPerformedBy())
	assert.Equal(t, "Customer approved", p.GetDetails())
	assert.Equal(t, int64(1718020800), p.GetCreatedAt())
}

func TestAuditLogResultToProto_Nil(t *testing.T) {
	assert.Nil(t, AuditLogResultToProto(nil))
}

func TestAuditLogResultToProto_EmptyCreatedAt(t *testing.T) {
	r := &command.AuditLogResult{
		ID:          "audit-1",
		CustomerID:  "cust-1",
		Action:      "approved",
		PerformedBy: "admin-1",
		Details:     "Customer approved",
		CreatedAt:   "",
	}

	p := AuditLogResultToProto(r)
	require.NotNil(t, p)
	assert.Equal(t, int64(0), p.GetCreatedAt())
}
