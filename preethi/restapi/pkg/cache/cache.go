// pkg/cache/cache.go
package cache

import "time"

// Cache interface defines methods for caching
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	GetAll() (map[string]interface{}, error)
}
