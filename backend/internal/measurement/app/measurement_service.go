package app

import (
	"context"
	"fmt"
	nodeapp "github.com/mjannello/smartcompost-webapp/backend/internal/node/app"
	"log"
	"time"

	measurementmodel "github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
)

type MeasurementService interface {
	GetMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurementmodel.Measurement, error)
	GetMeasurementByID(ctx context.Context, measurementID uint64) (measurementmodel.Measurement, error)
	UpdateMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error)
	DeleteMeasurement(ctx context.Context, measurementID uint64) (uint64, error)
	AddNodeMeasurements(ctx context.Context, nodeID uint64, measurement []measurementmodel.Measurement) ([]measurementmodel.Measurement, error)
}

type measurementService struct {
	measurementRepository measurementmodel.Repository
	nodeService           nodeapp.NodeService
}

func NewMeasurementService(mr measurementmodel.Repository, ns nodeapp.NodeService) MeasurementService {
	return &measurementService{
		measurementRepository: mr,
		nodeService:           ns}
}

func (ms *measurementService) GetMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurementmodel.Measurement, error) {
	measurements, err := ms.measurementRepository.GetAllMeasurementsByNodeID(ctx, nodeID)
	if err != nil {
		log.Printf("Error fetching measurements for node ID %d: %v", nodeID, err)
		return nil, fmt.Errorf("error getting measurements: %w", err)
	}
	log.Printf("Fetched %d measurements for node ID %d", len(measurements), nodeID)
	return measurements, nil
}

func (ms *measurementService) GetMeasurementByID(ctx context.Context, measurementID uint64) (measurementmodel.Measurement, error) {
	measurement, err := ms.measurementRepository.GetMeasurementByID(ctx, measurementID)
	if err != nil {
		log.Printf("Error fetching measurement by ID %d: %v", measurementID, err)
		return measurementmodel.Measurement{}, fmt.Errorf("error getting measurement by ID: %w", err)
	}
	log.Printf("Fetched measurement by ID %d: %+v", measurementID, measurement)
	return measurement, nil
}

func (ms *measurementService) UpdateMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	updatedMeasurement, err := ms.measurementRepository.UpdateMeasurement(ctx, measurement)
	if err != nil {
		log.Printf("Error updating measurement with ID %d: %v", measurement.ID, err)
		return measurementmodel.Measurement{}, fmt.Errorf("error updating measurement: %w", err)
	}
	log.Printf("Updated measurement: %+v", updatedMeasurement)
	return updatedMeasurement, nil
}

func (ms *measurementService) DeleteMeasurement(ctx context.Context, measurementID uint64) (uint64, error) {
	deletedID, err := ms.measurementRepository.DeleteMeasurement(ctx, measurementID)
	if err != nil {
		log.Printf("Error deleting measurement with ID %d: %v", measurementID, err)
		return 0, fmt.Errorf("error deleting measurement: %w", err)
	}
	log.Printf("Deleted measurement with ID %d", deletedID)
	return deletedID, nil
}

func (ms *measurementService) AddNodeMeasurements(ctx context.Context, nodeID uint64, measurements []measurementmodel.Measurement) ([]measurementmodel.Measurement, error) {
	// Validate that the node exists
	_, err := ms.nodeService.GetNodeByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	// Add each measurement
	var createdMeasurements []measurementmodel.Measurement
	for _, measurement := range measurements {
		measurement.NodeID = nodeID
		createdMeasurement, err := ms.measurementRepository.AddMeasurement(ctx, measurement)
		if err != nil {
			log.Printf("Error adding measurement: %v", err)
			return nil, fmt.Errorf("error adding measurement: %w", err)
		}
		createdMeasurements = append(createdMeasurements, createdMeasurement)
	}

	// Update the node's last_updated field
	lastUpdated := time.Now()
	err = ms.nodeService.UpdateNodeLastUpdated(ctx, nodeID, lastUpdated)
	if err != nil {
		return nil, fmt.Errorf("error updating node last_updated: %w", err)
	}

	log.Printf("Added measurements: %+v", createdMeasurements)
	return createdMeasurements, nil
}
