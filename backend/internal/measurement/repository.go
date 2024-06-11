package measurement

import "context"

type Repository interface {
	GetAllMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]Measurement, error)
}
