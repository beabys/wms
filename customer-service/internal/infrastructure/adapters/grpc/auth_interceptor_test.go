package grpcdapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNewAuthInterceptor(t *testing.T) {
	i := NewAuthInterceptor(nil, []string{}, nil)
	assert.NotNil(t, i)
}

func TestAuthInterceptor_ExemptMethods(t *testing.T) {
	i := NewAuthInterceptor(nil, []string{"/grpc.health.v1.Health/Check"}, nil)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	resp, err := i.Unary()(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/grpc.health.v1.Health/Check",
	}, handler)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestAuthInterceptor_MissingMetadata(t *testing.T) {
	i := NewAuthInterceptor(nil, []string{}, nil)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	_, err := i.Unary()(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/customer.v1.CustomerService/GetCustomer",
	}, handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing metadata")
}

func TestGetExemptMethods(t *testing.T) {
	methods := GetExemptMethods()
	assert.Contains(t, methods, "/grpc.health.v1.Health/Check")
	assert.Contains(t, methods, "/grpc.health.v1.Health/Watch")
}

func TestClaimsFromContext(t *testing.T) {
	claims := ClaimsFromContext(context.Background())
	assert.Nil(t, claims)

	claims = &Claims{UserID: "user-1", Role: "admin"}
	ctx := context.WithValue(context.Background(), claimsKey, claims)
	got := ClaimsFromContext(ctx)
	assert.NotNil(t, got)
	assert.Equal(t, "user-1", got.UserID)
	assert.Equal(t, "admin", got.Role)
}
