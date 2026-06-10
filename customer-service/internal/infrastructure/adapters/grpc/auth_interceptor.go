package grpcdapter

import (
	"context"
	"strings"

	"github.com/beabys/wms/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

// Claims represents JWT claims extracted from the token.
type Claims struct {
	UserID      string
	Role        string
	CustomerID  string
	Permissions []string
}

// ClaimsFromContext extracts JWT claims from context.
func ClaimsFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}

// AuthInterceptor handles JWT validation for gRPC requests.
type AuthInterceptor struct {
	authClient    *AuthClient
	exemptMethods []string
	logger        logger.Logger
}

// NewAuthInterceptor creates a new AuthInterceptor.
func NewAuthInterceptor(authClient *AuthClient, exemptMethods []string, log logger.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		authClient:    authClient,
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
				return handler(ctx, req)
			}
			return nil, status.Error(codes.Unauthenticated, "missing authorization")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		if tokenString == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		// Validate via auth-service
		validateResult, err := i.authClient.ValidateToken(ctx, tokenString)
		if err != nil {
			i.logger.Debug("auth interceptor: invalid token",
				logger.LogField{Key: "method", Value: info.FullMethod},
				logger.LogField{Key: "error", Value: err.Error()},
			)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Inject claims into context
		claims := &Claims{
			UserID:      validateResult.UserID,
			Role:        validateResult.Role,
			CustomerID:  validateResult.CustomerID,
			Permissions: validateResult.Permissions,
		}
		ctx = context.WithValue(ctx, claimsKey, claims)
		return handler(ctx, req)
	}
}

// GetExemptMethods returns the list of methods that should bypass auth.
func GetExemptMethods() []string {
	return []string{
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
		"/customer.v1.CustomerService/RegisterCustomer",
	}
}
