package structures

// Data structure to store key-val pairs.
type KeyVal[v any] struct {
	Key   string
	Value v
}

type Store[v any] struct {
	capacity    int
	filledCount int
	loadFactor  int
	bucket      [][]KeyVal[v]
	hashFunc    func(string) uint32
}

func NewStore[v any](capacity int, loadFactor int) *Store[v] {
	return &Store[v]{
		capacity:    capacity,
		bucket:      make([][]KeyVal[v], capacity),
		filledCount: 0,
		loadFactor:  loadFactor,
		hashFunc:    NaiveHash,
	}
}

func (ht *Store[v]) _hash(Key string) int {
	// Double hashing for better distribution at the expense of speed
	h1 := int(ht.hashFunc(Key) % uint32(ht.capacity))
	h2 := int(ht.hashFunc(Key) % uint32(ht.capacity))

	for i := 0; len(ht.bucket[h1]) > 0 && ht.bucket[h1][0].Key != Key; i++ {
		h1 = (h1 + i*h2 + (i*i*i-i)/6) % ht.capacity
		if i == ht.capacity {
			break
		}
	}

	return h1
}

func (ht *Store[v]) Set(Key string, Value v) {
	// While setting, check if the hash table is about to be filled
	// Calculate the load and if the load is about to reach the %load factor increase the bucket size and re-distribute the Key-val pairs
	load := ht.filledCount * 100 / ht.capacity

	if load >= ht.loadFactor {
		// If load reaches given load factor, double the bucket size
		ht.capacity = ht.capacity * 2
		temp := ht.bucket
		ht.bucket = make([][]KeyVal[v], ht.capacity)

		// Rehash every key in the bucket compensating new size
		for _, value := range temp {
			for _, w := range value {
				hash := ht._hash(w.Key)
				ht.bucket[hash] = append(ht.bucket[hash], KeyVal[v]{Key: w.Key, Value: w.Value})
			}
		}
	}

	hash := ht._hash(Key)

	// If Key exists, just update the value
	if len(ht.bucket[hash]) > 0 {
		for i, v := range ht.bucket[hash] {
			if v.Key == Key {
				ht.bucket[hash][i].Value = Value

				return
			}
		}

	}

	// If the slot is empty, append to it
	if len(ht.bucket[hash]) == 0 {
		ht.filledCount = ht.filledCount + 1
	}

	ht.bucket[hash] = append(ht.bucket[hash], KeyVal[v]{Key: Key, Value: Value})
}

func (ht *Store[v]) Get(Key string) *v {
	hash := ht._hash(Key)

	if len(ht.bucket[hash]) > 0 {
		for _, v := range ht.bucket[hash] {
			if v.Key == Key {
				return &v.Value
			}
		}
	}

	return nil
}

func (ht *Store[v]) Remove(Key string) bool {
	hash := ht._hash(Key)

	if len(ht.bucket[hash]) > 0 {
		for i, v := range ht.bucket[hash] {
			if v.Key == Key {
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
