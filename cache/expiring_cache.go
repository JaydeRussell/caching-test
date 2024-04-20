package cache

import (
	"time"

	"caching-test/internal/structures"
)

type ExpiringCache struct {
	bucket *structures.Store[[]byte]
}

func NewExpiringCache() *ExpiringCache {
	return &ExpiringCache{
		bucket: structures.NewStore[[]byte](1000, 50),
	}
}

// Set will store the key value pair with a given TTL.
func (c *ExpiringCache) Set(key, value []byte, ttl time.Duration) {
	c.bucket.Set(string(key), value)
}

// Get returns the value stored using `key`.
//
// If the key is not present value will be set to nil.
func (c *ExpiringCache) Get(key []byte) ([]byte, time.Duration) {
	v := *c.bucket.Get(string(key))
	if v == nil {
		return nil, -1
	}

	return v, 0
}
