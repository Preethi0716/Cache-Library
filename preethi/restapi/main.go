package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"preethi/go/src/preethi/restapi/pkg/api"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	unifiedCache, err := api.InitCache()
	if err != nil {
		log.Fatalf("Failed to initialize caches: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/cache/{key}", handleCacheRequest(unifiedCache)).Methods("GET", "POST", "DELETE")
	r.HandleFunc("/cache", handleGetAllCacheRequest(unifiedCache)).Methods("GET")

	log.Println("Server is starting on port 8080...")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCacheRequest(unifiedCache *api.UnifiedCache) http.HandlerFunc {
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

func handleGetAllCacheRequest(unifiedCache *api.UnifiedCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allEntries, err := api.GetAllCacheEntries(unifiedCache)
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

func getCacheValue(unifiedCache *api.UnifiedCache, key string, cacheType string) (string, error) {
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

	// Type assert the value to a string
	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value is not of type string")
	}
	return strValue, nil
}

func setCacheValue(unifiedCache *api.UnifiedCache, key string, value string, ttl time.Duration, cacheType string) error {
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

func deleteCacheValue(unifiedCache *api.UnifiedCache, key string, cacheType string) error {
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
