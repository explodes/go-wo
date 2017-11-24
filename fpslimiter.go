package wo

import (
	"time"
)

// FpsLimiter is a tool used to limit the number of frames
// drawn or executed to a given fps (frames per second).
type FpsLimiter struct {
	wait       time.Duration
	frameStart time.Time
}

// NewFpsLimiter creates a new FpsLimiter.
func NewFpsLimiter(maxFps float64) *FpsLimiter {
	fpsLimiter := &FpsLimiter{
		frameStart: time.Now(),
	}
	fpsLimiter.SetLimit(maxFps)
	return fpsLimiter
}

// Reset resets the internal start time so that an immediate call
// to StartFrame will return (close to) zero.
func (f *FpsLimiter) Reset() {
	f.frameStart = time.Now()
}

// StartFrame marks the beginning of a frame and returns the time
// since the last frame.
func (f *FpsLimiter) StartFrame() float64 {
	delta := time.Since(f.frameStart).Seconds()
	f.Reset()
	return delta
}

// WaitForNextFrame sleeps the amount of time required to limit
// to the specified fps.
func (f *FpsLimiter) WaitForNextFrame() {
	time.Sleep(f.wait - time.Since(f.frameStart))
}

// SetLimit sets the fps limit.
func (f *FpsLimiter) SetLimit(maxFps float64) {
	f.wait = time.Duration(float64(time.Second) / maxFps)
}

// CurrentFrameFps gets the fps of the current frame.
func (f *FpsLimiter) CurrentFrameFps() float64 {
	return 1 / time.Since(f.frameStart).Seconds()
}
