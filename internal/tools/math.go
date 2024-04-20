package tools

import (
	"caching-test/interfaces"
	"math"
)

func ApproximatelyEquals[t interfaces.Number](a, b, drift t) bool {
	return math.Abs(float64(a-b)) <= float64(drift)
}
