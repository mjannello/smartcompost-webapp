package app_test

import (
	"context"
	"fmt"
	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement/app"
	"github.com/mjannello/smartcompost-webapp/backend/pkg/clock"
	"github.com/mjannello/smartcompost-webapp/backend/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMeasurementService_GetMeasurementsByNodeID(t *testing.T) {

	type depFields struct {
		measurementRepositoryMock *test.MeasurementRepositoryMock
		nodeServiceMock           *test.NodeServiceMock
	}
	type input struct {
		nodeID uint64
	}
	type output struct {
		measurements []measurement.Measurement
		err          error
	}

	expectedMeasurements := []measurement.Measurement{
		{
			NodeID: uint64(1),
		},
	}

	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "get measurements by NodeID successfully",
			in:   input{nodeID: uint64(1)},
			on: func(df *depFields) {
				df.measurementRepositoryMock.On("GetAllMeasurementsByNodeID", uint64(1)).Return(expectedMeasurements, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, expectedMeasurements, out.measurements)
			},
		},
		{
			name: "error getting measurements by NodeID",
			in:   input{nodeID: uint64(1)},
			on: func(df *depFields) {
				df.measurementRepositoryMock.On("GetAllMeasurementsByNodeID", uint64(1)).Return(nil, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.Error(t, out.err)
				assert.ErrorContains(t, out.err, "error getting measurements: test")
				assert.Nil(t, out.measurements)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			measurementRepositoryMock := &test.MeasurementRepositoryMock{}
			nodeServiceMock := &test.NodeServiceMock{}
			clockMock := &clock.ClockMock{}
			s := app.NewMeasurementService(measurementRepositoryMock, nodeServiceMock, clockMock)

			df := &depFields{measurementRepositoryMock: measurementRepositoryMock, nodeServiceMock: nodeServiceMock}
			tt.on(df)

			// When
			resultMeasurements, err := s.GetMeasurementsByNodeID(context.Background(), tt.in.nodeID)

			// Then
			tt.assert(t, &output{resultMeasurements, err})
			measurementRepositoryMock.AssertExpectations(t)
			nodeServiceMock.AssertExpectations(t)

		})
	}
}

func TestMeasurementService_GetMeasurementByID(t *testing.T) {
	type depFields struct {
		measurementRepositoryMock *test.MeasurementRepositoryMock
		nodeServiceMock           *test.NodeServiceMock
	}
	type input struct {
		nodeID uint64
	}
	type output struct {
		measurement measurement.Measurement
		err         error
	}

	expectedMeasurement := measurement.Measurement{
		NodeID: uint64(1),
	}

	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "get measurement by ID successfully",
			in:   input{nodeID: uint64(1)},
			on: func(df *depFields) {
				df.measurementRepositoryMock.On("GetMeasurementByID", uint64(1)).Return(expectedMeasurement, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, expectedMeasurement, out.measurement)
			},
		},
		{
			name: "error getting measurement by ID",
			in:   input{nodeID: uint64(1)},
			on: func(df *depFields) {
				df.measurementRepositoryMock.On("GetMeasurementByID", uint64(1)).Return(measurement.Measurement{}, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.Error(t, out.err)
				assert.ErrorContains(t, out.err, "error getting measurement by ID: test")
				assert.Equal(t, measurement.Measurement{}, out.measurement)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			measurementRepositoryMock := &test.MeasurementRepositoryMock{}
			nodeServiceMock := &test.NodeServiceMock{}
			clockMock := &clock.ClockMock{}
			s := app.NewMeasurementService(measurementRepositoryMock, nodeServiceMock, clockMock)

			df := &depFields{measurementRepositoryMock: measurementRepositoryMock, nodeServiceMock: nodeServiceMock}
			tt.on(df)

			// When
			resultMeasurement, err := s.GetMeasurementByID(context.Background(), tt.in.nodeID)

			// Then
			tt.assert(t, &output{resultMeasurement, err})
			measurementRepositoryMock.AssertExpectations(t)
			nodeServiceMock.AssertExpectations(t)

		})
	}
}

