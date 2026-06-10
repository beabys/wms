package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	a := New()
	assert.NotNil(t, a)
	assert.Nil(t, a.Config)
	assert.Nil(t, a.Logger)
	assert.Nil(t, a.DB)
	assert.Nil(t, a.GrpcServer)
	assert.Nil(t, a.AuthClient)
	assert.Nil(t, a.CustomerRepo)
	assert.Nil(t, a.CompanyRepo)
	assert.Nil(t, a.CompanyRoleRepo)
}
