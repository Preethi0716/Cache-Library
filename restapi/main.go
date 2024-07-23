package main

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Cache structure to hold user data with timestamp
type Cache struct {
	ID        uint      `gorm:"primaryKey"`
	User      string    `gorm:"unique;not null" json:"UserId"`
	Password  string    `json:"Password"`
	Number    string    `json:"Access Number"`
	CreatedAt time.Time `json:"Created At"` // Timestamp field
}

// CacheClient interface for cache operations
type CacheClient interface {
	Get(key string) (Cache, bool)
	Set(key string, value Cache)
	Print() // Include the Print method
}

// LRUCache item structure with timestamp
type cacheItem struct {
	key       string
	value     Cache
	timestamp time.Time
}

// LRUCache structure
type LRUCache struct {
	capacity     int
	ttl          time.Duration
	items        map[string]*list.Element
	evictionList *list.List
	mu           sync.Mutex
}

// Initialize LRUCache
func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	return &LRUCache{
		capacity:     capacity,
		ttl:          ttl,
		items:        make(map[string]*list.Element),
		evictionList: list.New(),
	}
}

// Get an item from LRUCache
func (c *LRUCache) Get(key string) (Cache, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.items[key]; found {
		item := element.Value.(*cacheItem)
		if time.Since(item.timestamp) < c.ttl {
			item.timestamp = time.Now() // Update the timestamp
			c.evictionList.MoveToFront(element)
			return item.value, true
		}
		c.evictionList.Remove(element)
		delete(c.items, key)
	}
	return Cache{}, false
}

// Set an item in LRUCache
func (c *LRUCache) Set(key string, value Cache) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.items[key]; found {
		c.evictionList.MoveToFront(element)
		item := element.Value.(*cacheItem)
		item.value = value
		item.timestamp = time.Now()
	} else {
		if c.evictionList.Len() >= c.capacity {
			backElement := c.evictionList.Back()
			if backElement != nil {
				c.evictionList.Remove(backElement)
				delete(c.items, backElement.Value.(*cacheItem).key)
			}
		}
		item := &cacheItem{
			key:       key,
			value:     value,
			timestamp: time.Now(),
		}
		frontElement := c.evictionList.PushFront(item)
		c.items[key] = frontElement
	}
}

func (c *LRUCache) Print() {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("LRU Cache Content:")
	for e := c.evictionList.Front(); e != nil; e = e.Next() {
		item := e.Value.(*cacheItem)
		fmt.Printf("Key: %s, Value: %+v, Timestamp: %s\n", item.key, item.value, item.timestamp)
	}
}

// RedisCache structure
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// Initialize RedisCache
func NewRedisCache(addr string, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Get an item from RedisCache
func (r *RedisCache) Get(key string) (Cache, bool) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil || err != nil {
		return Cache{}, false
	}
	var cache Cache
	if err := json.Unmarshal([]byte(val), &cache); err != nil {
		return Cache{}, false
	}
	// Update the timestamp when accessed
	cache.CreatedAt = time.Now()
	r.Set(key, cache) // Save updated timestamp
	return cache, true
}

// Set an item in RedisCache
func (r *RedisCache) Set(key string, value Cache) {
	val, err := json.Marshal(value)
	if err != nil {
		return
	}
	r.client.Set(r.ctx, key, val, 5*time.Minute)
}

func (r *RedisCache) Print() {
	fmt.Println("Redis cache content print not supported.")
}

// Memcached structure
type Memcached struct {
	client *memcache.Client
}

// Initialize Memcached
func NewMemcached(addr string) *Memcached {
	client := memcache.New(addr)
	return &Memcached{
		client: client,
	}
}

// Get an item from Memcached
func (m *Memcached) Get(key string) (Cache, bool) {
	item, err := m.client.Get(key)
	if err != nil {
		return Cache{}, false
	}
	var cache Cache
	if err := json.Unmarshal(item.Value, &cache); err != nil {
		return Cache{}, false
	}
	// Update the timestamp when accessed
	cache.CreatedAt = time.Now()
	m.Set(key, cache) // Save updated timestamp
	return cache, true
}

// Set an item in Memcached
func (m *Memcached) Set(key string, value Cache) {
	val, err := json.Marshal(value)
	if err != nil {
		return
	}
	m.client.Set(&memcache.Item{
		Key:        key,
		Value:      val,
		Expiration: int32(5 * time.Minute / time.Second),
	})
}

