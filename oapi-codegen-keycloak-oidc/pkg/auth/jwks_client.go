package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

// JWKSClient is a client for fetching and caching JSON Web Key Sets (JWKS) a.k.a. public keys from a remote URL.
// It provides thread-safe access to public keys identified by key IDs (kid) and automatically
// refreshes the key set after a configurable cache TTL.
//
// Usage:
//
//	client := NewJWKSClient(jwksURL)
//	pubKey, err := client.GetPublicKey(keyID)
//
// Fields:
//   - jwksURL:   The URL from which to fetch the JWKS.
//   - keySet:    The cached set of JWKs.
//   - lastFetch: The timestamp of the last successful fetch.
//   - cacheTTL:  The duration for which the key set is cached.
//   - mutex:     Synchronizes access to the key set.
type JWKSClient struct {
	jwksURL   string
	keySet    jwk.Set
	lastFetch time.Time
	cacheTTL  time.Duration
	mutex     sync.RWMutex
	logger    *slog.Logger
}

const cacheTTLDefault = 5 * time.Minute // Default cache TTL for JWKS keys

// Creates a new JWKSClient with the specified JWKS URL and a default cache TTL of 5 minutes.
func NewJWKSClient(logger *slog.Logger, jwksURL string) *JWKSClient {
	return &JWKSClient{
		jwksURL:  jwksURL,
		cacheTTL: cacheTTLDefault,
		logger:   logger,
	}
}

// Retrieves the RSA public key for the given key ID. If the cache is expired or empty,
// it refreshes the key set from the JWKS endpoint.
func (c *JWKSClient) GetPublicKey(keyID string) (*rsa.PublicKey, error) {
	c.mutex.RLock()
	needsRefresh := c.keySet == nil || time.Since(c.lastFetch) > c.cacheTTL
	c.mutex.RUnlock()

	if needsRefresh {
		if err := c.refreshKeys(); err != nil {
			c.logger.Error("failed to refresh JWKS", "error", err)

			return nil, errors.New("failed to refresh JWKS")
		}
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	key, found := c.keySet.LookupKeyID(keyID)
	if !found {
		c.logger.Error("key with ID not found", "keyID", keyID)

		return nil, errors.New("key with ID not found")
	}

	var rawKey any
	if err := key.Raw(&rawKey); err != nil {
		c.logger.Error("failed to get raw key", "error", err)

		return nil, errors.New("failed to get key")
	}

	rsaKey, ok := rawKey.(*rsa.PublicKey)
	if !ok {
		c.logger.Error("key is not an RSA public key", "keyID", keyID)

		return nil, errors.New("key is not an RSA public key")
	}

	return rsaKey, nil
}

// Fetches the latest JWKS from the configured URL and updates the cache.
func (c *JWKSClient) refreshKeys() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	timeout := 10 * time.Second //nolint: mnd

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	keySet, err := jwk.Fetch(ctx, c.jwksURL)
	if err != nil {
		c.logger.Error("failed to fetch JWKS", "url", c.jwksURL, "error", err)

		return errors.New("failed to fetch JWKS")
	}

	c.keySet = keySet
	c.lastFetch = time.Now()

	return nil
}
