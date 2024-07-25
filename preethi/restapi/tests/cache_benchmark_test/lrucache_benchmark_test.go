package cache

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/Preethi0716/Cache-Library/preethi/restapi/pkg/cache"
)

func BenchmarkLRUCache(b *testing.B) {
	// Measure memory usage before benchmark
	var memStatsBefore runtime.MemStats
	runtime.ReadMemStats(&memStatsBefore)

	// Basic Set/Get Benchmark
	b.Run("Basic Set/Get", func(b *testing.B) {
		cache := cache.NewLRUCache(100)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Set(fmt.Sprintf("key%d", i), "value", time.Minute)
			cache.Get(fmt.Sprintf("key%d", i))
		}
	})

	// Eviction Benchmark
	b.Run("Eviction", func(b *testing.B) {
		cache := cache.NewLRUCache(100)
		b.ResetTimer()
		for i := 0; i < 200; i++ {
			cache.Set(fmt.Sprintf("key%d", i), "value", time.Minute)
		}
		// Ensure some gets to check eviction
		for i := 0; i < 100; i++ {
			cache.Get(fmt.Sprintf("key%d", i))
		}
	})

	// Cache Penetration Benchmark
	b.Run("Cache Penetration", func(b *testing.B) {
		cache := cache.NewLRUCache(100)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Get(fmt.Sprintf("key%d", i))
		}
	})

	// Cache Expiration Benchmark
	b.Run("Cache Expiration", func(b *testing.B) {
		cache := cache.NewLRUCache(100)
		cache.Set("key1", "value", time.Millisecond*100)
		time.Sleep(time.Millisecond * 200)
		b.ResetTimer()
		cache.Get("key1") // Should be expired
	})

	// Concurrency Benchmark
	b.Run("Concurrency", func(b *testing.B) {
		cache := cache.NewLRUCache(100)
		var wg sync.WaitGroup
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				cache.Set(fmt.Sprintf("key%d", i), "value", time.Minute)
				cache.Get(fmt.Sprintf("key%d", i))
			}(i)
		}
		wg.Wait()
	})

	// Large Data Sets Benchmark
	b.Run("Large Data Sets", func(b *testing.B) {
		cache := cache.NewLRUCache(10000) // Larger cache size
		for i := 0; i < 10000; i++ {
			cache.Set(fmt.Sprintf("key%d", i), "value", time.Minute)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Get(fmt.Sprintf("key%d", i%10000))
		}
	})

	// Measure memory usage after benchmark
	var memStatsAfter runtime.MemStats
	runtime.ReadMemStats(&memStatsAfter)
	fmt.Printf("Memory Usage: Before: %v, After: %v, Difference: %v\n",
		memStatsBefore.Alloc, memStatsAfter.Alloc, memStatsAfter.Alloc-memStatsBefore.Alloc)
}
