package config

import "time"

type CacheConfig struct {
	RedisAddr        string
	MemcachedServers []string
	MaxLRUSize       int
	DefaultTTL       time.Duration
}
