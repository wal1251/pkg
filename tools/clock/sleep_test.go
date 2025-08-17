package clock_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/wal1251/pkg/tools/clock"
)

func TestSleep(t *testing.T) {
	sleepDuration := 2 * time.Second
	fuzzFactor := 1.5
	expectedFuzz := math.Round(float64(sleepDuration) * fuzzFactor)
	expectedSleepDuration := sleepDuration + time.Duration(rand.Int63n(int64(expectedFuzz)))

	start := time.Now()
	actualSleepDuration := clock.SleepFuzz(sleepDuration)
	elapsed := time.Since(start)

	if elapsed > expectedSleepDuration {
		t.Errorf("Sleep time is less than expected. Expected: %v, actual: %v", expectedSleepDuration, elapsed)
	}

	if actualSleepDuration > expectedSleepDuration {
		t.Errorf("Returned sleep duration is less than expected. Expected: %v, actual: %v", expectedSleepDuration, actualSleepDuration)
	}
}
