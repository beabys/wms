package transformer

import (
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
)

// CustomerToProto converts a domain customer to a proto customer.
func CustomerToProto(c *model.Customer) *customerv1.Customer {
	if c == nil {
		return nil
	}
	return &customerv1.Customer{
		Id:             c.ID,
		CompanyName:    c.CompanyName,
		Email:          c.Email,
		Phone:          c.Phone,
		VatNumber:      c.VatNumber,
		Address:        c.Address,
		City:           c.City,
		PostalCode:     c.PostalCode,
		Country:        c.Country,
		Status:         string(c.Status),
		CompanyAdminId: c.CompanyAdminID,
		CreatedAt:      timeToProto(c.CreatedAt),
		UpdatedAt:      timeToProto(c.UpdatedAt),
	}
}

// CustomerFromProto converts a proto customer to a domain customer.
func CustomerFromProto(c *customerv1.Customer) *model.Customer {
	if c == nil {
		return nil
	}
	return &model.Customer{
		ID:             c.GetId(),
		CompanyName:    c.GetCompanyName(),
		Email:          c.GetEmail(),
		Phone:          c.GetPhone(),
		VatNumber:      c.GetVatNumber(),
		Address:        c.GetAddress(),
		City:           c.GetCity(),
		PostalCode:     c.GetPostalCode(),
		Country:        c.GetCountry(),
		Status:         model.CustomerStatus(c.GetStatus()),
		CompanyAdminID: c.GetCompanyAdminId(),
		CreatedAt:      protoToTime(c.GetCreatedAt()),
		UpdatedAt:      protoToTime(c.GetUpdatedAt()),
	}
}

// AuditEntryToProto converts a domain audit log to a proto audit entry.
func AuditEntryToProto(e *model.AuditLog) *customerv1.AuditEntry {
	if e == nil {
		return nil
	}
	return &customerv1.AuditEntry{
		Id:          e.ID,
		CustomerId:  e.CustomerID,
		Action:      string(e.Action),
		PerformedBy: e.PerformedBy,
		Details:     e.Details,
		CreatedAt:   e.CreatedAt.Unix(),
	}
}

// AuditLogResultToProto converts a command.AuditLogResult to a proto audit entry.
func AuditLogResultToProto(r *command.AuditLogResult) *customerv1.AuditEntry {
	if r == nil {
		return nil
	}
	return &customerv1.AuditEntry{
		Id:          r.ID,
		CustomerId:  r.CustomerID,
		Action:      r.Action,
		PerformedBy: r.PerformedBy,
		Details:     r.Details,
		CreatedAt:   parseAuditTimestamp(r.CreatedAt),
	}
}

// parseAuditTimestamp parses an RFC3339 string to Unix seconds.
func parseAuditTimestamp(ts string) int64 {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// AuditEntryFromProto converts a proto audit entry to a domain audit log.
func AuditEntryFromProto(e *customerv1.AuditEntry) *model.AuditLog {
	if e == nil {
		return nil
	}
	return &model.AuditLog{
		ID:          e.GetId(),
		CustomerID:  e.GetCustomerId(),
		Action:      model.AuditAction(e.GetAction()),
		PerformedBy: e.GetPerformedBy(),
		Details:     e.GetDetails(),
		CreatedAt:   time.Unix(e.GetCreatedAt(), 0),
	}
}

// CustomerResultToProto converts a command.CustomerResult to a proto customer.
func CustomerResultToProto(r *command.CustomerResult) *customerv1.Customer {
	if r == nil {
		return nil
	}
	return &customerv1.Customer{
		Id:             r.ID,
		CompanyName:    r.CompanyName,
		Email:          r.Email,
		Phone:          r.Phone,
		VatNumber:      r.VatNumber,
		Address:        r.Address,
		City:           r.City,
		PostalCode:     r.PostalCode,
		Country:        r.Country,
		Status:         r.Status,
		CompanyAdminId: r.CompanyAdminID,
		CreatedAt:      parseTimestamp(r.CreatedAt),
		UpdatedAt:      parseTimestamp(r.UpdatedAt),
	}
}

// timeToProto converts time.Time to proto Timestamp.
func timeToProto(t time.Time) *commonv1.Timestamp {
	if t.IsZero() {
		return nil
	}
	return &commonv1.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}

// protoToTime converts proto Timestamp to time.Time.
func protoToTime(ts *commonv1.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
}

// parseTimestamp parses an RFC3339 string into a proto Timestamp.
func parseTimestamp(ts string) *commonv1.Timestamp {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil
	}
	return &commonv1.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}
