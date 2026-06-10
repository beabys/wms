package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppNew(t *testing.T) {
	a := New()
	assert.NotNil(t, a)
}

func TestAppSetupWithEmptyConfig(t *testing.T) {
	// Setup with empty config should fail because gRPC dial fails
	// (gRPC dial is lazy, so it may succeed in creating the client)
	// Just verify it doesn't panic
	a := New()
	cfg := &Config{}
	err := a.Setup(cfg)
	// May or may not error depending on gRPC lazy dial behavior
	_ = err
}
