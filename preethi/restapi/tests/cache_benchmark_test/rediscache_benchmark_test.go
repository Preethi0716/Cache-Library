package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// Define a struct for your Redis cache
type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (rc *RedisCache) Set(key, value string, ttl time.Duration) error {
	return rc.client.Set(context.Background(), key, value, ttl).Err()
}

func (rc *RedisCache) Get(key string) (string, error) {
	return rc.client.Get(context.Background(), key).Result()
}

// Benchmark Basic Set/Get operations
func BenchmarkRedisCache_BasicOperations(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewRedisCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
		_, _ = cache.Get("key")
	}
}

// Benchmark Eviction
func BenchmarkRedisCache_Eviction(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewRedisCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Millisecond*100)
	}
}

// Benchmark Cache Penetration
func BenchmarkRedisCache_Penetration(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewRedisCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("nonexistentkey")
	}
}

// Benchmark Cache Expiration
func BenchmarkRedisCache_Expiration(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewRedisCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Millisecond*10)
		time.Sleep(time.Millisecond * 20) // Ensure expiration
		_, _ = cache.Get("key" + string(rune(i)))
	}
}

// Benchmark Concurrency
func BenchmarkRedisCache_Concurrency(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewRedisCache(client)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set("key", "value", time.Minute)
			_, _ = cache.Get("key")
		}
	})
}

// Benchmark Large Data Sets
func BenchmarkRedisCache_LargeDataSet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewRedisCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Minute)
	}
}

// Benchmark Memory Usage
func BenchmarkRedisCache_MemoryUsage(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewRedisCache(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
		_, _ = cache.Get("key")
	}
}

func BenchmarkAllRedisOperations(b *testing.B) {
	b.Run("BasicOperations", BenchmarkRedisCache_BasicOperations)
	b.Run("Eviction", BenchmarkRedisCache_Eviction)
	b.Run("Penetration", BenchmarkRedisCache_Penetration)
	b.Run("Expiration", BenchmarkRedisCache_Expiration)
	b.Run("Concurrency", BenchmarkRedisCache_Concurrency)
	b.Run("LargeDataSet", BenchmarkRedisCache_LargeDataSet)
	b.Run("MemoryUsage", BenchmarkRedisCache_MemoryUsage)
}
