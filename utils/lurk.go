package utils

import (
	"math/rand"
	"time"
)

const durationEntropy = 5 // In percent

// LurkDuration increases the duration by a random amount
// ranging from 0 to {DURATION_ENTROPOY}% of the duration
func LurkDuration(duration time.Duration) time.Duration {

	return duration + time.Duration(rand.Int63n(duration.Nanoseconds()*durationEntropy/100))
}
