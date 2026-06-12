package grpcdapter

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/usecase"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func setupInterceptorTest(t *testing.T) (*AuthInterceptor, *usecase.Service) {
	t.Helper()
	privPEM, pubPEM, err := usecase.GenerateKeyPair()
	require.NoError(t, err)

	jwtSvc, err := usecase.NewService(privPEM, pubPEM, 15*time.Minute, 7*24*time.Hour)
	require.NoError(t, err)

	log := logger.NewSlogLogger(slog.LevelDebug)
	exemptMethods := GetPublicKeyExemptMethods()
	interceptor := NewAuthInterceptor(jwtSvc, exemptMethods, log)
	return interceptor, jwtSvc
}

func TestAuthInterceptor_ExemptMethods(t *testing.T) {
	t.Run("login method is exempt", func(t *testing.T) {
		interceptor, _ := setupInterceptorTest(t)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		result, err := interceptor.Unary()(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.AuthService/Login",
		}, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("refresh token method is exempt", func(t *testing.T) {
		interceptor, _ := setupInterceptorTest(t)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		result, err := interceptor.Unary()(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.AuthService/RefreshToken",
		}, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("validate invite method is exempt", func(t *testing.T) {
		interceptor, _ := setupInterceptorTest(t)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		result, err := interceptor.Unary()(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.AuthService/ValidateInvite",
		}, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("health check method is exempt", func(t *testing.T) {
		interceptor, _ := setupInterceptorTest(t)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		result, err := interceptor.Unary()(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/grpc.health.v1.Health/Check",
		}, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", result)
	})
}

func TestAuthInterceptor_NoAuth(t *testing.T) {
	t.Run("missing metadata returns error", func(t *testing.T) {
		interceptor, _ := setupInterceptorTest(t)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		// Context without metadata
		_, err := interceptor.Unary()(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.UserService/AssignRole",
		}, handler)
		require.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("missing authorization header returns error", func(t *testing.T) {
		interceptor, _ := setupInterceptorTest(t)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("some-key", "some-value"))
		_, err := interceptor.Unary()(ctx, nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.UserService/AssignRole",
		}, handler)
		require.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})
}

func TestAuthInterceptor_ValidToken(t *testing.T) {
	t.Run("valid token passes through", func(t *testing.T) {
		interceptor, jwtSvc := setupInterceptorTest(t)

		user := &model.User{ID: "user-1", Role: "admin", Email: "test@example.com", Name: "Test"}
		token, err := jwtSvc.GenerateAccessToken(user)
		require.NoError(t, err)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			claims := ClaimsFromContext(ctx)
			require.NotNil(t, claims)
			assert.Equal(t, "user-1", claims.UserID)
			assert.Equal(t, "admin", claims.Role)
			return "success", nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		result, err := interceptor.Unary()(ctx, nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.UserService/GetUser",
		}, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("valid service API key passes through", func(t *testing.T) {
		interceptor, _ := setupInterceptorTest(t)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-service-api-key", "some-api-key"))
		result, err := interceptor.Unary()(ctx, nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.UserService/AssignRole",
		}, handler)
		require.NoError(t, err)
		assert.Equal(t, "success", result)
	})
}

func TestAuthInterceptor_InvalidToken(t *testing.T) {
	t.Run("invalid token returns error", func(t *testing.T) {
		interceptor, _ := setupInterceptorTest(t)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer invalid-token"))
		_, err := interceptor.Unary()(ctx, nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.UserService/GetUser",
		}, handler)
		require.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("expired token returns error", func(t *testing.T) {
		privPEM, pubPEM, err := usecase.GenerateKeyPair()
		require.NoError(t, err)

		// Create service with -1 minute TTL
		jwtSvc, err := usecase.NewService(privPEM, pubPEM, -1*time.Minute, 7*24*time.Hour)
		require.NoError(t, err)

		log := logger.NewSlogLogger(slog.LevelDebug)
		interceptor := NewAuthInterceptor(jwtSvc, GetPublicKeyExemptMethods(), log)

		user := &model.User{ID: "user-1", Role: "admin"}
		token, err := jwtSvc.GenerateAccessToken(user)
		require.NoError(t, err)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "success", nil
		}

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		_, err = interceptor.Unary()(ctx, nil, &grpc.UnaryServerInfo{
			FullMethod: "/auth.v1.UserService/DeleteUser",
		}, handler)
		require.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})
}

func TestAuthInterceptor_ClaimsFromContext(t *testing.T) {
	t.Run("nil when no claims in context", func(t *testing.T) {
		claims := ClaimsFromContext(context.Background())
		assert.Nil(t, claims)
	})
}
