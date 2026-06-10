package grpcdapter

import (
	"testing"

	"github.com/beabys/wms/customer-service/internal/application/customer/usecase"
	"github.com/stretchr/testify/assert"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
)

func TestNewPermissionServer(t *testing.T) {
	uc := &usecase.CompanyUseCase{}
	s := NewPermissionServer(uc)
	assert.NotNil(t, s)
}

func TestPermissionServer_ImplementsInterface(t *testing.T) {
	uc := &usecase.CompanyUseCase{}
	s := NewPermissionServer(uc)
	var _ customerv1.PermissionServiceServer = s
}
