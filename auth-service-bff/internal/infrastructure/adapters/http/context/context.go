package httpctx

// ContextKey is the type used for context keys in request-scoped values.
// Using a dedicated type prevents collisions with context keys from other packages.
type ContextKey string

const (
	// ContextKeyJWT is the context key for the JWT token extracted from the Authorization header.
	ContextKeyJWT ContextKey = "jwt_token"
)
