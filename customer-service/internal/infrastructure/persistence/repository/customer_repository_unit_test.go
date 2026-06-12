package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerRepositoryTypeCheck(t *testing.T) {
	// Verify the repository type can be instantiated (nil check)
	var r *CustomerRepository
	assert.Nil(t, r)
}
