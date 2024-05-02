package structures_test

import (
	"caching-test/internal/structures"
	"testing"
	"time"
)

func BenchmarkStoreSet(b *testing.B) {
	hashTable := structures.NewStore[structures.CacheValue](1000, 50)
	key := []byte("hello")

	for i := 0; i < b.N; i++ {
		hashTable.Set(key, structures.CacheValue{Value: []byte("world"), Expiration: time.Now().Add(time.Hour * 1)}, time.Hour)
	}
}

func BenchmarkStoreGet(b *testing.B) {
	hashTable := structures.NewStore[structures.CacheValue](1000, 50)
	key := []byte("hello")
	hashTable.Set(key, structures.CacheValue{Value: []byte("world"), Expiration: time.Now().Add(time.Hour * 1)}, time.Hour)
	for i := 0; i < b.N; i++ {
		hashTable.Get(key)
	}
}
