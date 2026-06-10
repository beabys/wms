package grpcdapter

import (
	"context"
	"testing"

	"github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/http/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthClient(t *testing.T) {
	// grpc.NewClient is lazy — doesn't connect until first RPC
	client, err := NewAuthClient("localhost:9999")
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.NotNil(t, client.conn)
	assert.NotNil(t, client.authService)
	assert.NotNil(t, client.userService)
	defer client.Close()
}

func TestNewAuthClientEmptyTarget(t *testing.T) {
	// An empty target should still create a client (lazy dial)
	client, err := NewAuthClient("")
	require.NoError(t, err)
	assert.NotNil(t, client)
	defer client.Close()
}

func TestAuthClientClose(t *testing.T) {
	client, err := NewAuthClient("localhost:9999")
	require.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)
}

func TestForwardJWTCtxWithToken(t *testing.T) {
	client, err := NewAuthClient("localhost:9999")
	require.NoError(t, err)
	defer client.Close()

	ctx := context.WithValue(context.Background(), httpctx.ContextKeyJWT, "test-jwt-token")
	resultCtx := client.forwardJWTCtx(ctx)
	assert.NotNil(t, resultCtx)
	assert.NotEqual(t, ctx, resultCtx) // metadata modified
}

func TestForwardJWTCtxWithoutToken(t *testing.T) {
	client, err := NewAuthClient("localhost:9999")
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	resultCtx := client.forwardJWTCtx(ctx)
	assert.Equal(t, ctx, resultCtx) // unchanged
}

func TestForwardJWTCtxWithWrongType(t *testing.T) {
	client, err := NewAuthClient("localhost:9999")
	require.NoError(t, err)
	defer client.Close()

	ctx := context.WithValue(context.Background(), httpctx.ContextKeyJWT, 123) // wrong type
	resultCtx := client.forwardJWTCtx(ctx)
	assert.Equal(t, ctx, resultCtx) // unchanged, value not a string
}
