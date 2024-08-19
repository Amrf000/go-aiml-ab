package ab

import (
	"time"
)

type Timer struct {
	startTime time.Time
}

func NewTimer() *Timer {
	return &Timer{
		startTime: time.Now(),
	}
}

func (t *Timer) Start() {
	t.startTime = time.Now()
}

func (t *Timer) ElapsedTimeMillis() int64 {
	return time.Since(t.startTime).Milliseconds() + 1
}

func (t *Timer) ElapsedRestartMs() int64 {
	ms := time.Since(t.startTime).Milliseconds() + 1
	t.Start()
	return ms
}

func (t *Timer) ElapsedTimeSecs() float64 {
	return float64(t.ElapsedTimeMillis()) / 1000.0
}

func (t *Timer) ElapsedTimeMins() float64 {
	return t.ElapsedTimeSecs() / 60.0
}
