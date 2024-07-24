package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func handleCacheRequest(unifiedCache *UnifiedCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		cacheType := r.URL.Query().Get("cache")

		switch r.Method {
		case "GET":
			value, err := getCacheValue(unifiedCache, key, cacheType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Write([]byte(value))
		case "POST":
			var requestBody map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			value, ok := requestBody["value"].(string)
			if !ok {
				http.Error(w, "Invalid value format", http.StatusBadRequest)
				return
			}
			ttl := time.Minute // Example TTL
			err := setCacheValue(unifiedCache, key, value, ttl, cacheType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		case "DELETE":
			err := deleteCacheValue(unifiedCache, key, cacheType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleGetAllCacheRequest(unifiedCache *UnifiedCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allEntries, err := GetAllCacheEntries(unifiedCache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(allEntries)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func getCacheValue(unifiedCache *UnifiedCache, key string, cacheType string) (string, error) {
	switch cacheType {
	case "redis":
		value, err := unifiedCache.RedisCache.Get(key)
		if err != nil {
			return "", err
		}
		return value.(string), nil
	case "memcached":
		value, err := unifiedCache.MemcachedCache.Get(key)
		if err != nil {
			return "", err
		}
		return value.(string), nil
	case "lru":
		value, err := unifiedCache.InMemoryCache.Get(key)
		if err != nil {
			return "", err
		}
		return value.(string), nil
	default:
		return "", fmt.Errorf("unknown cache type: %s", cacheType)
	}
}

func setCacheValue(unifiedCache *UnifiedCache, key string, value string, ttl time.Duration, cacheType string) error {
	switch cacheType {
	case "redis":
		return unifiedCache.RedisCache.Set(key, value, ttl)
	case "memcached":
		return unifiedCache.MemcachedCache.Set(key, value, ttl)
	case "lru":
		return unifiedCache.InMemoryCache.Set(key, value, ttl)
	default:
		return fmt.Errorf("unknown cache type: %s", cacheType)
	}
}

func deleteCacheValue(unifiedCache *UnifiedCache, key string, cacheType string) error {
	switch cacheType {
	case "redis":
		return unifiedCache.RedisCache.Delete(key)
	case "memcached":
		return unifiedCache.MemcachedCache.Delete(key)
	case "lru":
		return unifiedCache.InMemoryCache.Delete(key)
	default:
		return fmt.Errorf("unknown cache type: %s", cacheType)
	}
}
