//to define the mechanisms of cache

package cache

import "time"

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	GetAll() (map[string]interface{}, error)
}
