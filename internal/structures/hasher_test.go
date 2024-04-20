package structures_test

import (
	"caching-test/internal/structures"
	"testing"
)

func BenchmarkNaiveHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		structures.NaiveHash("apple")
	}
}
