package cache

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedCache represents a Memcached cache
type MemcachedCache struct {
	client *memcache.Client
}

// NewMemcachedCache creates a new MemcachedCache
func NewMemcachedCache(address string) (*MemcachedCache, error) {
	client := memcache.New(address)
	if err := client.Ping(); err != nil {
		return nil, err
	}
	return &MemcachedCache{client: client}, nil
}

// Set sets a value in the cache with an optional TTL
func (c *MemcachedCache) Set(key string, value interface{}, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value.(string)),
		Expiration: int32(ttl.Seconds()),
	}
	return c.client.Set(item)
}

// Get gets a value from the cache
func (c *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}
	return string(item.Value), nil
}

// Delete deletes a value from the cache
func (c *MemcachedCache) Delete(key string) error {
	return c.client.Delete(key)
}

// GetAll retrieves all values from the Memcached cache (not generally supported)
func (c *MemcachedCache) GetAll() (map[string]interface{}, error) {
	// Memcached does not support GetAll in the same way as an in-memory cache.
	return map[string]interface{}{}, nil
}
