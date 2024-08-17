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
	AddMeasurementAP(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	measurementService measurementapp.MeasurementService
}

func NewMeasurementHandler(measurementService measurementapp.MeasurementService) Handler {
	return &handler{
		measurementService: measurementService,
	}
}

// GetMeasurementsByNode
// @Summary Get measurements by node serial number
// @Description Get all measurements for a specific node by its serial number
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param serialNumber path string true "Node Serial Number"
// @Success 200 {array} MeasurementRestModel
// @Failure 500 {string} string "Error getting measurements"
// @Router /api/nodes/{serialNumber}/measurements [get]
func (h *handler) GetMeasurementsByNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serialNumber := vars["serialNumber"]

	ctx := r.Context()
	measurements, err := h.measurementService.GetMeasurementsByNode(ctx, serialNumber)
	if err != nil {
		log.Printf("[Handler] GetMeasurementsByNode - Error getting measurements: %s", err.Error())
		http.Error(w, "Error getting measurements", http.StatusInternalServerError)
		return
	}

	var serializedMeasurements []MeasurementRestModel
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

	log.Printf("[Handler] GetMeasurementsByNode - Measurements fetched successfully for Node with Serial Number: %s", serialNumber)
}

// UpdateMeasurement
// @Summary Update a measurement
// @Description Update a specific measurement by its ID and node ID
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param nodeID path int true "Node ID"
// @Param measurementID path int true "Measurement ID"
// @Param request body MeasurementRestModel true "Measurement data"
// @Success 200 {object} MeasurementRestModel
// @Failure 400 {string} string "Invalid request body or parameters"
// @Failure 500 {string} string "Error updating measurement"
// @Router /api/nodes/{nodeID}/measurements/{measurementID} [patch]
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

