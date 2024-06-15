package measurement

import (
	"context"
)

type Repository interface {
	GetAllMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]Measurement, error)
	GetMeasurementByID(ctx context.Context, measurementID uint64) (Measurement, error)
	UpdateMeasurement(ctx context.Context, measurement Measurement) (Measurement, error)
	DeleteMeasurement(ctx context.Context, measurementID uint64) (uint64, error)
	AddMeasurement(ctx context.Context, measurement Measurement) (Measurement, error)
}
