package httpctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKeyJWT(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextKeyJWT, "test-jwt-token")
	val, ok := ctx.Value(ContextKeyJWT).(string)
	assert.True(t, ok)
	assert.Equal(t, "test-jwt-token", val)
}

func TestContextKeyJWTWithWrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextKeyJWT, 12345)
	_, ok := ctx.Value(ContextKeyJWT).(string)
	assert.False(t, ok)
}

func TestContextKeyJWTEmpty(t *testing.T) {
	// Context without JWT should return nil
	ctx := context.Background()
	assert.Nil(t, ctx.Value(ContextKeyJWT))
}

func TestContextKeyType(t *testing.T) {
	var key ContextKey = "jwt_token"
	assert.Equal(t, ContextKey("jwt_token"), key)
	assert.Equal(t, ContextKeyJWT, key)
}

func TestContextKeyString(t *testing.T) {
	assert.Equal(t, "jwt_token", string(ContextKeyJWT))
}
