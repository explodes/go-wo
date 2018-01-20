package wo

import (
	"time"
)

// FpsLimiter is a tool used to limit the number of frames
// drawn or executed to a given fps (frames per second).
type FpsLimiter struct {
	clock clock

	// wait is duration to wait each frame to keep a
	// consistent FPS
	wait time.Duration

	// frameStart marks the time a frame was started
	// with StartFrame()
	frameStart time.Time
}

// NewFpsLimiter creates a new FpsLimiter.
func NewFpsLimiter(maxFps float64) *FpsLimiter {
	return newFpsLimiterClock(maxFps, &systemClock{})
}

// NewFpsLimiter creates a new FpsLimiter with the given clock
func newFpsLimiterClock(maxFps float64, clock clock) *FpsLimiter {
	fpsLimiter := &FpsLimiter{
		clock: clock,
	}
	fpsLimiter.Reset()
	fpsLimiter.SetLimit(maxFps)
	return fpsLimiter
}

// Reset resets the internal start time so that an immediate call
// to StartFrame will return (close to) zero.
func (f *FpsLimiter) Reset() {
	f.frameStart = f.clock.Now()
}

// StartFrame marks the beginning of a frame and returns the time
// since the last frame.
func (f *FpsLimiter) StartFrame() float64 {
	delta := f.clock.Since(f.frameStart).Seconds()
	f.Reset()
	return delta
}

// WaitForNextFrame sleeps the amount of time required to limit
// to the specified fps.
func (f *FpsLimiter) WaitForNextFrame() {
	f.clock.Sleep(f.wait - f.clock.Since(f.frameStart))
}

// SetLimit sets the fps limit.
func (f *FpsLimiter) SetLimit(maxFps float64) {
	f.wait = time.Duration(float64(time.Second) / maxFps)
}

// CurrentFrameFps gets the fps of the current frame.
func (f *FpsLimiter) CurrentFrameFps() float64 {
	return 1 / f.clock.Since(f.frameStart).Seconds()
}
