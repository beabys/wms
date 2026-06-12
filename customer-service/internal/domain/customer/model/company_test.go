package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCompanyStruct(t *testing.T) {
	now := time.Now()
	c := &Company{
		ID:         "comp-1",
		CustomerID: "cust-1",
		Name:       "Test Company",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	assert.Equal(t, "comp-1", c.ID)
	assert.Equal(t, "cust-1", c.CustomerID)
	assert.Equal(t, "Test Company", c.Name)
	assert.False(t, c.CreatedAt.IsZero())
	assert.False(t, c.UpdatedAt.IsZero())
}
