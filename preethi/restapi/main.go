package main

import (
	"log"
	"net/http"

	"github.com/Preethi0716/Cache-Library/preethi/restapi/pkg/api"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the caches
	unifiedCache, err := api.InitCache()
	if err != nil {
		log.Fatalf("Failed to initialize caches: %v", err)
	}

	r := mux.NewRouter()

	// Register handlers
	r.HandleFunc("/cache/{key}", api.HandleCacheRequest(unifiedCache)).Methods("GET", "DELETE", "POST")
	r.HandleFunc("/cache", api.HandleGetAllCacheRequest(unifiedCache)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

//Inmemory ::
// post -- http://localhost:8080/cache/d6
// get -- http://localhost:8080/cache/d4?cache=inMemory
// delete -- http://localhost:8080/cache/d7?cache=inMemory

// Redis ::
// post -- http://localhost:8080/cache/d6
// get -- http://localhost:8080/cache/d4?cache=redis
// delete -- http://localhost:8080/cache/d7?cache=redis

// memcached ::
// post -- http://localhost:8080/cache/d6
// get -- http://localhost:8080/cache/d4?cache=memcached
// delete -- http://localhost:8080/cache/d7?cache=memcached
