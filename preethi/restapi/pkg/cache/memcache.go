//Memcached cache syntax to perform cache mechanisms

package cache

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedCache struct {
	client *memcache.Client
}

func NewMemcachedCache(address string) (*MemcachedCache, error) {
	client := memcache.New(address)
	if err := client.Ping(); err != nil {
		return nil, err
	}
	return &MemcachedCache{client: client}, nil
}

func (c *MemcachedCache) Set(key string, value interface{}, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value.(string)),
		Expiration: int32(ttl.Seconds()),
	}
	return c.client.Set(item)
}

func (c *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}
	return string(item.Value), nil
}

func (c *MemcachedCache) Delete(key string) error {
	return c.client.Delete(key)
}

func (c *MemcachedCache) GetAll() (map[string]interface{}, error) {
	// Memcached does not support GetAll in the same way as an in-memory cache.
	return map[string]interface{}{}, nil
}
