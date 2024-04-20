package structures_test

import (
	"caching-test/internal/structures"
	"testing"
	"time"
)

func BenchmarkStoreSet(b *testing.B) {
	hashTable := structures.NewStore[structures.CacheValue](1000, 50)

	for i := 0; i < b.N; i++ {
		hashTable.Set("hello", structures.CacheValue{Value: []byte("world"), Expiration: time.Now().Add(time.Hour * 1)})
	}
}

func BenchmarkStoreGet(b *testing.B) {
	hashTable := structures.NewStore[structures.CacheValue](1000, 50)
	hashTable.Set("hello", structures.CacheValue{Value: []byte("world"), Expiration: time.Now().Add(time.Hour * 1)})

	for i := 0; i < b.N; i++ {
		hashTable.Get("hello")
	}
}
