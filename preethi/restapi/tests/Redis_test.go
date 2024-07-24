package tests

import (
	"preethi/go/src/preethi/restapi/pkg/cache"
	"testing"
	"time"
)

func TestRedisCache_SetGet(t *testing.T) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create Redis cache: %v", err)
	}

	err = client.Set("key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	value, err := client.Get("key1")
	if err != nil || value != "value1" {
		t.Fatalf("Expected value1, got %v", value)
	}
}

func TestRedisCache_Delete(t *testing.T) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create Redis cache: %v", err)
	}

	err = client.Set("key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	err = client.Delete("key1")
	if err != nil {
		t.Fatalf("Failed to delete cache: %v", err)
	}

	_, err = client.Get("key1")
	if err == nil {
		t.Fatal("Expected an error for a deleted key")
	}
}

func BenchmarkRedisCache_Set(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("Failed to create Redis cache: %v", err)
	}

	for i := 0; i < b.N; i++ {
		err := client.Set("key", "value", time.Minute)
		if err != nil {
			b.Fatalf("Failed to set cache: %v", err)
		}
	}
}

func BenchmarkRedisCache_Get(b *testing.B) {
	client, err := cache.NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("Failed to create Redis cache: %v", err)
	}

	err = client.Set("key", "value", time.Minute)
	if err != nil {
		b.Fatalf("Failed to set cache: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := client.Get("key")
		if err != nil {
			b.Fatalf("Failed to get cache: %v", err)
		}
	}
}
