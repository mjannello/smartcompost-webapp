package app

import (
	"context"
	"fmt"
	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
)

type MeasurementService interface {
	GetMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurement.Measurement, error)
	AddMeasurements(ctx context.Context, measurements []measurement.Measurement) error
}

type measurementService struct {
	repo measurement.Repository
}

func NewService(repo measurement.Repository) MeasurementService {
	return &measurementService{repo: repo}
}

func (ms *measurementService) GetMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurement.Measurement, error) {
	measurements, err := ms.repo.GetAllMeasurementsByNodeID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("could not get measurements: %w", err)
	}
	return measurements, nil
}

func (ms *measurementService) AddMeasurements(ctx context.Context, measurements []measurement.Measurement) error {
	for _, m := range measurements {
		if err := ms.repo.AddMeasurement(ctx, m); err != nil {
			return fmt.Errorf("error adding measurement: %w", err)
		}
	}
	return nil
}