// DeleteMeasurement
// @Summary Delete a measurement
// @Description Delete a specific measurement by its ID and node serial number
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param serialNumber path string true "Node Serial Number"
// @Param measurementID path int true "Measurement ID"
// @Success 204
// @Failure 400 {string} string "Invalid measurementID"
// @Failure 404 {string} string "Measurement not found"
// @Failure 500 {string} string "Error deleting measurement"
// @Router /api/nodes/{serialNumber}/measurements/{measurementID} [delete]
func (h *handler) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serialNumber := vars["serialNumber"]

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

	deletedID, err := h.measurementService.DeleteMeasurement(ctx, measurement, serialNumber)
	if err != nil {
		log.Printf("[Handler] DeleteMeasurement - Error deleting measurement: %s", err.Error())
		http.Error(w, "Error deleting measurement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("[Handler] DeleteMeasurement - Measurement deleted successfully. ID: %d", deletedID)
}

// AddMeasurement
// @Summary Add new measurements
// @Description Add one or more measurements to a specific node by its serial number
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param serialNumber path string true "Node Serial Number"
// @Param request body NodeMeasurementRestModel true "Measurement data"
// @Success 201 {array} MeasurementRestModel
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Error adding measurements"
// @Router /api/nodes/{serialNumber}/measurements [post]
func (h *handler) AddMeasurement(w http.ResponseWriter, r *http.Request) {
	// Extract nodeID from URI params
	vars := mux.Vars(r)
	serialNumber, ok := vars["serialNumber"]
	if !ok {
		log.Printf("[Handler] AddMeasurement - serialNumber not provided in URI")
		http.Error(w, "serialNumber not provided", http.StatusBadRequest)
		return
	}

	// Decode request body
	var measurementRest NodeMeasurementRestModel
	if err := json.NewDecoder(r.Body).Decode(&measurementRest); err != nil {
		log.Printf("[Handler] AddMeasurement - Invalid request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert RestModel to AppModel
	measurements := RestNodeMeasurementModelToApp(measurementRest)

	// Add measurements to the node
	ctx := r.Context()
	createdMeasurements, err := h.measurementService.AddNodeMeasurements(ctx, serialNumber, measurementRest.LastUpdated, measurements)
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

// AddMeasurementAP
// @Summary Add measurements to multiple nodes
// @Description Add measurements to multiple nodes and update the last updated date for each node
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param serialNumber path string true "AP Serial Number"
// @Param request body NodesMeasurementsRestModel true "Nodes Measurements data"
// @Success 201
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Error adding measurements"
// @Router /api/nodes/{serialNumber}/measurements/ap [post]
func (h *handler) AddMeasurementAP(w http.ResponseWriter, r *http.Request) {
	// Extract nodeID from URI params
	vars := mux.Vars(r)
	serialNumber, ok := vars["serialNumber"]
	if !ok {
		log.Printf("[Handler] AddMeasurementAP - serialNumber not provided in URI")
		http.Error(w, "serialNumber not provided", http.StatusBadRequest)
		return
	}

	// Decode request body
	var nodesMeasurementsRest NodesMeasurementsRestModel
	if err := json.NewDecoder(r.Body).Decode(&nodesMeasurementsRest); err != nil {
		log.Printf("[Handler] AddMeasurementAP - Invalid request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add measurements to the node
	ctx := r.Context()
	var allCreatedMeasurements []measurementmodel.Measurement
	for _, nm := range nodesMeasurementsRest.NodesMeasurements {
		measurements := RestNodeMeasurementModelToApp(nm)
		log.Printf("measurements: %+v", measurements)
		createdMeasurements, err := h.measurementService.AddNodeMeasurements(ctx, nm.SerialNumber, nm.LastUpdated, measurements)
		if err != nil {
			log.Printf("[Handler] AddMeasurementAP - Error adding measurements for node %s: %s", nm.SerialNumber, err.Error())
			http.Error(w, "Error adding measurements", http.StatusInternalServerError)
			return
		}
		allCreatedMeasurements = append(allCreatedMeasurements, createdMeasurements...)
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("[Handler] AddMeasurementAP - %d Measurements added successfully to AP %s", len(allCreatedMeasurements), serialNumber)

	// Modify LastUpdated date for the AP node
	err := h.measurementService.UpdateAPLastUpdated(ctx, serialNumber, nodesMeasurementsRest.LastUpdated)
	if err != nil {
		log.Printf("[Handler] AddMeasurementAP - Error updating last updated for node %s: %s", serialNumber, err.Error())
		return
	}
}

func AppToNodeMeasurementsRestModel(measurements []measurementmodel.Measurement) []MeasurementRestModel {
	createdMeasurementRestModels := make([]MeasurementRestModel, len(measurements))
	for i, m := range measurements {
		createdMeasurementRestModels[i] = AppToRestNodeMeasurementModel(m)
	}
	return createdMeasurementRestModels
}

func AppToRestNodeMeasurementModel(measurement measurementmodel.Measurement) MeasurementRestModel {
	return MeasurementRestModel{
		ID:        measurement.ID,
		Value:     measurement.Value,
		Timestamp: measurement.Timestamp,
		Type:      measurement.Type,
	}
}

func AppToRestMeasurementModel(m measurementmodel.Measurement) NodeMeasurementRestModel {
	restModel := NodeMeasurementRestModel{
		LastUpdated:  time.Now().UTC(),
		Measurements: []MeasurementRestModel{},
	}

	nodeMeasurement := MeasurementRestModel{
		Value:     m.Value,
		Timestamp: m.Timestamp,
		Type:      m.Type,
	}
	restModel.Measurements = append(restModel.Measurements, nodeMeasurement)

	return restModel
}

func RestNodeMeasurementModelToApp(nodeMeasurementRestModel NodeMeasurementRestModel) []measurementmodel.Measurement {
	measurements := make([]measurementmodel.Measurement, len(nodeMeasurementRestModel.Measurements))
	for i, nm := range nodeMeasurementRestModel.Measurements {
		measurements[i] = RestMeasurementModelToApp(nm)
	}
	return measurements
}

func RestMeasurementModelToApp(measurementRestModel MeasurementRestModel) measurementmodel.Measurement {
	return measurementmodel.Measurement{
		Value:     measurementRestModel.Value,
		Timestamp: measurementRestModel.Timestamp,
		Type:      measurementRestModel.Type,
	}
}

type NodesMeasurementsRestModel struct {
	LastUpdated       time.Time                  `json:"last_updated"`
	NodesMeasurements []NodeMeasurementRestModel `json:"nodes_measurements"`
}

type NodeMeasurementRestModel struct {
	SerialNumber string                 `json:"serial_number"`
	LastUpdated  time.Time              `json:"last_updated"`
	Measurements []MeasurementRestModel `json:"measurements"`
}

type MeasurementRestModel struct {
	ID        uint64    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}
