package grpcdapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGRPCServer(t *testing.T) {
	gs := &GRPCServer{}
	assert.NotNil(t, gs)
}

func TestGRPCServer_NilServer(t *testing.T) {
	var gs *GRPCServer
	assert.Nil(t, gs)
}
