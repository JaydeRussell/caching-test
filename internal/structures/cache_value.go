package structures

import "time"

type CacheValue struct {
	Value      []byte
	Expiration time.Time
}
