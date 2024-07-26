//go test ./tests/ -v

package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Preethi0716/Cache-Library/preethi/restapi/pkg/cache"
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

func TestLRUCache_Eviction(t *testing.T) {
	cache := cache.NewLRUCache(2)
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	cache.Set("key3", "value3", time.Minute)

	value, err := cache.Get("key1")
	if err == nil {
		t.Fatal("Expected an error for an evicted key")
	}

	value, err = cache.Get("key2")
	if err != nil || value != "value2" {
		t.Fatalf("Expected value2, got %v", value)
	}

	value, err = cache.Get("key3")
	if err != nil || value != "value3" {
		t.Fatalf("Expected value3, got %v", value)
	}
}

func TestLRUCache_Concurrency(t *testing.T) {
	capacity := 1000
	cache := cache.NewLRUCache(capacity)

	var wg sync.WaitGroup
	var mu sync.Mutex
	failedKeys := []string{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			err := cache.Set(key, "value", time.Minute)
			if err != nil {
				mu.Lock()
				failedKeys = append(failedKeys, key)
				mu.Unlock()
				t.Errorf("Failed to set %v: %v", key, err)
				return
			}
			value, err := cache.Get(key)
			if err != nil {
				mu.Lock()
				failedKeys = append(failedKeys, key)
				mu.Unlock()
				t.Errorf("Failed to get %v: %v", key, err)
				return
			}
			expectedValue := "value"
			if value != expectedValue {
				mu.Lock()
				failedKeys = append(failedKeys, key)
				mu.Unlock()
				t.Errorf("Expected %v, got %v for %v", expectedValue, value, key)
			}
		}(i)
	}
	wg.Wait()

	if len(failedKeys) > 0 {
		t.Fatalf("Failed to get keys: %v", failedKeys)
	}
}

func TestLRUCache_TTL_Stress(t *testing.T) {
	capacity := 100
	cache := cache.NewLRUCache(capacity)

	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 10*time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)

	for i := 0; i < 1000; i++ {
		_, err := cache.Get(fmt.Sprintf("key%d", i))
		if err == nil {
			t.Fatalf("Expected an error for expired key%d", i)
		}
	}
}

func TestLRUCache_UpdateValue(t *testing.T) {
	cache := cache.NewLRUCache(2)

	err := cache.Set("key1", "initialValue", time.Minute)
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

	value, err = cache.Get("key1")
	if err != nil || value != "updatedValue" {
		t.Fatalf("Expected updatedValue, got %v", value)
	}
}

func TestLRUCache_DifferentPorts(t *testing.T) {
	ports := []int{11211, 6379, 8080}

	for _, port := range ports {
		t.Run(fmt.Sprintf("port%d", port), func(t *testing.T) {
			cache := cache.NewLRUCache(2)

			err := cache.Set("key1", "value1", time.Minute)
			if err != nil {
				t.Fatalf("Failed to set value on port %d: %v", port, err)
			}

			value, err := cache.Get("key1")
			if err != nil || value != "value1" {
				t.Fatalf("Expected value1, got %v on port %d", value, port)
			}
		})
	}
}
