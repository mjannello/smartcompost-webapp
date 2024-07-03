package test

import (
	"context"
	measurementmodel "github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
	nodemodel "github.com/mjannello/smartcompost-webapp/backend/internal/node"
	"github.com/stretchr/testify/mock"
	"time"
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

func (mr *MeasurementRepositoryMock) AddMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	//TODO implement me
	panic("implement me")
}

func (mr *MeasurementRepositoryMock) GetAllMeasurementsByNodeID(_ context.Context, nodeID uint64) ([]measurementmodel.Measurement, error) {
	args := mr.Called(nodeID)
	measurements, _ := args.Get(0).([]measurementmodel.Measurement)
	e, _ := args.Get(1).(error)
	return measurements, e
}

type NodeServiceMock struct {
	mock.Mock
}

func (n NodeServiceMock) GetNodes(ctx context.Context) ([]nodemodel.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (n NodeServiceMock) GetNodeByID(ctx context.Context, nodeID uint64) (nodemodel.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (n NodeServiceMock) UpdateNode(ctx context.Context, node nodemodel.Node) (nodemodel.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (n NodeServiceMock) UpdateNodeLastUpdated(ctx context.Context, nodeID uint64, lastUpdated time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (n NodeServiceMock) DeleteNode(ctx context.Context, nodeID uint64) (uint64, error) {
	//TODO implement me
	panic("implement me")
}
