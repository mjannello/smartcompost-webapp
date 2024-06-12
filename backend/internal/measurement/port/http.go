package port

import (
	"encoding/json"
	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	measurementapp "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/app"
)

type Handler interface {
	GetMeasurementsByNodeID(w http.ResponseWriter, r *http.Request)
	AddMeasurement(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	measurementService measurementapp.MeasurementService
}

func NewHTTPHandler(ms measurementapp.MeasurementService) Handler {
	return &handler{measurementService: ms}
}

type AddMeasurementRequest struct {
	NodeType         string                 `json:"node_type"`
	LastUpdated      string                 `json:"last_updated"`
	NodeMeasurements []MeasurementRestModel `json:"node_measurements"`
}

type MeasurementRestModel struct {
	ID        uint64  `json:"id"`
	NodeID    uint64  `json:"node_id"`
	Value     float64 `json:"value"`
	Timestamp string  `json:"timestamp"`
	Type      string  `json:"type"`
}

func (h *handler) GetMeasurementsByNodeID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeIDStr, ok := vars["nodeID"]
	if !ok {
		http.Error(w, "nodeID is required", http.StatusBadRequest)
		return
	}

	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid nodeID", http.StatusBadRequest)
		return
	}

	measurements, err := h.measurementService.GetMeasurementsByNodeID(r.Context(), nodeID)
	if err != nil {
		http.Error(w, "Error getting measurements", http.StatusInternalServerError)
		return
	}

	var serializedMeasurements []MeasurementRestModel
	for _, m := range measurements {
		serializedMeasurement := AppToRestMeasurementModel(m)
		serializedMeasurements = append(serializedMeasurements, serializedMeasurement)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(serializedMeasurements); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *handler) AddMeasurement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeIDStr := vars["node_id"]
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid nodeID", http.StatusBadRequest)
		return
	}

	var reqBody AddMeasurementRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	measurements := make([]measurement.Measurement, len(reqBody.NodeMeasurements))
	for i, m := range reqBody.NodeMeasurements {
		measurements[i] = RestToAppMeasurementModel(m)
		measurements[i].NodeID = nodeID
	}
	var logger = log.New(os.Stdout, "smartcompost-webapp ", log.LstdFlags|log.Lshortfile)
	logger.Println(measurements)
	if err := h.measurementService.AddMeasurements(r.Context(), measurements); err != nil {
		http.Error(w, "error adding measurements", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Measurements added successfully"))
}
func AppToRestMeasurementModel(m measurement.Measurement) MeasurementRestModel {
	return MeasurementRestModel{
		ID:        m.ID,
		NodeID:    m.NodeID,
		Value:     m.Value,
		Timestamp: m.Timestamp.Format("2006-01-02 15:04:05"),
		Type:      m.Type,
	}
}

func RestToAppMeasurementModel(m MeasurementRestModel) measurement.Measurement {
	timestamp, _ := time.Parse(time.RFC3339, m.Timestamp)
	return measurement.Measurement{
		ID:        m.ID,
		NodeID:    m.NodeID,
		Value:     m.Value,
		Timestamp: timestamp,
		Type:      m.Type,
	}
}
