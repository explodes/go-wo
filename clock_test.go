package wo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSystemClock_Now(t *testing.T) {
	const epsilon = 3 * float64(time.Microsecond)

	t.Parallel()

	clock := &systemClock{}

	now := time.Now()
	reported := clock.Now()

	assert.InEpsilon(t, now.UnixNano(), reported.UnixNano(), epsilon, "reported now %v not close enough to the real time %v", reported, now)
}

func TestSystemClock_Since(t *testing.T) {
	const (
		nearZero = time.Nanosecond
		epsilon  = 3 * float64(time.Microsecond)
	)

	t.Parallel()

	clock := &systemClock{}

	now := time.Now()
	since := clock.Since(now)

	assert.InEpsilon(t, nearZero, float64(since), epsilon, "reported since %v not close enough to zero ", since)
}

func TestSystemClock_Sleep(t *testing.T) {
	const (
		sleepTime = time.Second
		epsilon   = 3 * float64(time.Microsecond)
	)

	t.Parallel()

	clock := &systemClock{}

	now := time.Now()
	clock.Sleep(sleepTime)
	since := clock.Since(now)

	assert.InEpsilon(t, sleepTime, float64(since), epsilon, "reported sleep time %v not close enough to zero ", since)
}