func TestMeasurementService_UpdateMeasurement(t *testing.T) {
	type depFields struct {
		measurementRepositoryMock *test.MeasurementRepositoryMock
		nodeServiceMock           *test.NodeServiceMock
	}
	type input struct {
		ms measurement.Measurement
	}
	type output struct {
		measurement measurement.Measurement
		err         error
	}

	inputMeasurement := measurement.Measurement{
		NodeID: uint64(1),
		Value:  20.0,
	}
	expectedMeasurement := measurement.Measurement{
		NodeID: uint64(1),
		Value:  25.0,
	}

	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "update measurement successfully",
			in:   input{ms: inputMeasurement},
			on: func(df *depFields) {
				df.measurementRepositoryMock.On("UpdateMeasurement", inputMeasurement).Return(expectedMeasurement, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, expectedMeasurement, out.measurement)
			},
		},
		{
			name: "error updating measurement by ID",
			in:   input{ms: inputMeasurement},
			on: func(df *depFields) {
				df.measurementRepositoryMock.On("UpdateMeasurement", inputMeasurement).Return(measurement.Measurement{}, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.Error(t, out.err)
				assert.ErrorContains(t, out.err, "error updating measurement: test")
				assert.Equal(t, measurement.Measurement{}, out.measurement)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			measurementRepositoryMock := &test.MeasurementRepositoryMock{}
			nodeServiceMock := &test.NodeServiceMock{}
			clockMock := &clock.ClockMock{}
			s := app.NewMeasurementService(measurementRepositoryMock, nodeServiceMock, clockMock)

			df := &depFields{measurementRepositoryMock: measurementRepositoryMock, nodeServiceMock: nodeServiceMock}
			tt.on(df)

			// When
			resultMeasurement, err := s.UpdateMeasurement(context.Background(), tt.in.ms)

			// Then
			tt.assert(t, &output{resultMeasurement, err})
			measurementRepositoryMock.AssertExpectations(t)
			nodeServiceMock.AssertExpectations(t)

		})
	}
}

func TestMeasurementService_DeleteMeasurement(t *testing.T) {
	type depFields struct {
		measurementRepositoryMock *test.MeasurementRepositoryMock
		nodeServiceMock           *test.NodeServiceMock
	}
	type input struct {
		measurementID uint64
	}
	type output struct {
		deletedMeasurementID uint64
		err                  error
	}

	expectedMeasurementID := uint64(1)

	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "delete measurement by ID successfully",
			in:   input{measurementID: expectedMeasurementID},
			on: func(df *depFields) {
				df.measurementRepositoryMock.On("DeleteMeasurement", expectedMeasurementID).Return(expectedMeasurementID, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, expectedMeasurementID, out.deletedMeasurementID)
			},
		},
		{
			name: "error deleting measurement by ID",
			in:   input{measurementID: expectedMeasurementID},
			on: func(df *depFields) {
				df.measurementRepositoryMock.On("DeleteMeasurement", expectedMeasurementID).Return(uint64(0), fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.Error(t, out.err)
				assert.ErrorContains(t, out.err, "error deleting measurement: test")
				assert.Equal(t, uint64(0), out.deletedMeasurementID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			measurementRepositoryMock := &test.MeasurementRepositoryMock{}
			nodeServiceMock := &test.NodeServiceMock{}
			clockMock := &clock.ClockMock{}
			s := app.NewMeasurementService(measurementRepositoryMock, nodeServiceMock, clockMock)

			df := &depFields{measurementRepositoryMock: measurementRepositoryMock, nodeServiceMock: nodeServiceMock}
			tt.on(df)

			// When
			resultMeasurementDeletedID, err := s.DeleteMeasurement(context.Background(), tt.in.measurementID)

			// Then
			tt.assert(t, &output{resultMeasurementDeletedID, err})
			measurementRepositoryMock.AssertExpectations(t)
			nodeServiceMock.AssertExpectations(t)

		})
	}
}

