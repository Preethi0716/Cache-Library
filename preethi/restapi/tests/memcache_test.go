package tests

import (
	"preethi/go/src/preethi/restapi/pkg/cache"
	"testing"
	"time"
)

const memcachedAddress = "localhost:11211"

func TestMemcachedCache_SetGet(t *testing.T) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		t.Fatalf("Failed to create Memcached cache: %v", err)
	}
	cache.Set("key1", "value1", time.Minute)

	value, err := cache.Get("key1")
	if err != nil || value != "value1" {
		t.Fatalf("Expected value1, got %v, error: %v", value, err)
	}
}

func TestMemcachedCache_Delete(t *testing.T) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		t.Fatalf("Failed to create Memcached cache: %v", err)
	}
	cache.Set("key1", "value1", time.Minute)
	cache.Delete("key1")
	_, err = cache.Get("key1")
	if err == nil {
		t.Fatal("Expected an error for a deleted key")
	}
}

func BenchmarkMemcachedCache_Set(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached cache: %v", err)
	}
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
	}
}

func BenchmarkMemcachedCache_Get(b *testing.B) {
	cache, err := cache.NewMemcachedCache(memcachedAddress)
	if err != nil {
		b.Fatalf("Failed to create Memcached cache: %v", err)
	}
	cache.Set("key", "value", time.Minute)
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}
