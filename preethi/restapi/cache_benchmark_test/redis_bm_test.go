package tests

import (
	"testing"
	"time"

	"github.com/Preethi0716/Cache-Library/preethi/restapi/pkg/cache"
)

func BenchmarkRedisCache_BasicOperations(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("failed to create Redis cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Set("key", "value", time.Minute)
		_, _ = client.Get("key")
	}
}

func BenchmarkRedisCache_Eviction(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("failed to create Redis cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Set("key"+string(rune(i)), "value", time.Millisecond*100)
	}
}

func BenchmarkRedisCache_Penetration(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("failed to create Redis cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = client.Get("nonexistentkey")
	}
}

func BenchmarkRedisCache_Expiration(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("failed to create Redis cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Set("key"+string(rune(i)), "value", time.Millisecond*10)
		time.Sleep(time.Millisecond * 20)
		_, _ = client.Get("key" + string(rune(i)))
	}
}

func BenchmarkRedisCache_Concurrency(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("failed to create Redis cache: %v", err)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client.Set("key", "value", time.Minute)
			_, _ = client.Get("key")
		}
	})
}

func BenchmarkRedisCache_LargeDataSet(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("failed to create Redis cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Set("key"+string(rune(i)), "value", time.Minute)
	}
}

func BenchmarkRedisCache_MemoryUsage(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("failed to create Redis cache: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Set("key", "value", time.Minute)
		_, _ = client.Get("key")
	}
}
