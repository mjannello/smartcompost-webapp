package port

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	measurementapp "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/app"
)

type Handler interface {
	GetMeasurementsByNodeID(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	measurementService measurementapp.MeasurementService
}

func NewHTTPHandler(ms measurementapp.MeasurementService) Handler {
	return &handler{measurementService: ms}
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(measurements); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
