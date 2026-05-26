package authinterceptor_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/beabys/wms/pkg/authinterceptor"
)

// generateTestKey creates an RSA key pair for testing.
func generateTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key
}

// signToken creates a signed JWT with the given claims.
func signToken(t *testing.T, key *rsa.PrivateKey, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	require.NoError(t, err)
	return signed
}

// testHandler is a simple gRPC handler that returns a success marker.
var testHandler = func(ctx context.Context, req interface{}) (interface{}, error) {
	return map[string]string{"status": "ok"}, nil
}

func TestJWTAuthInterceptor_ExemptMethod(t *testing.T) {
	key := generateTestKey(t)
	interceptor := authinterceptor.JWTAuthInterceptor(&key.PublicKey, []string{"/auth.AuthService/Login"})

	info := &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/Login"}
	resp, err := interceptor(context.Background(), nil, info, testHandler)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"status": "ok"}, resp)
}

func TestJWTAuthInterceptor_MissingMetadata(t *testing.T) {
	key := generateTestKey(t)
	interceptor := authinterceptor.JWTAuthInterceptor(&key.PublicKey, nil)

	info := &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/ValidateToken"}
	_, err := interceptor(context.Background(), nil, info, testHandler)
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestJWTAuthInterceptor_MissingAuthHeader(t *testing.T) {
	key := generateTestKey(t)
	interceptor := authinterceptor.JWTAuthInterceptor(&key.PublicKey, nil)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("other", "value"))
	info := &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/ValidateToken"}
	_, err := interceptor(ctx, nil, info, testHandler)
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestJWTAuthInterceptor_ValidToken(t *testing.T) {
	key := generateTestKey(t)
	token := signToken(t, key, jwt.MapClaims{
		"sub":  "user_123",
		"role": "customer",
	})

	interceptor := authinterceptor.JWTAuthInterceptor(&key.PublicKey, nil)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
	info := &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/ValidateToken"}
	resp, err := interceptor(ctx, nil, info, testHandler)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"status": "ok"}, resp)
}

func TestJWTAuthInterceptor_InvalidToken(t *testing.T) {
	key := generateTestKey(t)
	// Sign with different key
	otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	token := signToken(t, otherKey, jwt.MapClaims{
		"sub":  "user_123",
		"role": "customer",
	})

	interceptor := authinterceptor.JWTAuthInterceptor(&key.PublicKey, nil)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
	info := &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/ValidateToken"}
	_, err = interceptor(ctx, nil, info, testHandler)
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestJWTAuthInterceptor_ExpiredToken(t *testing.T) {
	key := generateTestKey(t)
	token := signToken(t, key, jwt.MapClaims{
		"sub": "user_123",
		"exp": time.Now().Add(-1 * time.Hour).Unix(), // expired 1 hour ago
	})

	interceptor := authinterceptor.JWTAuthInterceptor(&key.PublicKey, nil)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
	info := &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/ValidateToken"}
	_, err := interceptor(ctx, nil, info, testHandler)
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestJWTAuthInterceptor_InvalidSigningMethod(t *testing.T) {
	key := generateTestKey(t)
	// Use HMAC instead of RSA
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user_123",
	})
	signed, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	interceptor := authinterceptor.JWTAuthInterceptor(&key.PublicKey, nil)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+signed))
	info := &grpc.UnaryServerInfo{FullMethod: "/auth.AuthService/ValidateToken"}
	_, err = interceptor(ctx, nil, info, testHandler)
	require.Error(t, err)
}

func TestServiceAuthInterceptor_ValidKey(t *testing.T) {
	interceptor := authinterceptor.ServiceAuthInterceptor([]string{"key1", "key2"})

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-service-api-key", "key1"))
	info := &grpc.UnaryServerInfo{FullMethod: "/customer.CustomerService/GetCustomer"}
	resp, err := interceptor(ctx, nil, info, testHandler)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"status": "ok"}, resp)
}

func TestServiceAuthInterceptor_InvalidKey(t *testing.T) {
	interceptor := authinterceptor.ServiceAuthInterceptor([]string{"key1", "key2"})

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-service-api-key", "evil-key"))
	info := &grpc.UnaryServerInfo{FullMethod: "/customer.CustomerService/GetCustomer"}
	_, err := interceptor(ctx, nil, info, testHandler)
	require.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, st.Code())
}

func TestServiceAuthInterceptor_MissingKey(t *testing.T) {
	interceptor := authinterceptor.ServiceAuthInterceptor([]string{"key1"})

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("other", "value"))
	info := &grpc.UnaryServerInfo{FullMethod: "/customer.CustomerService/GetCustomer"}
	_, err := interceptor(ctx, nil, info, testHandler)
	require.Error(t, err)
}

func TestClientAuthInterceptor_AddsToken(t *testing.T) {
	interceptor := authinterceptor.ClientAuthInterceptor("test-jwt-token")

	var capturedCtx context.Context
	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		capturedCtx = ctx
		return nil
	}

	err := interceptor(context.Background(), "/test.Method", nil, nil, nil, mockInvoker)
	require.NoError(t, err)

	md, ok := metadata.FromOutgoingContext(capturedCtx)
	assert.True(t, ok)
	assert.Equal(t, []string{"Bearer test-jwt-token"}, md["authorization"])
}

func TestServiceClientAuthInterceptor_AddsAPIKey(t *testing.T) {
	interceptor := authinterceptor.ServiceClientAuthInterceptor("svc-api-key")

	var capturedCtx context.Context
	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		capturedCtx = ctx
		return nil
	}

	err := interceptor(context.Background(), "/test.Method", nil, nil, nil, mockInvoker)
	require.NoError(t, err)

	md, ok := metadata.FromOutgoingContext(capturedCtx)
	assert.True(t, ok)
	assert.Equal(t, []string{"svc-api-key"}, md["x-service-api-key"])
}
