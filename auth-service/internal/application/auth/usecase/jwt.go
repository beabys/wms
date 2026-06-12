package usecase

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims for this application.
type Claims struct {
	jwt.RegisteredClaims
	UserID      string   `json:"sub"`
	Role        string   `json:"role"`
	CustomerID  string   `json:"customer_id,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// Service handles JWT token generation and validation using RS256.
type Service struct {
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewService creates a new JWT service with RSA keys.
func NewService(privateKeyPEM, publicKeyPEM string, accessTTL, refreshTTL time.Duration) (*Service, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &Service{
		privateKey:      privateKey,
		publicKey:       publicKey,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}, nil
}

// GenerateAccessToken creates a signed JWT access token for the given user.
func (s *Service) GenerateAccessToken(user *model.User) (string, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			Issuer:    "auth-service",
			Subject:   user.ID,
		},
		UserID:      user.ID,
		Role:        user.Role,
		Permissions: []string{},
	}
	if user.CustomerID != nil {
		claims.CustomerID = *user.CustomerID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return signed, nil
}

// GenerateRefreshToken creates a cryptographically random refresh token.
func (s *Service) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	// Return as hex-encoded string
	return fmt.Sprintf("%x", b), nil
}

// ValidateToken parses and validates a JWT access token, returning the claims.
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// GetPublicKeyPEM returns the public key in PEM format.
func (s *Service) GetPublicKeyPEM() (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(s.publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return string(pem.EncodeToMemory(pemBlock)), nil
}

// GetAccessTokenTTL returns the access token TTL in seconds.
func (s *Service) GetAccessTokenTTL() int64 {
	return int64(s.accessTokenTTL.Seconds())
}

// normalizePEM ensures PEM data has proper newlines (handles single-line env var format).
func normalizePEM(s string) string {
	s = strings.TrimSpace(s)
	// Insert newline after header: "-----BEGIN ... -----" → "-----BEGIN ... -----\n"
	// The header ends at the second "-----" after "BEGIN"
	if strings.HasPrefix(s, "-----BEGIN ") {
		idx := strings.Index(s[len("-----BEGIN "):], "-----")
		if idx >= 0 {
			headerEnd := len("-----BEGIN ") + idx + len("-----")
			s = s[:headerEnd] + "\n" + s[headerEnd:]
		}
	}
	// Insert newline before footer: "-----END ..." → "\n-----END ..."
	footerIdx := strings.Index(s, "-----END ")
	if footerIdx > 0 && s[footerIdx-1] != '\n' {
		s = s[:footerIdx] + "\n" + s[footerIdx:]
	}
	return s
}

func parsePrivateKey(pemData string) (*rsa.PrivateKey, error) {
	pemData = normalizePEM(pemData)
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("no PEM data found")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA private key")
	}
	return rsaKey, nil
}

func parsePublicKey(pemData string) (*rsa.PublicKey, error) {
	// Normalize PEM: ensure newlines after header/before footer (for single-line env var format)
	pemData = normalizePEM(pemData)
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("no PEM data found")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA public key")
	}
	return rsaKey, nil
}

// GenerateKeyPair generates a new RSA key pair for testing.
func GenerateKeyPair() (privatePEM, publicPEM string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key: %w", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %w", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(privPEM), string(pubPEM), nil
}
