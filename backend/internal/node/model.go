package node

import (
	"time"

	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
)

type Node struct {
	ID           uint64
	Description  string
	Type         string
	LastUpdated  time.Time
	Measurements []measurement.Measurement
}
