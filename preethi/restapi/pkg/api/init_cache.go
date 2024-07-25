// for initializing the cache

package api

import (
	"fmt"
	"preethi/go/src/preethi/restapi/pkg/cache"
)

func InitCache() (*UnifiedCache, error) {
	inMemoryCache := cache.NewLRUCache(5)
	if inMemoryCache == nil {
		return nil, fmt.Errorf("failed to initialize in-memory cache")
	}

	redisCache, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis cache: %w", err)
	}

	memcachedCache, err := cache.NewMemcachedCache("localhost:11211")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Memcached cache: %w", err)
	}

	return NewUnifiedCache(inMemoryCache, redisCache, memcachedCache), nil
}
