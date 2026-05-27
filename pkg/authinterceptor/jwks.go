package authinterceptor

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// JWKSClient fetches and caches RSA public keys from a JWKS endpoint.
type JWKSClient struct {
	BaseURL    string
	HTTPClient *http.Client

	mu          sync.RWMutex
	cachedKey   *rsa.PublicKey
	cachedAt    time.Time
	cacheTTL    time.Duration
}

// NewJWKSClient creates a new JWKSClient.
func NewJWKSClient(baseURL string) *JWKSClient {
	return &JWKSClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cacheTTL: 5 * time.Minute,
	}
}

// GetPublicKey returns the cached RSA public key, refreshing from the remote
// endpoint if the cache is stale.
func (c *JWKSClient) GetPublicKey(ctx context.Context) (*rsa.PublicKey, error) {
	c.mu.RLock()
	if c.cachedKey != nil && time.Since(c.cachedAt) < c.cacheTTL {
		defer c.mu.RUnlock()
		return c.cachedKey, nil
	}
	c.mu.RUnlock()

	return c.refresh(ctx)
}

// refresh fetches the public key from the remote endpoint and updates the cache.
func (c *JWKSClient) refresh(ctx context.Context) (*rsa.PublicKey, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.cachedKey != nil && time.Since(c.cachedAt) < c.cacheTTL {
		return c.cachedKey, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/v1/keys", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch keys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Expect PEM-encoded public key in response body
	buf := make([]byte, 4096)
	n, _ := resp.Body.Read(buf)
	block, _ := pem.Decode(buf[:n])
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA")
	}

	c.cachedKey = rsaKey
	c.cachedAt = time.Now()

	return rsaKey, nil
}
