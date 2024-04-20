package interfaces

import "time"

type Number interface {
	time.Duration | int64 | int | float64 | float32
}
