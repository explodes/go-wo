package wo

import (
	"testing"

	"time"

	"fmt"

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

func BenchmarkFpsLimiter_Loop(b *testing.B) {
	loopSizes := []int{
		10, 100, 1000, 10000,
	}
	for _, loopSize := range loopSizes {
		loopSize := loopSize
		b.Run(fmt.Sprintf("loopSize[%d]", loopSize), func(b *testing.B) {
			fps := newFpsLimiterClock(60, NewFakeClock())
			for i := 0; i < b.N; i++ {
				for j := 0; j < loopSize; j++ {
					fps.StartFrame()
					fps.WaitForNextFrame()
				}
			}
		})
	}
}
