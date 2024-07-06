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
	GetMeasurementsByNode(w http.ResponseWriter, r *http.Request)
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

func (h *handler) GetMeasurementsByNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fabricCode := vars["fabricCode"]

	ctx := r.Context()
	measurements, err := h.measurementService.GetMeasurementsByNode(ctx, fabricCode)
	if err != nil {
		log.Printf("[Handler] GetMeasurementsByNode - Error getting measurements: %s", err.Error())
		http.Error(w, "Error getting measurements", http.StatusInternalServerError)
		return
	}

	var serializedMeasurements []NodeMeasurementRestModel
	for _, m := range measurements {
		serializedMeasurement := AppToRestNodeMeasurementModel(m)
		serializedMeasurements = append(serializedMeasurements, serializedMeasurement)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(serializedMeasurements); err != nil {
		log.Printf("[Handler] GetMeasurementsByNode - Error encoding response: %s", err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("[Handler] GetMeasurementsByNode - Measurements fetched successfully for Node with Fabric Code: %s", fabricCode)
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
	fabricCode, ok := vars["fabricCode"]
	if !ok {
		log.Printf("[Handler] AddMeasurement - fabricCode not provided in URI")
		http.Error(w, "fabricCode not provided", http.StatusBadRequest)
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
	createdMeasurements, err := h.measurementService.AddNodeMeasurements(ctx, fabricCode, measurementRest.LastUpdated, measurements)
	if err != nil {
		log.Printf("[Handler] AddMeasurement - Error adding measurements: %s", err.Error())
		http.Error(w, "Error adding measurements", http.StatusInternalServerError)
		return
	}

	// Convert created measurements to REST model
	createdMeasurementRestModels := AppToNodeMeasurementsRestModel(createdMeasurements)

	// Respond with created measurements
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdMeasurementRestModels); err != nil {
		log.Printf("[Handler] AddMeasurement - Error encoding response: %s", err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("[Handler] AddMeasurement - Measurements added successfully")
}

func AppToNodeMeasurementsRestModel(measurements []measurementmodel.Measurement) []NodeMeasurementRestModel {
	createdMeasurementRestModels := make([]NodeMeasurementRestModel, len(measurements))
	for i, m := range measurements {
		createdMeasurementRestModels[i] = AppToRestNodeMeasurementModel(m)
	}
	return createdMeasurementRestModels

}

func AppToRestNodeMeasurementModel(measurement measurementmodel.Measurement) NodeMeasurementRestModel {
	return NodeMeasurementRestModel{
		ID:        measurement.ID,
		Value:     measurement.Value,
		Timestamp: measurement.Timestamp,
		Type:      measurement.Type,
	}
}

func AppToRestMeasurementModel(m measurementmodel.Measurement) MeasurementRestModel {
	restModel := MeasurementRestModel{
		LastUpdated:      time.Now().UTC(),
		NodeMeasurements: []NodeMeasurementRestModel{},
	}

	nodeMeasurement := NodeMeasurementRestModel{
		Value:     m.Value,
		Timestamp: m.Timestamp,
		Type:      m.Type,
	}
	restModel.NodeMeasurements = append(restModel.NodeMeasurements, nodeMeasurement)

	return restModel
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

type MeasurementRestModel struct {
	LastUpdated      time.Time                  `json:"last_updated"`
	NodeMeasurements []NodeMeasurementRestModel `json:"node_measurements"`
}

type NodeMeasurementRestModel struct {
	ID        uint64    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}
