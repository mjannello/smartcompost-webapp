package port_test

import (
	"github.com/gorilla/mux"
	measurementmodel "github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement/port"
	"github.com/mjannello/smartcompost-webapp/backend/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMeasurementsByNodeID(t *testing.T) {
	type input struct {
		serialNumber string
	}
	type output struct {
		response *http.Response
	}
	type depFields struct {
		service *test.MeasurementServiceMock
	}

	timeNow := time.Date(2024, 06, 20, 20, 15, 30, 0, time.UTC)
	expectedMeasurements := []measurementmodel.Measurement{
		test.MakeMeasurement(uint64(10), uint64(1), 20.0, "Web", timeNow),
		test.MakeMeasurement(uint64(20), uint64(1), 25.0, "Web", timeNow),
	}

	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "valid serialNumber, no measurements found",
			in:   input{serialNumber: "abcd"},
			on: func(df *depFields) {
				df.service.On("GetMeasurementsByNode", "abcd").Return([]measurementmodel.Measurement{}, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.Equal(t, http.StatusOK, out.response.StatusCode)
			},
		},
		{
			name: "valid serialNumber, measurements found",
			in:   input{serialNumber: "abcd"},
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
			measurementService := test.MeasurementServiceMock{}
			handler := port.NewMeasurementHandler(&measurementService)

			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/test", nil)
			req = mux.SetURLVars(req, map[string]string{"serialNumber": tt.in.serialNumber})

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