func TestMeasurementService_AddNodeMeasurements(t *testing.T) {

	type depFields struct {
		measurementRepositoryMock *test.MeasurementRepositoryMock
		nodeServiceMock           *test.NodeServiceMock
		clockMock                 *clock.ClockMock
	}
	type input struct {
		fabricCode   string
		measurements []measurement.Measurement
	}
	type output struct {
		addedMeasurements []measurement.Measurement
		err               error
	}

	timeNow := time.Date(2024, 06, 20, 20, 15, 30, 0, time.UTC)
	expectedMeasurements := []measurement.Measurement{
		{
			ID:        uint64(10),
			NodeID:    uint64(1),
			Value:     20.0,
			Type:      "Web",
			Timestamp: timeNow,
		},
		{
			ID:        uint64(20),
			NodeID:    uint64(1),
			Value:     25.0,
			Type:      "Web",
			Timestamp: timeNow,
		},
	}

	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "add measurements by fabricCode successfully",
			in:   input{fabricCode: "abcd", measurements: expectedMeasurements},
			on: func(df *depFields) {
				df.nodeServiceMock.On("GetNodeIDByFabricCode", "abcd").Return(uint64(1), nil)
				df.measurementRepositoryMock.On("AddMeasurement", expectedMeasurements[0]).Return(expectedMeasurements[0], nil).Once()
				df.measurementRepositoryMock.On("AddMeasurement", expectedMeasurements[1]).Return(expectedMeasurements[1], nil).Once()
				df.clockMock.On("Time").Return(timeNow)
				df.nodeServiceMock.On("UpdateNodeLastUpdated", uint64(1), timeNow).Return(nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, expectedMeasurements, out.addedMeasurements)
			},
		},
		{
			name: "error getting node by ID",
			in:   input{fabricCode: "abcd", measurements: expectedMeasurements},
			on: func(df *depFields) {
				df.nodeServiceMock.On("GetNodeIDByFabricCode", "abcd").Return(0, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.Error(t, out.err)
				assert.ErrorContains(t, out.err, "node not found: test")
				assert.Nil(t, out.addedMeasurements)
			},
		},
		{
			name: "error adding measurement",
			in:   input{fabricCode: "abcd", measurements: expectedMeasurements},
			on: func(df *depFields) {
				df.nodeServiceMock.On("GetNodeIDByFabricCode", "abcd").Return(uint64(1), nil)
				df.measurementRepositoryMock.On("AddMeasurement", expectedMeasurements[0]).Return(expectedMeasurements[0], nil).Once()
				df.measurementRepositoryMock.On("AddMeasurement", expectedMeasurements[1]).Return(expectedMeasurements[1], fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.Error(t, out.err)
				assert.ErrorContains(t, out.err, "error adding measurement: test")
				assert.Nil(t, out.addedMeasurements)
			},
		},
		{
			name: "error updating last updated time",
			in:   input{fabricCode: "abcd", measurements: expectedMeasurements},
			on: func(df *depFields) {
				df.nodeServiceMock.On("GetNodeIDByFabricCode", "abcd").Return(uint64(1), nil)
				df.measurementRepositoryMock.On("AddMeasurement", expectedMeasurements[0]).Return(expectedMeasurements[0], nil).Once()
				df.measurementRepositoryMock.On("AddMeasurement", expectedMeasurements[1]).Return(expectedMeasurements[1], nil).Once()
				df.clockMock.On("Time").Return(timeNow)
				df.nodeServiceMock.On("UpdateNodeLastUpdated", uint64(1), timeNow).Return(fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.Error(t, out.err)
				assert.ErrorContains(t, out.err, "error updating node last_updated: test")
				assert.Nil(t, out.addedMeasurements)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			measurementRepositoryMock := &test.MeasurementRepositoryMock{}
			nodeServiceMock := &test.NodeServiceMock{}
			clockMock := &clock.ClockMock{}
			s := app.NewMeasurementService(measurementRepositoryMock, nodeServiceMock, clockMock)

			df := &depFields{measurementRepositoryMock: measurementRepositoryMock, nodeServiceMock: nodeServiceMock, clockMock: clockMock}
			tt.on(df)

			// When
			resultMeasurements, err := s.AddNodeMeasurements(context.Background(), tt.in.fabricCode, tt.in.measurements)

			// Then
			tt.assert(t, &output{resultMeasurements, err})
			measurementRepositoryMock.AssertExpectations(t)
			nodeServiceMock.AssertExpectations(t)
			clockMock.AssertExpectations(t)

		})
	}
}
