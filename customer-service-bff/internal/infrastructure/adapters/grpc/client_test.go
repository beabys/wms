package grpcdapter

import (
	"context"
	"testing"

	"github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/http/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCustomerClient(t *testing.T) {
	client, err := NewCustomerClient("localhost:9999")
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.NotNil(t, client.conn)
	assert.NotNil(t, client.customerService)
	assert.NotNil(t, client.permissionService)
	defer client.Close()
}

func TestNewCustomerClientEmptyTarget(t *testing.T) {
	client, err := NewCustomerClient("")
	require.NoError(t, err)
	assert.NotNil(t, client)
	defer client.Close()
}

func TestCustomerClientClose(t *testing.T) {
	client, err := NewCustomerClient("localhost:9999")
	require.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)
}

func TestForwardJWTCtxWithToken(t *testing.T) {
	client, err := NewCustomerClient("localhost:9999")
	require.NoError(t, err)
	defer client.Close()

	ctx := context.WithValue(context.Background(), httpctx.ContextKeyJWT, "test-jwt-token")
	resultCtx := client.forwardJWTCtx(ctx)
	assert.NotNil(t, resultCtx)
	assert.NotEqual(t, ctx, resultCtx)
}

func TestForwardJWTCtxWithoutToken(t *testing.T) {
	client, err := NewCustomerClient("localhost:9999")
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	resultCtx := client.forwardJWTCtx(ctx)
	assert.Equal(t, ctx, resultCtx)
}

func TestForwardJWTCtxWithWrongType(t *testing.T) {
	client, err := NewCustomerClient("localhost:9999")
	require.NoError(t, err)
	defer client.Close()

	ctx := context.WithValue(context.Background(), httpctx.ContextKeyJWT, 123)
	resultCtx := client.forwardJWTCtx(ctx)
	assert.Equal(t, ctx, resultCtx)
}
