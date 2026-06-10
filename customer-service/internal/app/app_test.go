package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReturnsEmptyApp(t *testing.T) {
	a := New()
	assert.NotNil(t, a)
}

func TestAppRunWithoutSetupReturnsError(t *testing.T) {
	a := New()
	err := a.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grpc server not configured")
}

func TestAppShutdownWithoutSetupDoesNotPanic(t *testing.T) {
	a := New()
	assert.NotPanics(t, func() {
		a.Shutdown()
	})
}

func TestAppTypesCompile(t *testing.T) {
	var a *App
	assert.Nil(t, a)
}
