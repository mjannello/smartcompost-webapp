package http

import (
	"github.com/gorilla/mux"
	measurementport "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/port"
	nodeport "github.com/mjannello/smartcompost-webapp/backend/internal/node/port"
	"net/http"
)

type RouterHandler interface {
	RouteURLs(router *mux.Router)
}

func NewRouterHandler(nodeHandler nodeport.Handler, measurementHandler measurementport.Handler) RouterHandler {
	return &routerHandler{
		nodeHandler:        nodeHandler,
		measurementHandler: measurementHandler,
	}
}

type routerHandler struct {
	nodeHandler        nodeport.Handler
	measurementHandler measurementport.Handler
}

func (r *routerHandler) RouteURLs(router *mux.Router) {
	prefix := "/api"
	nodesPrefix := prefix + "/nodes"

	// Nodes
	router.HandleFunc(nodesPrefix, r.nodeHandler.GetNodes).Methods(http.MethodGet)
	router.HandleFunc(nodesPrefix+"/{nodeID}", r.nodeHandler.GetNodeByID).Methods(http.MethodGet)
	router.HandleFunc(nodesPrefix, r.nodeHandler.CreateNode).Methods(http.MethodPost)
	router.HandleFunc(nodesPrefix+"/{nodeID}", r.nodeHandler.UpdateNode).Methods(http.MethodPut)
	router.HandleFunc(nodesPrefix+"/{nodeID}", r.nodeHandler.DeleteNode).Methods(http.MethodDelete)

	// Measurements
	measurementsSuffix := "/{serialNumber}/measurements"
	measurementsPrefix := nodesPrefix + measurementsSuffix
	router.HandleFunc(measurementsPrefix, r.measurementHandler.GetMeasurementsByNode).Methods(http.MethodGet)
	//router.HandleFunc(measurementsPrefix+"/{measurementID}", r.measurementHandler.GetMeasurementByID).Methods(http.MethodGet)
	router.HandleFunc(measurementsPrefix, r.measurementHandler.AddMeasurement).Methods(http.MethodPost)
	router.HandleFunc(measurementsPrefix+"/{measurementID}", r.measurementHandler.UpdateMeasurement).Methods(http.MethodPut)
	router.HandleFunc(measurementsPrefix+"/{measurementID}", r.measurementHandler.DeleteMeasurement).Methods(http.MethodDelete)

	// Access Point
	accessPointPrefix := prefix + "/ap"
	apMeasurementsPrefix := accessPointPrefix + measurementsSuffix
	router.HandleFunc(apMeasurementsPrefix, r.measurementHandler.AddMeasurementAP).Methods(http.MethodPost)
}
