package structures

import (
	"math/rand/v2"
	"time"
)

// Data structure to store key-val pairs.
type keyVal[v any] struct {
	Key        []byte
	Value      v
	Expiration time.Time
}

type Store[v any] struct {
	capacity    int
	filledCount int
	loadFactor  int
	bucket      [][]keyVal[v]
	hashFunc    func([]byte) uint32
}

func NewStore[v any](capacity int, loadFactor int) *Store[v] {
	return &Store[v]{
		capacity:    capacity,
		bucket:      make([][]keyVal[v], capacity),
		filledCount: 0,
		loadFactor:  loadFactor,
		hashFunc:    NaiveHash,
	}
}

func (ht *Store[v]) _hash(Key []byte) int {
	// Double hashing for better distribution at the expense of speed
	h1 := int(ht.hashFunc(Key) % uint32(ht.capacity))
	h2 := int(ht.hashFunc(Key) % uint32(ht.capacity))

	for i := 0; len(ht.bucket[h1]) > 0 && string(ht.bucket[h1][0].Key) != string(Key); i++ {
		h1 = (h1 + i*h2 + (i*i*i-i)/6) % ht.capacity
		if i == ht.capacity {
			break
		}
	}

	return h1
}

// Set will store the key value pair with a given TTL.
func (ht *Store[v]) Set(Key []byte, Value v, ttl time.Duration) {
	// While setting, check if the hash table is about to be filled
	// Calculate the load and if the load is about to reach the %load factor increase the bucket size and re-distribute the Key-val pairs
	load := ht.filledCount * 100 / ht.capacity

	if load >= ht.loadFactor {
		// If load reaches given load factor, double the bucket size
		ht.capacity = ht.capacity * 2
		temp := ht.bucket
		ht.bucket = make([][]keyVal[v], ht.capacity)

		// Rehash every key in the bucket compensating new size
		for _, value := range temp {
			for _, w := range value {
				hash := ht._hash(w.Key)
				ht.bucket[hash] = append(ht.bucket[hash], keyVal[v]{Key: w.Key, Value: w.Value, Expiration: time.Now().Add(ttl)})
			}
		}
	}

	hash := ht._hash(Key)

	// If Key exists, just update the value
	if len(ht.bucket[hash]) > 0 {
		for i, v := range ht.bucket[hash] {
			if string(v.Key) == string(Key) {
				ht.bucket[hash][i].Value = Value
				return
			}
		}

	}

	// If the slot is empty, append to it
	if len(ht.bucket[hash]) == 0 {
		ht.filledCount = ht.filledCount + 1
	}

	ht.bucket[hash] = append(ht.bucket[hash], keyVal[v]{Key: Key, Value: Value})
}

// Get returns the value stored using `key`.
//
// If the key is not present value will be set to nil.
func (ht *Store[v]) Get(key []byte) (*v, time.Duration) {
	hash := ht._hash(key)

	if len(ht.bucket[hash]) > 0 {
		for _, v := range ht.bucket[hash] {
			if string(v.Key) == string(key) {
				return &v.Value, time.Until(v.Expiration)
			}
		}
	}

	return nil, 0
}

func (ht *Store[v]) Remove(Key []byte) bool {
	hash := ht._hash(Key)

	if len(ht.bucket[hash]) > 0 {
		for i, v := range ht.bucket[hash] {
			if string(v.Key) == string(Key) {
				// Remove the match item from the slot
				ht.bucket[hash] = append(ht.bucket[hash][:i], ht.bucket[hash][:i+1]...)

				// If removing item empties the slot decrease filledCount
				if len(ht.bucket[hash]) == 0 {
					ht.filledCount = ht.filledCount - 1
				}
				return true
			}
		}
	}

	return false
}

func (ht *Store[v]) Length() int {
	return ht.filledCount
}

func (ht *Store[v]) GetRandomValues(num int) (s []keyVal[v]) {
	s = make([]keyVal[v], num)
	h := ht.ToSlice()

	if len(h) <= num {
		return h
	}

	for i := 0; i < num; i++ {
		s = append(s, h[rand.IntN(len(h))])
	}

	return
}

func (ht *Store[v]) ToSlice() (s []keyVal[v]) {
	s = make([]keyVal[v], ht.filledCount)

	for _, v := range ht.bucket {
		for _, v2 := range v {
			if string(v2.Key) != "" {
				s = append(s, v2)
			}
		}
	}

	return
}
