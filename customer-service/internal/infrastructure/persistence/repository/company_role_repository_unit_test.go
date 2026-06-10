package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompanyRoleRepositoryTypeCheck(t *testing.T) {
	var r *CompanyRoleRepository
	assert.Nil(t, r)
}
