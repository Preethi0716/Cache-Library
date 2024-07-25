package cache

import (
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Define a struct for your Memcached cache
type MemcachedCache struct {
	client *memcache.Client
}

func NewMemcachedCache(client *memcache.Client) *MemcachedCache {
	return &MemcachedCache{client: client}
}

func (mc *MemcachedCache) Set(key, value string, ttl time.Duration) error {
	return mc.client.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: int32(ttl.Seconds())})
}

func (mc *MemcachedCache) Get(key string) (string, error) {
	item, err := mc.client.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

func (mc *MemcachedCache) Delete(key string) error {
	return mc.client.Delete(key)
}

// Benchmark Basic Set/Get operations
func BenchmarkMemcachedCache_BasicOperations(b *testing.B) {
	client := memcache.New("localhost:11211")
	cache := NewMemcachedCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
		_, _ = cache.Get("key")
	}
}

// Benchmark Eviction
func BenchmarkMemcachedCache_Eviction(b *testing.B) {
	client := memcache.New("localhost:11211")
	cache := NewMemcachedCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(i), "value", time.Millisecond*100)
	}
}

// Benchmark Cache Penetration
func BenchmarkMemcachedCache_Penetration(b *testing.B) {
	client := memcache.New("localhost:11211")
	cache := NewMemcachedCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("nonexistentkey")
	}
}

// Benchmark Cache Expiration
func BenchmarkMemcachedCache_Expiration(b *testing.B) {
	client := memcache.New("localhost:11211")
	cache := NewMemcachedCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(i), "value", time.Millisecond*10)
		time.Sleep(time.Millisecond * 20) // Ensure expiration
		_, _ = cache.Get("key" + string(i))
	}
}

// Benchmark Concurrency
func BenchmarkMemcachedCache_Concurrency(b *testing.B) {
	client := memcache.New("localhost:11211")
	cache := NewMemcachedCache(client)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set("key", "value", time.Minute)
			_, _ = cache.Get("key")
		}
	})
}

// Benchmark Large Data Sets
func BenchmarkMemcachedCache_LargeDataSet(b *testing.B) {
	client := memcache.New("localhost:11211")
	cache := NewMemcachedCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(i), "value", time.Minute)
	}
}

// Benchmark Memory Usage
func BenchmarkMemcachedCache_MemoryUsage(b *testing.B) {
	client := memcache.New("localhost:11211")
	cache := NewMemcachedCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
		_, _ = cache.Get("key")
	}
}

func BenchmarkAllMemcachedOperations(b *testing.B) {
	b.Run("BasicOperations", BenchmarkMemcachedCache_BasicOperations)
	b.Run("Eviction", BenchmarkMemcachedCache_Eviction)
	b.Run("Penetration", BenchmarkMemcachedCache_Penetration)
	b.Run("Expiration", BenchmarkMemcachedCache_Expiration)
	b.Run("Concurrency", BenchmarkMemcachedCache_Concurrency)
	b.Run("LargeDataSet", BenchmarkMemcachedCache_LargeDataSet)
	b.Run("MemoryUsage", BenchmarkMemcachedCache_MemoryUsage)
}
