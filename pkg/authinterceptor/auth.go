package authinterceptor

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	// ClaimsKey is the context key for injected JWT claims.
	ClaimsKey contextKey = "auth.claims"
	// APIKeyHeader is the metadata key for service-to-service auth.
	APIKeyHeader = "x-service-api-key"
	// AuthorizationHeader is the metadata key for user JWT.
	AuthorizationHeader = "authorization"
)

// JWTAuthInterceptor returns a gRPC UnaryServerInterceptor that validates
// JWT tokens on all requests except those in exemptMethods.
func JWTAuthInterceptor(publicKey *rsa.PublicKey, exemptMethods []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for exempt methods
		for _, m := range exemptMethods {
			if info.FullMethod == m {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get(AuthorizationHeader)
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		ctx = context.WithValue(ctx, ClaimsKey, claims)
		return handler(ctx, req)
	}
}

// ServiceAuthInterceptor returns a gRPC UnaryServerInterceptor that validates
// API keys for service-to-service communication.
func ServiceAuthInterceptor(validAPIKeys []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		apiKeys := md.Get(APIKeyHeader)
		if len(apiKeys) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing service api key")
		}

		valid := false
		for _, k := range validAPIKeys {
			if apiKeys[0] == k {
				valid = true
				break
			}
		}

		if !valid {
			return nil, status.Error(codes.PermissionDenied, "invalid service api key")
		}

		ctx = context.WithValue(ctx, ClaimsKey, jwt.MapClaims{
			"service":   true,
			"api_key":   apiKeys[0],
		})
		return handler(ctx, req)
	}
}

// ClientAuthInterceptor adds a JWT to outgoing gRPC metadata.
func ClientAuthInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, AuthorizationHeader, "Bearer "+token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// ServiceClientAuthInterceptor adds an API key to outgoing gRPC metadata.
func ServiceClientAuthInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, APIKeyHeader, apiKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
