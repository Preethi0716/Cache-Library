package tests

import (
	"preethi/go/src/preethi/restapi/pkg/cache"
	"testing"
	"time"
)

func TestLRUCache_SetGet(t *testing.T) {
	cache := cache.NewLRUCache(2)
	cache.Set("key1", "value1", time.Minute)

	value, err := cache.Get("key1")
	if err != nil || value != "value1" {
		t.Fatalf("Expected value1, got %v", value)
	}
}

func TestLRUCache_Delete(t *testing.T) {
	cache := cache.NewLRUCache(2)
	cache.Set("key1", "value1", time.Minute)
	cache.Delete("key1")
	_, err := cache.Get("key1")
	if err == nil {
		t.Fatal("Expected an error for a deleted key")
	}
}

func BenchmarkLRUCache_Set(b *testing.B) {
	cache := cache.NewLRUCache(2)
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
	}
}

func BenchmarkLRUCache_Get(b *testing.B) {
	cache := cache.NewLRUCache(2)
	cache.Set("key", "value", time.Minute)
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}
