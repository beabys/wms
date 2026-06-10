package grpcdapter

import (
	"context"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/usecase"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
)

// PermissionServer implements customerv1.PermissionServiceServer.
type PermissionServer struct {
	customerv1.UnimplementedPermissionServiceServer
	companyUC *usecase.CompanyUseCase
}

// NewPermissionServer creates a new PermissionServer.
func NewPermissionServer(companyUC *usecase.CompanyUseCase) *PermissionServer {
	return &PermissionServer{
		companyUC: companyUC,
	}
}

// AssignCompanyRole assigns a role with permissions to a user.
func (s *PermissionServer) AssignCompanyRole(ctx context.Context, req *customerv1.AssignCompanyRoleRequest) (*customerv1.AssignCompanyRoleResponse, error) {
	err := s.companyUC.AssignCompanyRole(ctx, command.AssignCompanyRoleCommand{
		CustomerID:  req.GetCustomerId(),
		UserID:      req.GetUserId(),
		RoleName:    req.GetRoleName(),
		Permissions: req.GetPermissions(),
		AssignedBy:  req.GetAssignedBy(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.AssignCompanyRoleResponse{
		Success: true,
	}, nil
}

// ListCompanyRoles returns all roles for a customer.
func (s *PermissionServer) ListCompanyRoles(ctx context.Context, req *customerv1.ListCompanyRolesRequest) (*customerv1.ListCompanyRolesResponse, error) {
	roles, err := s.companyUC.ListCompanyRoles(ctx, req.GetCustomerId())
	if err != nil {
		return nil, mapError(err)
	}

	protoRoles := make([]*customerv1.CompanyRole, 0, len(roles))
	for _, r := range roles {
		protoRoles = append(protoRoles, &customerv1.CompanyRole{
			Id:          r.ID,
			CustomerId:  r.CustomerID,
			Name:        r.Name,
			Permissions: r.Permissions,
			IsDefault:   r.IsDefault,
		})
	}

	return &customerv1.ListCompanyRolesResponse{
		Roles: protoRoles,
	}, nil
}

// GetUserPermissions retrieves aggregated permissions for a user.
func (s *PermissionServer) GetUserPermissions(ctx context.Context, req *customerv1.GetUserPermissionsRequest) (*customerv1.GetUserPermissionsResponse, error) {
	permissions, err := s.companyUC.GetUserPermissions(ctx, req.GetUserId(), req.GetCustomerId())
	if err != nil {
		return nil, mapError(err)
	}

	return &customerv1.GetUserPermissionsResponse{
		Permissions: permissions,
	}, nil
}
