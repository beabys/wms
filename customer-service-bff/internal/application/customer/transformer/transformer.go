package transformer

import (
	"fmt"
	"time"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"

	"github.com/beabys/wms/customer-service-bff/internal/domain/model"
)

// CustomerFromGrpc converts a proto Customer to a model CustomerResponse.
func CustomerFromGrpc(c *customerv1.Customer) *model.CustomerResponse {
	if c == nil {
		return nil
	}
	resp := &model.CustomerResponse{
		ID:             c.GetId(),
		CompanyName:    c.GetCompanyName(),
		Email:          c.GetEmail(),
		Phone:          c.GetPhone(),
		VatNumber:      c.GetVatNumber(),
		Address:        c.GetAddress(),
		City:           c.GetCity(),
		PostalCode:     c.GetPostalCode(),
		Country:        c.GetCountry(),
		Status:         c.GetStatus(),
		CompanyAdminID: c.GetCompanyAdminId(),
	}
	if ts := c.GetCreatedAt(); ts != nil {
		resp.CreatedAt = formatTimestamp(ts)
	}
	return resp
}

// CustomerListFromGrpc converts a slice of proto Customers to a slice of model CustomerResponse.
func CustomerListFromGrpc(customers []*customerv1.Customer) []model.CustomerResponse {
	if customers == nil {
		return nil
	}
	out := make([]model.CustomerResponse, 0, len(customers))
	for _, c := range customers {
		if cr := CustomerFromGrpc(c); cr != nil {
			out = append(out, *cr)
		}
	}
	return out
}

// AuditEntryFromGrpc converts a proto AuditEntry to a model AuditEntry.
func AuditEntryFromGrpc(e *customerv1.AuditEntry) model.AuditEntry {
	if e == nil {
		return model.AuditEntry{}
	}
	return model.AuditEntry{
		ID:          e.GetId(),
		CustomerID:  e.GetCustomerId(),
		Action:      e.GetAction(),
		PerformedBy: e.GetPerformedBy(),
		Details:     e.GetDetails(),
		CreatedAt:   fmt.Sprintf("%d", e.GetCreatedAt()),
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

// CompanyRoleFromGrpc converts a proto CompanyRole to a model CompanyRoleResponse.
func CompanyRoleFromGrpc(r *customerv1.CompanyRole) model.CompanyRoleResponse {
	if r == nil {
		return model.CompanyRoleResponse{}
	}
	return model.CompanyRoleResponse{
		ID:          r.GetId(),
		CustomerID:  r.GetCustomerId(),
		Name:        r.GetName(),
		Permissions: r.GetPermissions(),
		IsDefault:   r.GetIsDefault(),
	}
}

func formatTimestamp(ts *commonv1.Timestamp) string {
	if ts == nil {
		return ""
	}
	return time.Unix(ts.GetSeconds(), int64(ts.GetNanos())).UTC().Format(time.RFC3339)
}
