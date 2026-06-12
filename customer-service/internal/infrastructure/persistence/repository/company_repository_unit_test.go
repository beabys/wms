package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompanyRepositoryTypeCheck(t *testing.T) {
	var r *CompanyRepository
	assert.Nil(t, r)
}
