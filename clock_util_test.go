package wo

import (
	"sync"
	"time"
)

type fakeClock struct {
	*sync.Mutex
	now time.Time
}

func NewFakeClock() *fakeClock {
	return &fakeClock{
		Mutex: &sync.Mutex{},
		now:   time.Time{},
	}
}

func (c *fakeClock) Sleep(d time.Duration) {
	if d <= 0 {
		return
	}
	c.Advance(d)
}

func (c *fakeClock) Now() time.Time {
	return c.now
}

func (c *fakeClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

func (c *fakeClock) Advance(d time.Duration) {
	c.Lock()
	c.now = c.now.Add(d)
	c.Unlock()
}

func (c *fakeClock) AdvanceSeconds(dt float64) {
	c.Lock()
	d := time.Duration(int64(dt * float64(time.Second)))
	c.now = c.now.Add(d)
	c.Unlock()
}

func (c *fakeClock) Set(when time.Time) {
	c.now = when
}

func (c *fakeClock) ElapsedDuration() time.Duration {
	return time.Duration(c.now.UnixNano())
}

func (c *fakeClock) ElapsedSeconds() float64 {
	return c.now.Sub(time.Time{}).Seconds()
}
