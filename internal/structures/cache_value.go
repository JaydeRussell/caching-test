package structures

import "time"

type cacheValue struct {
	value      []byte
	expiration time.Time
}
