package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache represents a Redis cache
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new RedisCache
func NewRedisCache(address string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisCache{client: client}, nil
}

// Set sets a value in the cache with an optional TTL
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(context.Background(), key, value, ttl).Err()
}

// Get gets a value from the cache
func (c *RedisCache) Get(key string) (interface{}, error) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Delete deletes a value from the cache
func (c *RedisCache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

// GetAll retrieves all values from the Redis cache (not generally supported)
func (c *RedisCache) GetAll() (map[string]interface{}, error) {
	// Redis does not support GetAll in the same way as an in-memory cache.
	return map[string]interface{}{}, nil
}
