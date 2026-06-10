package grpcdapter

import (
	"context"
	"strings"

	"github.com/beabys/wms/auth-service/internal/application/auth/usecase"
	"github.com/beabys/wms/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

// ClaimsFromContext extracts JWT claims from context.
func ClaimsFromContext(ctx context.Context) *usecase.Claims {
	claims, ok := ctx.Value(claimsKey).(*usecase.Claims)
	if !ok {
		return nil
	}
	return claims
}

// AuthInterceptor handles JWT validation for gRPC requests.
type AuthInterceptor struct {
	jwtSvc        *usecase.Service
	exemptMethods []string
	logger        logger.Logger
}

// NewAuthInterceptor creates a new AuthInterceptor.
func NewAuthInterceptor(jwtSvc *usecase.Service, exemptMethods []string, log logger.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		jwtSvc:        jwtSvc,
		exemptMethods: exemptMethods,
		logger:        log,
	}
}

// Unary returns a unary server interceptor that validates JWT tokens.
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Check if method is exempt from auth
		for _, m := range i.exemptMethods {
			if info.FullMethod == m {
				return handler(ctx, req)
			}
		}

		// Extract JWT from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Get token from "authorization" header (format: "Bearer <token>")
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			// Check for service-to-service API key
			apiKey := md.Get("x-service-api-key")
			if len(apiKey) > 0 && apiKey[0] != "" {
				// Service-to-service: skip JWT validation for now
				// In production, validate the API key against a store
				return handler(ctx, req)
			}
			return nil, status.Error(codes.Unauthenticated, "missing authorization")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		if tokenString == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		claims, err := i.jwtSvc.ValidateToken(tokenString)
		if err != nil {
			i.logger.Debug("auth interceptor: invalid token",
				logger.LogField{Key: "method", Value: info.FullMethod},
				logger.LogField{Key: "error", Value: err.Error()},
			)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Inject claims into context
		ctx = context.WithValue(ctx, claimsKey, claims)
		return handler(ctx, req)
	}
}

// GetPublicKeyExemptMethods returns the list of methods that should bypass auth.
func GetPublicKeyExemptMethods() []string {
	return []string{
		"/auth.v1.AuthService/Login",
		"/auth.v1.AuthService/RefreshToken",
		"/auth.v1.AuthService/GetPublicKey",
		"/auth.v1.AuthService/ValidateToken",
		"/auth.v1.AuthService/ValidateInvite",
		"/auth.v1.UserService/CreateUser",
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
	}
}
