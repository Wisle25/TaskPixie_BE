package cache

import "time"

// Cache interface defines methods for setting and getting cache values.
// This interface is used when the use case wants to cache something
type Cache interface {
	// SetCache sets a value in the cache with an expiration duration.
	SetCache(key string, value interface{}, expiration time.Duration)

	// GetCache retrieves a value from the cache by key.
	GetCache(key string) interface{}

	// DeleteCache removing cache
	DeleteCache(key string)
}
