package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Preethi0716/Cache-Library/preethi/restapi/pkg/cache"
)

func TestRedisCache_SetGet(t *testing.T) {
	cache, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create Redis cache: %v", err)
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

func TestRedisCache_Delete(t *testing.T) {
	cache, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create Redis cache: %v", err)
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

func TestRedisCache_Concurrency(t *testing.T) {
	c, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create Redis cache: %v", err)
	}

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			c.Set(key, "value", time.Minute)
			if _, err := c.Get(key); err != nil {
				t.Fatalf("Expected to get %v", key)
			}
		}(i)
	}
	wg.Wait()
}

func TestRedisCache_ConcurrencyWithTTL(t *testing.T) {
	cache, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create Redis cache: %v", err)
	}

	var wg sync.WaitGroup

	const numOperations = 1000
	const ttl = 1 * time.Second

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			err := cache.Set(key, "value", ttl)
			if err != nil {
				t.Errorf("Failed to set %v: %v", key, err)
				return
			}

			value, err := cache.Get(key)
			if err != nil || value != "value" {
				t.Errorf("Expected value, got %v for key %v, error: %v", value, key, err)
			}

			time.Sleep(ttl + 1*time.Second)

			value, err = cache.Get(key)
			if err == nil {
				t.Errorf("Expected cache miss error for key %v after TTL expiration, got: %v", key, value)
			}
			if value != nil {
				t.Errorf("Expected nil value for key %v after TTL expiration, got: %v", key, value)
			}
		}(i)
	}
	wg.Wait()
}

func TestRedisCache_UpdateValue(t *testing.T) {
	cache, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create Redis cache: %v", err)
	}

	err = cache.Set("key1", "initialValue", time.Minute)
	if err != nil {
		t.Fatalf("Failed to set initial value: %v", err)
	}

	value, err := cache.Get("key1")
	if err != nil || value != "initialValue" {
		t.Fatalf("Expected initialValue, got %v", value)
	}

	err = cache.Set("key1", "updatedValue", time.Minute)
	if err != nil {
		t.Fatalf("Failed to update value: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	value, err = cache.Get("key1")
	if err != nil || value != "updatedValue" {
		t.Fatalf("Expected updatedValue, got %v", value)
	}
}
