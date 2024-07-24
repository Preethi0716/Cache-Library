package cache

import "errors"

var ErrCacheMiss = errors.New("cache miss")

var ErrKeyNotFound = errors.New("key not found")
