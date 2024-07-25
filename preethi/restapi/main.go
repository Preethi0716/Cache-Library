package main

import (
	"log"
	"net/http"
	"preethi/go/src/preethi/restapi/pkg/api"

	"github.com/gorilla/mux"
)

func main() {
	unifiedCache, err := api.InitCache()
	if err != nil {
		log.Fatalf("Failed to initialize caches: %v", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/cache/{key}", api.HandleCacheRequest(unifiedCache)).Methods("GET", "POST", "DELETE")
	r.HandleFunc("/cache", api.HandleGetAllCacheRequest(unifiedCache)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

//Inmemory ::
// post -- http://localhost:8080/cache/d6?cache=inMemory
// get -- http://localhost:8080/cache/d4?cache=inMemory
// delete -- http://localhost:8080/cache/d7?cache=inMemory

// Redis ::
// post -- http://localhost:8080/cache/d6?cache=redis
// get -- http://localhost:8080/cache/d4?cache=redis
// delete -- http://localhost:8080/cache/d7?cache=redis

// memcached ::
// post -- http://localhost:8080/cache/d6?cache=memcached
// get -- http://localhost:8080/cache/d4?cache=memcached
// delete -- http://localhost:8080/cache/d7?cache=memcached
