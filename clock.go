package wo

import "time"

type clock interface {
	Now() time.Time
	Sleep(time.Duration)
	Since(t time.Time) time.Duration
}

type systemClock struct {
}

func (s *systemClock) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

func (s *systemClock) Now() time.Time {
	return time.Now()
}

func (s *systemClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}
