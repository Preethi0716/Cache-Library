package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type CacheItem struct {
	key        string
	value      interface{}
	expiration time.Time
}

type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, found := c.items[key]; found {
		c.list.MoveToFront(element)
		element.Value.(*CacheItem).value = value
		element.Value.(*CacheItem).expiration = time.Now().Add(ttl)
		return nil
	}

	if c.list.Len() >= c.capacity {
		c.evict()
	}

	item := &CacheItem{
		key:        key,
		value:      value,
		expiration: time.Now().Add(ttl),
	}
	element := c.list.PushFront(item)
	c.items[key] = element
	return nil
}

func (c *LRUCache) Get(key string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, found := c.items[key]; found {
		if element.Value.(*CacheItem).expiration.After(time.Now()) {
			c.list.MoveToFront(element)
			return element.Value.(*CacheItem).value, nil
		}
		c.list.Remove(element)
		delete(c.items, key)
		return nil, errors.New("cache miss")
	}
	return nil, errors.New("cache miss")
}

func (c *LRUCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, found := c.items[key]; found {
		c.list.Remove(element)
		delete(c.items, key)
		return nil
	}
	return errors.New("cache miss")
}

func (c *LRUCache) GetAll() (map[string]interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	allItems := make(map[string]interface{})
	for key, element := range c.items {
		if element.Value.(*CacheItem).expiration.After(time.Now()) {
			allItems[key] = element.Value.(*CacheItem).value
		}
	}
	return allItems, nil
}

func (c *LRUCache) evict() {
	if element := c.list.Back(); element != nil {
		c.list.Remove(element)
		delete(c.items, element.Value.(*CacheItem).key)
	}
}
