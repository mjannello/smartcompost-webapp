package test

import (
	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
	"time"
)

func MakeMeasurement(id, nodeID uint64, value float64, mtype string, timestamp time.Time) measurement.Measurement {
	return measurement.Measurement{
		ID:        id,
		NodeID:    nodeID,
		Value:     value,
		Type:      mtype,
		Timestamp: timestamp,
	}
}
