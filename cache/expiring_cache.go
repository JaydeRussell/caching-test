package cache

import (
	"time"

	"caching-test/internal/structures"
)

type ExpiringCache struct {
	bucket *structures.Bucket
}

func NewExpiringCache() *ExpiringCache {
	return &ExpiringCache{
		bucket: structures.NewBucket(),
	}
}

// Set will store the key value pair with a given TTL.
func (c *ExpiringCache) Set(key, value []byte, ttl time.Duration) {
	c.bucket.Set(key, value, ttl)
}

// Get returns the value stored using `key`.
//
// If the key is not present value will be set to nil.
func (c *ExpiringCache) Get(key []byte) (value []byte, ttl time.Duration) {
	ttl = c.bucket.TTL(key)
	if ttl <= 0 {
		return nil, -1
	}

	value = c.bucket.Get(key)
	return
}
