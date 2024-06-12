package measurement

import "time"

type Measurement struct {
	ID        uint64
	NodeID    uint64
	Type      string
	Value     float64
	Timestamp time.Time
}
