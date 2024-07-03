package test

import (
	"context"
	measurementmodel "github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
	"github.com/stretchr/testify/mock"
)

type MeasurementRepositoryMock struct {
	mock.Mock
}

func (mr *MeasurementRepositoryMock) GetMeasurementByID(_ context.Context, measurementID uint64) (measurementmodel.Measurement, error) {
	args := mr.Called(measurementID)
	measurement, _ := args.Get(0).(measurementmodel.Measurement)
	e, _ := args.Get(1).(error)
	return measurement, e
}

func (mr *MeasurementRepositoryMock) UpdateMeasurement(_ context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	args := mr.Called(measurement)
	ms, _ := args.Get(0).(measurementmodel.Measurement)
	e, _ := args.Get(1).(error)
	return ms, e
}

func (mr *MeasurementRepositoryMock) DeleteMeasurement(_ context.Context, measurementID uint64) (uint64, error) {
	args := mr.Called(measurementID)
	deletedMeasurementID, _ := args.Get(0).(uint64)
	e, _ := args.Get(1).(error)
	return deletedMeasurementID, e
}

func (mr *MeasurementRepositoryMock) AddMeasurement(_ context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	args := mr.Called(measurement)
	ms, _ := args.Get(0).(measurementmodel.Measurement)
	e, _ := args.Get(1).(error)
	return ms, e
}

func (mr *MeasurementRepositoryMock) GetAllMeasurementsByNodeID(_ context.Context, nodeID uint64) ([]measurementmodel.Measurement, error) {
	args := mr.Called(nodeID)
	measurements, _ := args.Get(0).([]measurementmodel.Measurement)
	e, _ := args.Get(1).(error)
	return measurements, e
}
