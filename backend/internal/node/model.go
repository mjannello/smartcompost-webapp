package node

import (
	"time"

	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
)

type Node struct {
	ID           uint64
	SerialNumber string
	Description  string
	Model        string
	DateCreated  time.Time
	LastUpdated  time.Time
	Measurements []measurement.Measurement
}
