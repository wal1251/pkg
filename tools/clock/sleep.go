package clock

import (
	"math"
	"math/rand"
	"time"
)

const sleepFuzzFactor = 0.1

// SleepFuzz аналогично time.Sleep, но время сна вычисляется с фазингом около 10% от целевого значения.
func SleepFuzz(sleepDuration time.Duration) time.Duration {
	val := math.Round(float64(sleepDuration) * sleepFuzzFactor)
	fuzz := time.Duration(rand.Int63n(int64(val))) //nolint:gosec
	sleepDuration += fuzz

	time.Sleep(sleepDuration)

	return sleepDuration
}
