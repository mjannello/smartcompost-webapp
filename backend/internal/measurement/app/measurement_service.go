package app

import (
	"context"
	"fmt"
	"log"

	measurementmodel "github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
)

type MeasurementService interface {
	GetMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurementmodel.Measurement, error)
	GetMeasurementByID(ctx context.Context, measurementID uint64) (measurementmodel.Measurement, error)
	UpdateMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error)
	DeleteMeasurement(ctx context.Context, measurementID uint64) (uint64, error)
	AddMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error)
}

type measurementService struct {
	repo measurementmodel.Repository
}

func NewMeasurementService(repo measurementmodel.Repository) MeasurementService {
	return &measurementService{repo: repo}
}

func (ms *measurementService) GetMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurementmodel.Measurement, error) {
	measurements, err := ms.repo.GetAllMeasurementsByNodeID(ctx, nodeID)
	if err != nil {
		log.Printf("Error fetching measurements for node ID %d: %v", nodeID, err)
		return nil, fmt.Errorf("error getting measurements: %w", err)
	}
	log.Printf("Fetched %d measurements for node ID %d", len(measurements), nodeID)
	return measurements, nil
}

func (ms *measurementService) GetMeasurementByID(ctx context.Context, measurementID uint64) (measurementmodel.Measurement, error) {
	measurement, err := ms.repo.GetMeasurementByID(ctx, measurementID)
	if err != nil {
		log.Printf("Error fetching measurement by ID %d: %v", measurementID, err)
		return measurementmodel.Measurement{}, fmt.Errorf("error getting measurement by ID: %w", err)
	}
	log.Printf("Fetched measurement by ID %d: %+v", measurementID, measurement)
	return measurement, nil
}

func (ms *measurementService) UpdateMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	updatedMeasurement, err := ms.repo.UpdateMeasurement(ctx, measurement)
	if err != nil {
		log.Printf("Error updating measurement with ID %d: %v", measurement.ID, err)
		return measurementmodel.Measurement{}, fmt.Errorf("error updating measurement: %w", err)
	}
	log.Printf("Updated measurement: %+v", updatedMeasurement)
	return updatedMeasurement, nil
}

func (ms *measurementService) DeleteMeasurement(ctx context.Context, measurementID uint64) (uint64, error) {
	deletedID, err := ms.repo.DeleteMeasurement(ctx, measurementID)
	if err != nil {
		log.Printf("Error deleting measurement with ID %d: %v", measurementID, err)
		return 0, fmt.Errorf("error deleting measurement: %w", err)
	}
	log.Printf("Deleted measurement with ID %d", deletedID)
	return deletedID, nil
}

func (ms *measurementService) AddMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	createdMeasurement, err := ms.repo.AddMeasurement(ctx, measurement)
	if err != nil {
		log.Printf("Error adding measurement: %v", err)
		return measurementmodel.Measurement{}, fmt.Errorf("error adding measurement: %w", err)
	}
	log.Printf("Added measurement: %+v", createdMeasurement)
	return createdMeasurement, nil
}
