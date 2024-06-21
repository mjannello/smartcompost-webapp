package port

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	measurementmodel "github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
	measurementapp "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/app"
)

type Handler interface {
	GetMeasurementsByNodeID(w http.ResponseWriter, r *http.Request)
	GetMeasurementByID(w http.ResponseWriter, r *http.Request)
	UpdateMeasurement(w http.ResponseWriter, r *http.Request)
	DeleteMeasurement(w http.ResponseWriter, r *http.Request)
	AddMeasurement(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	measurementService measurementapp.MeasurementService
}

func NewMeasurementHandler(measurementService measurementapp.MeasurementService) Handler {
	return &handler{
		measurementService: measurementService,
	}
}

func (h *handler) GetMeasurementsByNodeID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeIDStr := vars["nodeID"]
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] GetMeasurementsByNodeID - Invalid nodeID")
		http.Error(w, "Invalid nodeID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	measurements, err := h.measurementService.GetMeasurementsByNodeID(ctx, nodeID)
	if err != nil {
		log.Printf("[Handler] GetMeasurementsByNodeID - Error getting measurements: %s", err.Error())
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
		log.Printf("[Handler] GetMeasurementsByNodeID - Error encoding response: %s", err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("[Handler] GetMeasurementsByNodeID - Measurements fetched successfully for NodeID: %d", nodeID)
}

func (h *handler) GetMeasurementByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeIDStr := vars["nodeID"]
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] GetMeasurementByID - Invalid nodeID")
		http.Error(w, "Invalid nodeID", http.StatusBadRequest)
		return
	}

	measurementIDStr := vars["measurementID"]
	measurementID, err := strconv.ParseUint(measurementIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] GetMeasurementByID - Invalid measurementID")
		http.Error(w, "Invalid measurementID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Obtain the measurement from the service
	measurement, err := h.measurementService.GetMeasurementByID(ctx, measurementID)
	if err != nil {
		log.Printf("[Handler] GetMeasurementByID - Error getting measurement: %s", err.Error())
		http.Error(w, "Error getting measurement", http.StatusInternalServerError)
		return
	}

	// Validate that the measurement belongs to the correct node
	if measurement.NodeID != nodeID {
		log.Println("[Handler] GetMeasurementByID - Measurement does not belong to the specified node")
		http.Error(w, "Measurement does not belong to the specified node", http.StatusNotFound)
		return
	}

	// Respond with the measurement
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(measurement); err != nil {
		log.Printf("[Handler] GetMeasurementByID - Error encoding response: %s", err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("[Handler] GetMeasurementByID - Measurement fetched successfully. ID: %d", measurementID)
}

func (h *handler) UpdateMeasurement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeIDStr := vars["nodeID"]
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] UpdateMeasurement - Invalid nodeID")
		http.Error(w, "Invalid nodeID", http.StatusBadRequest)
		return
	}

	measurementIDStr := vars["measurementID"]
	measurementID, err := strconv.ParseUint(measurementIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] UpdateMeasurement - Invalid measurementID")
		http.Error(w, "Invalid measurementID", http.StatusBadRequest)
		return
	}

	var measurement measurementmodel.Measurement
	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		log.Printf("[Handler] UpdateMeasurement - Invalid request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	measurement.ID = measurementID
	measurement.NodeID = nodeID

	ctx := r.Context()
	updatedMeasurement, err := h.measurementService.UpdateMeasurement(ctx, measurement)
	if err != nil {
		log.Printf("[Handler] UpdateMeasurement - Error updating measurement: %s", err.Error())
		http.Error(w, "Error updating measurement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("[Handler] UpdateMeasurement - Measurement updated successfully. ID: %d", updatedMeasurement.ID)
}

func (h *handler) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeIDStr := vars["nodeID"]
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] DeleteMeasurement - Invalid nodeID")
		http.Error(w, "Invalid nodeID", http.StatusBadRequest)
		return
	}

	measurementIDStr := vars["measurementID"]
	measurementID, err := strconv.ParseUint(measurementIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] DeleteMeasurement - Invalid measurementID")
		http.Error(w, "Invalid measurementID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	measurement, err := h.measurementService.GetMeasurementByID(ctx, measurementID)
	if err != nil {
		log.Printf("[Handler] DeleteMeasurement - Error getting measurement: %s", err.Error())
		http.Error(w, "Error getting measurement", http.StatusInternalServerError)
		return
	}

	if measurement.NodeID != nodeID {
		log.Println("[Handler] DeleteMeasurement - Measurement does not belong to the specified node")
		http.Error(w, "Measurement does not belong to the specified node", http.StatusNotFound)
		return
	}

	deletedID, err := h.measurementService.DeleteMeasurement(ctx, measurementID)
	if err != nil {
		log.Printf("[Handler] DeleteMeasurement - Error deleting measurement: %s", err.Error())
		http.Error(w, "Error deleting measurement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("[Handler] DeleteMeasurement - Measurement deleted successfully. ID: %d", deletedID)
}

func (h *handler) AddMeasurement(w http.ResponseWriter, r *http.Request) {
	// Extract nodeID from URI params
	vars := mux.Vars(r)
	nodeIDStr, ok := vars["nodeID"]
	if !ok {
		log.Printf("[Handler] AddMeasurement - nodeID not provided in URI")
		http.Error(w, "nodeID not provided", http.StatusBadRequest)
		return
	}

	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		log.Printf("[Handler] AddMeasurement - Invalid nodeID: %s", err.Error())
		http.Error(w, "Invalid nodeID", http.StatusBadRequest)
		return
	}

	// Decode request body
	var measurementRest MeasurementRestModel
	if err := json.NewDecoder(r.Body).Decode(&measurementRest); err != nil {
		log.Printf("[Handler] AddMeasurement - Invalid request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert RestModel to AppModel
	measurements := RestMeasurementModelToApp(measurementRest)

	// Add measurements to the node
	ctx := r.Context()
	createdMeasurements, err := h.measurementService.AddNodeMeasurements(ctx, nodeID, measurements)
	if err != nil {
		log.Printf("[Handler] AddMeasurement - Error adding measurements: %s", err.Error())
		http.Error(w, "Error adding measurements", http.StatusInternalServerError)
		return
	}

	// Respond with created measurements
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdMeasurements); err != nil {
		log.Printf("[Handler] AddMeasurement - Error encoding response: %s", err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("[Handler] AddMeasurement - Measurements added successfully")
}

func RestMeasurementModelToApp(measurementRestModel MeasurementRestModel) []measurementmodel.Measurement {
	measurements := make([]measurementmodel.Measurement, len(measurementRestModel.NodeMeasurements))
	for i, nm := range measurementRestModel.NodeMeasurements {
		measurements[i] = measurementmodel.Measurement{
			Value:     nm.Value,
			Timestamp: nm.Timestamp,
			Type:      nm.Type,
		}
	}
	return measurements
}

func AppToRestMeasurementModel(m measurementmodel.Measurement) MeasurementRestModel {
	restModel := MeasurementRestModel{
		NodeType:         m.Type,
		LastUpdated:      time.Now().UTC(),
		NodeMeasurements: []NodeMeasurementModel{},
	}

	nodeMeasurement := NodeMeasurementModel{
		Value:     m.Value,
		Timestamp: m.Timestamp,
		Type:      m.Type,
	}
	restModel.NodeMeasurements = append(restModel.NodeMeasurements, nodeMeasurement)

	return restModel
}

type MeasurementRestModel struct {
	ID               uint64                 `json:"id,omitempty"`
	NodeType         string                 `json:"node_type"`
	LastUpdated      time.Time              `json:"last_updated"`
	NodeMeasurements []NodeMeasurementModel `json:"node_measurements"`
}

type NodeMeasurementModel struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}
