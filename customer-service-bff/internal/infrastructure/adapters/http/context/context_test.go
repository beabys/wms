package httpctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKeyJWT(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextKeyJWT, "test-token")
	token, ok := ctx.Value(ContextKeyJWT).(string)
	assert.True(t, ok)
	assert.Equal(t, "test-token", token)
}

func TestContextKeyJWTMissing(t *testing.T) {
	ctx := context.Background()
	token, ok := ctx.Value(ContextKeyJWT).(string)
	assert.False(t, ok)
	assert.Empty(t, token)
}

func TestContextKeyTypeSafety(t *testing.T) {
	// Ensure ContextKey is a distinct type and doesn't collide
	ctx := context.WithValue(context.Background(), ContextKeyJWT, "jwt-value")
	// A string key should not collide with our typed key
	ctx2 := context.WithValue(ctx, "jwt_token", "different-value")
	val, ok := ctx2.Value(ContextKeyJWT).(string)
	assert.True(t, ok)
	assert.Equal(t, "jwt-value", val)
}
