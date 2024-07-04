package port_test

import (
	"context"
	"github.com/gorilla/mux"
	measurementmodel "github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MeasurementServiceMock struct {
	mock.Mock
}

func (m *MeasurementServiceMock) GetMeasurementsByNode(_ context.Context, fabricCode string) ([]measurementmodel.Measurement, error) {
	args := m.Called(fabricCode)
	return args.Get(0).([]measurementmodel.Measurement), args.Error(1)
}

func (m *MeasurementServiceMock) GetMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurementmodel.Measurement, error) {
	args := m.Called(ctx, nodeID)
	return args.Get(0).([]measurementmodel.Measurement), args.Error(1)
}

func (m *MeasurementServiceMock) GetMeasurementByID(ctx context.Context, measurementID uint64) (measurementmodel.Measurement, error) {
	args := m.Called(ctx, measurementID)
	return args.Get(0).(measurementmodel.Measurement), args.Error(1)
}

func (m *MeasurementServiceMock) UpdateMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	args := m.Called(ctx, measurement)
	return args.Get(0).(measurementmodel.Measurement), args.Error(1)
}

func (m *MeasurementServiceMock) DeleteMeasurement(ctx context.Context, measurementID uint64) (uint64, error) {
	args := m.Called(ctx, measurementID)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MeasurementServiceMock) AddNodeMeasurements(ctx context.Context, fabricCode string, measurements []measurementmodel.Measurement) ([]measurementmodel.Measurement, error) {
	args := m.Called(ctx, fabricCode, measurements)
	return args.Get(0).([]measurementmodel.Measurement), args.Error(1)
}

func TestGetMeasurementsByNodeID(t *testing.T) {
	type input struct {
		fabricCode string
	}
	type output struct {
		response *http.Response
	}
	type depFields struct {
		service *MeasurementServiceMock
	}

	timeNow := time.Date(2024, 06, 20, 20, 15, 30, 0, time.UTC)
	expectedMeasurements := []measurementmodel.Measurement{
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
			name: "valid fabricCode, no measurements found",
			in:   input{fabricCode: "abcd"},
			on: func(df *depFields) {
				df.service.On("GetMeasurementsByNode", "abcd").Return([]measurementmodel.Measurement{}, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.Equal(t, http.StatusOK, out.response.StatusCode)
			},
		},
		{
			name: "valid fabricCode, measurements found",
			in:   input{fabricCode: "abcd"},
			on: func(df *depFields) {
				df.service.On("GetMeasurementsByNode", "abcd").Return(expectedMeasurements, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.Equal(t, http.StatusOK, out.response.StatusCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			measurementService := MeasurementServiceMock{}
			handler := port.NewMeasurementHandler(&measurementService)

			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/test", nil)
			req = mux.SetURLVars(req, map[string]string{"fabricCode": tt.in.fabricCode})

			writer := httptest.NewRecorder()
			f := &depFields{service: &measurementService}
			tt.on(f)

			// When
			handler.GetMeasurementsByNode(writer, req)
			response := writer.Result()

			// Then
			tt.assert(t, &output{response})
			measurementService.AssertExpectations(t)
		})
	}
}
