package clock

import (
	"time"
)

type Clock interface {
	Time() time.Time
}

type clock struct{}

func NewClock() Clock { return &clock{} }

func (c clock) Time() time.Time { return time.Now() }
