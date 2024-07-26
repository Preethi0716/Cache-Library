package tests

import (
	"testing"
	"time"

	"github.com/Preethi0716/Cache-Library/preethi/restapi/pkg/cache"
)

func BenchmarkLRUCache_BasicOperations(b *testing.B) {
	cache := cache.NewLRUCache(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
		_, _ = cache.Get("key")
	}
}

func BenchmarkLRUCache_Eviction(b *testing.B) {
	cache := cache.NewLRUCache(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Millisecond*100)
	}
}

func BenchmarkLRUCache_Penetration(b *testing.B) {
	cache := cache.NewLRUCache(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("nonexistentkey")
	}
}

func BenchmarkLRUCache_Expiration(b *testing.B) {
	cache := cache.NewLRUCache(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Millisecond*10)
		time.Sleep(time.Millisecond * 20)
		_, _ = cache.Get("key" + string(rune(i)))
	}
}

func BenchmarkLRUCache_Concurrency(b *testing.B) {
	cache := cache.NewLRUCache(1000)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set("key", "value", time.Minute)
			_, _ = cache.Get("key")
		}
	})
}

func BenchmarkLRUCache_LargeDataSet(b *testing.B) {
	cache := cache.NewLRUCache(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), "value", time.Minute)
	}
}

func BenchmarkLRUCache_MemoryUsage(b *testing.B) {
	cache := cache.NewLRUCache(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
		_, _ = cache.Get("key")
	}
}
