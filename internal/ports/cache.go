package ports

import "context"

// CacheService defines the interface for caching operations
type CacheService interface {
	// Set stores a value in cache
	Set(ctx context.Context, key string, value interface{}, ttl int) error

	// Get retrieves a value from cache
	Get(ctx context.Context, key string, dest interface{}) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// Close closes the cache connection
	Close() error
}
