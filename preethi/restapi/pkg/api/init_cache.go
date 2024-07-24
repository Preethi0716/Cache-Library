package api

import (
	"fmt"
	"preethi/go/src/preethi/restapi/pkg/cache"
)

// InitCache initializes the caches and returns a UnifiedCache instance
func InitCache() (*UnifiedCache, error) {
	// Initialize LRU Cache with size 100
	inMemoryCache := cache.NewLRUCache(5)
	if inMemoryCache == nil {
		return nil, fmt.Errorf("failed to initialize in-memory cache")
	}

	// Initialize Redis Cache with address
	redisCache, err := cache.NewRedisCache("localhost:6379") // Adjust the address as needed
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis cache: %w", err)
	}

	// Initialize Memcached Cache with address
	memcachedCache, err := cache.NewMemcachedCache("localhost:11211") // Adjust the address as needed
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Memcached cache: %w", err)
	}

	return NewUnifiedCache(inMemoryCache, redisCache, memcachedCache), nil
}
