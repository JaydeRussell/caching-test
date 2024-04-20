package structures

import "time"

type Bucket struct {
	// TODO actually create this definition to meet criteria
	bucket map[string]cacheValue
}

func NewBucket() *Bucket {
	return &Bucket{
		bucket: map[string]cacheValue{},
	}
}

func (b *Bucket) Set(key, value []byte, ttl time.Duration) {
	if key == nil {
		panic("key cannot be nil")
	}

	b.bucket[string(key)] = cacheValue{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (b *Bucket) Get(key []byte) []byte {
	return b.get(string(key))
}

func (b Bucket) TTL(key []byte) time.Duration {
	if v, ok := b.bucket[string(key)]; ok {
		return time.Until(v.expiration)
	}

	return -1
}

func (b *Bucket) get(key string) []byte {
	return b.bucket[key].value
}

func (b *Bucket) Delete(key string) (deleted []byte) {
	deleted = b.get(key)

	delete(b.bucket, key)
	return
}
