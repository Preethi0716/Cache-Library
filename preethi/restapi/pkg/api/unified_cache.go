package api

import (
	"fmt"
	"preethi/go/src/preethi/restapi/pkg/cache"
	"time"
)

type UnifiedCache struct {
	InMemoryCache  cache.Cache
	RedisCache     cache.Cache
	MemcachedCache cache.Cache
}

// NewUnifiedCache creates a new UnifiedCache instance.
func NewUnifiedCache(inMemoryCache, redisCache, memcachedCache cache.Cache) *UnifiedCache {
	return &UnifiedCache{
		InMemoryCache:  inMemoryCache,
		RedisCache:     redisCache,
		MemcachedCache: memcachedCache,
	}
}

func (uc *UnifiedCache) Get(key string) (string, error) {
	// Example method to get a value from different caches
	if value, err := uc.InMemoryCache.Get(key); err == nil {
		return value.(string), nil
	}
	if value, err := uc.RedisCache.Get(key); err == nil {
		return value.(string), nil
	}
	if value, err := uc.MemcachedCache.Get(key); err == nil {
		return value.(string), nil
	}
	return "", fmt.Errorf("key not found in any cache")
}

func (uc *UnifiedCache) Set(key string, value string, ttl time.Duration) error {
	// Example method to set a value in different caches
	if err := uc.InMemoryCache.Set(key, value, ttl); err != nil {
		return err
	}
	if err := uc.RedisCache.Set(key, value, ttl); err != nil {
		return err
	}
	if err := uc.MemcachedCache.Set(key, value, ttl); err != nil {
		return err
	}
	return nil
}

func (uc *UnifiedCache) Delete(key string) error {
	// Example method to delete a value from different caches
	if err := uc.InMemoryCache.Delete(key); err != nil {
		return err
	}
	if err := uc.RedisCache.Delete(key); err != nil {
		return err
	}
	if err := uc.MemcachedCache.Delete(key); err != nil {
		return err
	}
	return nil
}

func GetAllCacheEntries(uc *UnifiedCache) (map[string]interface{}, error) {
	entries := make(map[string]interface{})

	// Fetch from in-memory cache
	inMemoryEntries, err := uc.InMemoryCache.GetAll()
	if err != nil {
		return nil, err
	}
	for k, v := range inMemoryEntries {
		entries[k] = v
	}

	// Fetch from Redis cache
	redisEntries, err := uc.RedisCache.GetAll()
	if err != nil {
		return nil, err
	}
	for k, v := range redisEntries {
		entries[k] = v
	}

	// Fetch from Memcached cache
	memcachedEntries, err := uc.MemcachedCache.GetAll()
	if err != nil {
		return nil, err
	}
	for k, v := range memcachedEntries {
		entries[k] = v
	}

	return entries, nil
}