func (m *Memcached) Print() {
	fmt.Println("Memcached cache content print not supported.")
}

var db *gorm.DB
var cache CacheClient

// Initialize database connection
func initDB() (*gorm.DB, error) {
	dsn := "user=postgres password=Preethi dbname=cachedb sslmode=disable" // Update with your DB credentials
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to the database.")
	db.AutoMigrate(&Cache{}) // Auto migrate to create the table structure with unique constraint
	return db, nil
}

// Authentication middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Username")
		password := r.Header.Get("Password")

		if username == "" || password == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var cacheEntry Cache
		result := db.Where("user = ?", username).First(&cacheEntry)

		if result.Error != nil || cacheEntry.Password != password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Update cache with the latest user details
		cache.Set(cacheEntry.Number, cacheEntry)

		// If authentication is successful, proceed with the request
		next.ServeHTTP(w, r)
	})
}

// Get all users from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUsers called")
	w.Header().Set("Content-type", "application/json")
	var caches []Cache
	if err := db.Find(&caches).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(caches)
}

// Get a specific user by access number
func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUser called")
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	cacheKey := params["id"]

	if cacheValue, found := cache.Get(cacheKey); found {
		json.NewEncoder(w).Encode(cacheValue)
		return
	}

	var cacheEntry Cache
	if err := db.Where("number = ?", cacheKey).First(&cacheEntry).Error; err != nil {
		http.NotFound(w, r)
		return
	}
	cache.Set(cacheKey, cacheEntry) // Update the cache
	json.NewEncoder(w).Encode(cacheEntry)
}

// Create a new user
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateUsers called") // Add this line
	w.Header().Set("Content-type", "application/json")
	var newCache Cache
	if err := json.NewDecoder(r.Body).Decode(&newCache); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if a user with the same username already exists
	var existingUser Cache
	result := db.Where("user = ?", newCache.User).First(&existingUser)

	if result.Error == nil {
		// User with the same username already exists
		http.Error(w, "User with this username already exists", http.StatusConflict)
		return
	}

	// Generate a new access number
	var count int64
	if err := db.Model(&Cache{}).Count(&count).Error; err != nil {
		http.Error(w, "Failed to count users", http.StatusInternalServerError)
		return
	}

	newCache.Number = fmt.Sprintf("%06d", count+1) // Ensure 6-digit access number
	newCache.CreatedAt = time.Now()                // Set the timestamp

	if err := db.Create(&newCache).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	cache.Set(newCache.Number, newCache) // Add to cache
	json.NewEncoder(w).Encode(newCache)
}

// Delete a user by access number
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteUsers called") // Add this line
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	cacheKey := params["id"]
	if err := db.Where("number = ?", cacheKey).Delete(&Cache{}).Error; err != nil {
		http.NotFound(w, r)
		return
	}
	if _, found := cache.Get(cacheKey); found {
		cache.Set(cacheKey, Cache{}) // Changed from Add to Set
	}
	w.WriteHeader(http.StatusNoContent)
}

// PrintCache prints the cache content
func PrintCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	cache.Print() // Call the Print method on the cache

	w.WriteHeader(http.StatusOK)
}

func main() {
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Successfully connected to the database.")

	// Initialize cache (adjust according to your setup)
	cacheType := "Redis" // Change to your configured cache type
	switch cacheType {
	case "LRU":
		cache = NewLRUCache(5, 5*time.Minute)
	case "Redis":
		cache = NewRedisCache("localhost:6379", "", 0)
	case "Memcached":
		cache = NewMemcached("localhost:11211")
	default:
		log.Fatalf("Unsupported cache type: %s", cacheType)
	}

	// Initialize the router and define routes
	router := mux.NewRouter()
	router.HandleFunc("/cache", GetUsers).Methods("GET")            // Get all users
	router.HandleFunc("/cache/{id}", GetUser).Methods("GET")        // Get user by access number
	router.HandleFunc("/cache", CreateUsers).Methods("POST")        // Create a new user
	router.HandleFunc("/cache/{id}", DeleteUsers).Methods("DELETE") // Delete a user by access number
	router.HandleFunc("/cache/print", PrintCache).Methods("GET")    // Print cache content

	// Apply authentication middleware to protected routes
	//authRouter := authMiddleware(router) // Comment out the middleware for testing

	// Start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}
