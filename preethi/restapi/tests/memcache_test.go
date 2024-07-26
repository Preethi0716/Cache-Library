package tests

import (
	"testing"
	"time"

	"github.com/Preethi0716/Cache-Library/preethi/restapi/pkg/cache"
)

func TestMemcachedCache_SetGet(t *testing.T) {
	cache, err := cache.NewMemcachedCache("localhost:11211")
	if err != nil {
		t.Fatalf("Failed to create Memcached cache: %v", err)
	}

	err = cache.Set("key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	value, err := cache.Get("key1")
	if err != nil || value != "value1" {
		t.Fatalf("Expected value1, got %v", value)
	}
}

func TestMemcachedCache_Delete(t *testing.T) {
	cache, err := cache.NewMemcachedCache("localhost:11211")
	if err != nil {
		t.Fatalf("Failed to create Memcached cache: %v", err)
	}

	err = cache.Set("key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	err = cache.Delete("key1")
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	_, err = cache.Get("key1")
	if err == nil {
		t.Fatal("Expected an error for a deleted key")
	}
}

func TestMemcachedCache_TTL_Expiration(t *testing.T) {
	cache, err := cache.NewMemcachedCache("localhost:11211")
	if err != nil {
		t.Fatalf("Failed to create Memcached cache: %v", err)
	}

	err = cache.Set("key1", "value1", 1*time.Second)
	if err != nil {
		t.Fatalf("Failed to set value with TTL: %v", err)
	}

	value, err := cache.Get("key1")
	if err != nil || value != "value1" {
		t.Fatalf("Expected value1, got %v, error: %v", value, err)
	}

	time.Sleep(2 * time.Second)

	value, err = cache.Get("key1")

	if err != nil && err.Error() != "memcache: cache miss" {
		t.Fatalf("Expected cache miss error, got: %v", err)
	}

	if value != nil {
		t.Fatalf("Expected empty value, got: %v", value)
	}
}
