package wo

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

const (
	fpsTestFps       = 60
	fpsTestFrameRate = 1. / 60.

	fpsTestTimeEpsilon = 0.0000001
)

type fpsTest func(*testing.T, *FpsLimiter, *fakeClock)

func TestFpsLimiter(t *testing.T) {
	tests := []fpsTest{
		testFpsLimiter_CurrentFrameFps,
		testFpsLimiter_Reset,
		testFpsLimiter_SetLimit,
		testFpsLimiter_StartFrame,
		testFpsLimiter_WaitForNextFrame,
	}
	for _, test := range tests {
		test := test
		t.Run(funcName(t, test), func(t *testing.T) {
			t.Parallel()
			clock := NewFakeClock()
			fps := newFpsLimiterClock(60, clock)
			test(t, fps, clock)
		})
	}
}

func TestNewFpsLimiter(t *testing.T) {
	fps := NewFpsLimiter(60)

	assert.NotNil(t, fps)
}

func testFpsLimiter_CurrentFrameFps(t *testing.T, fps *FpsLimiter, clock *fakeClock) {
	clock.AdvanceSeconds(fpsTestFrameRate)

	framesPerSecond := fps.CurrentFrameFps()

	assert.InEpsilon(t, fpsTestFps, framesPerSecond, fpsTestTimeEpsilon)
}

func testFpsLimiter_Reset(t *testing.T, fps *FpsLimiter, clock *fakeClock) {
	clock.AdvanceSeconds(2)
	fps.Reset()

	fpsTime := fps.frameStart.Sub(time.Time{}).Seconds()

	assert.InEpsilon(t, 2, fpsTime, fpsTestTimeEpsilon)
}

func testFpsLimiter_SetLimit(t *testing.T, fps *FpsLimiter, clock *fakeClock) {
	fps.SetLimit(30)
	fps.WaitForNextFrame()

	framesPerSecond := fps.CurrentFrameFps()

	assert.InEpsilon(t, 30, framesPerSecond, fpsTestTimeEpsilon)
}

func testFpsLimiter_StartFrame(t *testing.T, fps *FpsLimiter, clock *fakeClock) {
	clock.AdvanceSeconds(1)

	diff := fps.StartFrame()

	assert.InEpsilon(t, 1, diff, fpsTestTimeEpsilon)
}

func testFpsLimiter_WaitForNextFrame(t *testing.T, fps *FpsLimiter, clock *fakeClock) {
	fps.WaitForNextFrame()

	now := clock.ElapsedDuration().Seconds()

	assert.InEpsilon(t, fpsTestFrameRate, now, fpsTestTimeEpsilon)
}
