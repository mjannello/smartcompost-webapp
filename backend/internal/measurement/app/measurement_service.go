package app

import (
	"context"
	"fmt"
	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
)

type MeasurementService interface {
	GetMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurement.Measurement, error)
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
