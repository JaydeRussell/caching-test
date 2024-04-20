package cache_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"caching-test/cache"
	"caching-test/internal/tools"
)

func getBaseTestCache() (testCache *cache.ExpiringCache) {
	testCache = cache.NewExpiringCache()

	testCache.Set([]byte("hello"), []byte("world"), time.Hour*1)
	return
}

func TestSet(t *testing.T) {
	testCases := []struct {
		description string
		c           *cache.ExpiringCache
		key         []byte
		value       []byte
		ttl         time.Duration

		expectedValue      []byte
		expectedTTL        time.Duration
		allowableTTLMargin time.Duration
	}{
		{
			description: "basic set call",
			c:           cache.NewExpiringCache(),
			key:         []byte("hello"),
			value:       []byte("world"),
			ttl:         time.Hour * 1,

			expectedValue:      []byte("world"),
			expectedTTL:        time.Hour * 1,
			allowableTTLMargin: time.Millisecond * 20,
		},
	}

	for _, tc := range testCases {
		c := tc.c
		key := tc.key

		c.Set(tc.key, tc.value, tc.ttl)

		actualValue, actualTTL := c.Get(key)

		if string(actualValue) != string(tc.expectedValue) {
			log.Printf("%s failed because actual value did not equal expected. Actual: %#v Expected: %#v",
				tc.description, actualValue, tc.expectedValue)
			t.FailNow()
		}

		if !tools.ApproximatelyEquals(actualTTL, tc.expectedTTL, tc.allowableTTLMargin) {
			log.Printf("%s failed because actual TTL did not equal expected. Actual: %#v Expected: %#v acceptable Margin: %#v",
				tc.description, actualValue, tc.expectedValue, tc.allowableTTLMargin)
			t.FailNow()
		}
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		description string
		c           *cache.ExpiringCache
		key         []byte

		expectedValue      []byte
		expectedTTL        time.Duration
		allowableTTLMargin time.Duration
	}{
		{
			description: "basic get call",
			c:           getBaseTestCache(),
			key:         []byte("hello"),

			expectedValue:      []byte("world"),
			expectedTTL:        time.Hour * 1,
			allowableTTLMargin: time.Millisecond * 20,
		},
	}

	for _, tc := range testCases {
		c := tc.c
		key := tc.key

		actualValue, actualTTL := c.Get(key)

		if string(actualValue) != string(tc.expectedValue) {
			log.Printf("%s failed because actual value did not equal expected. Actual: %#v Expected: %#v",
				tc.description, actualValue, tc.expectedValue)
			t.FailNow()
		}

		if !tools.ApproximatelyEquals(actualTTL, tc.expectedTTL, tc.allowableTTLMargin) {
			log.Printf("%s failed because actual TTL did not equal expected. Actual: %#v Expected: %#v acceptable Margin: %#v",
				tc.description, actualValue, tc.expectedValue, tc.allowableTTLMargin)
			t.FailNow()
		}
	}
}

func BenchmarkSet(b *testing.B) {
	testCache := getBaseTestCache()

	for i := 0; i < b.N; i++ {
		testCache.Set(
			[]byte(fmt.Sprintf("Key %d", i)),
			[]byte(fmt.Sprintf("Value %d", i)),
			time.Hour*1,
		)
	}
}

func BenchmarkGet(b *testing.B) {
	testCache := getBaseTestCache()
	key := []byte("hello")

	for i := 0; i < b.N; i++ {
		testCache.Get(key)
	}
}
