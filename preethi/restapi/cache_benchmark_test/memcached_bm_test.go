package tests

import (
	"testing"
	"time"

	"github.com/Preethi0716/Cache-Library/preethi/restapi/pkg/cache"
)

const memcachedAddress = "localhost:11211"

func BenchmarkMemcachedCache_BasicOperations(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached Cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
		_, _ = cache.Get("key")
	}
}

func BenchmarkMemcachedCache_Eviction(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached Cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Millisecond*100)
	}
}

func BenchmarkMemcachedCache_Penetration(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached Cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("nonexistentkey")
	}
}

func BenchmarkMemcachedCache_Expiration(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached Cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Millisecond*10)
		time.Sleep(time.Millisecond * 20)
		_, _ = cache.Get("key" + string(rune(i)))
	}
}

func BenchmarkMemcachedCache_Concurrency(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached Cache: %v", err)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set("key", "value", time.Minute)
			_, _ = cache.Get("key")
		}
	})
}

func BenchmarkMemcachedCache_LargeDataSet(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached Cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Minute)
	}
}

func BenchmarkMemcachedCache_MemoryUsage(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached Cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
		_, _ = cache.Get("key")
	}
}
