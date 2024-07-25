package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"preethi/go/src/preethi/restapi/pkg/cache"
	"time"

	"github.com/gorilla/mux"
)

type UnifiedCache struct {
	InMemoryCache  cache.Cache
	RedisCache     cache.Cache
	MemcachedCache cache.Cache
}

func NewUnifiedCache(inMemoryCache, redisCache, memcachedCache cache.Cache) *UnifiedCache {
	return &UnifiedCache{
		InMemoryCache:  inMemoryCache,
		RedisCache:     redisCache,
		MemcachedCache: memcachedCache,
	}
}

func HandleCacheRequest(unifiedCache *UnifiedCache) http.HandlerFunc {
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
			ttl := time.Minute
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

func HandleGetAllCacheRequest(unifiedCache *UnifiedCache) http.HandlerFunc {
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
	var value interface{}
	var err error

	switch cacheType {
	case "inMemory":
		value, err = unifiedCache.InMemoryCache.Get(key)
	case "redis":
		value, err = unifiedCache.RedisCache.Get(key)
	case "memcached":
		value, err = unifiedCache.MemcachedCache.Get(key)
	default:
		return "", fmt.Errorf("invalid cache type")
	}

	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value is not of type string")
	}
	return strValue, nil
}

func setCacheValue(unifiedCache *UnifiedCache, key string, value string, ttl time.Duration, cacheType string) error {
	switch cacheType {
	case "inMemory":
		return unifiedCache.InMemoryCache.Set(key, value, ttl)
	case "redis":
		return unifiedCache.RedisCache.Set(key, value, ttl)
	case "memcached":
		return unifiedCache.MemcachedCache.Set(key, value, ttl)
	default:
		return fmt.Errorf("invalid cache type")
	}
}

func deleteCacheValue(unifiedCache *UnifiedCache, key string, cacheType string) error {
	switch cacheType {
	case "inMemory":
		return unifiedCache.InMemoryCache.Delete(key)
	case "redis":
		return unifiedCache.RedisCache.Delete(key)
	case "memcached":
		return unifiedCache.MemcachedCache.Delete(key)
	default:
		return fmt.Errorf("invalid cache type")
	}
}

func GetAllCacheEntries(unifiedCache *UnifiedCache) (map[string]interface{}, error) {
	allEntries := make(map[string]interface{})

	if unifiedCache.InMemoryCache != nil {
		lruEntries, err := unifiedCache.InMemoryCache.GetAll()
		if err != nil {
			return nil, err
		}
		for k, v := range lruEntries {
			allEntries[k] = v
		}
	}

	if unifiedCache.RedisCache != nil {
		redisEntries, err := unifiedCache.RedisCache.GetAll()
		if err != nil {
			return nil, err
		}
		for k, v := range redisEntries {
			allEntries[k] = v
		}
	}

	if unifiedCache.MemcachedCache != nil {
		memcachedEntries, err := unifiedCache.MemcachedCache.GetAll()
		if err != nil {
			return nil, err
		}
		for k, v := range memcachedEntries {
			allEntries[k] = v
		}
	}

	return allEntries, nil
}
