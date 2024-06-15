package http

import (
	"github.com/gorilla/mux"
	measurementport "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/port"
	nodeport "github.com/mjannello/smartcompost-webapp/backend/internal/node/port"
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
	router.HandleFunc(nodesPrefix, r.nodeHandler.GetNodes).Methods("GET")
	router.HandleFunc(nodesPrefix+"/{nodeID}", r.nodeHandler.GetNodeByID).Methods("GET")
	router.HandleFunc(nodesPrefix+"/{nodeID}", r.nodeHandler.UpdateNode).Methods("PUT")
	router.HandleFunc(nodesPrefix+"/{nodeID}", r.nodeHandler.DeleteNode).Methods("DELETE")

	// Measurements
	measurementsPrefix := nodesPrefix + "/{nodeID}/measurements"
	router.HandleFunc(measurementsPrefix, r.measurementHandler.GetMeasurementsByNodeID).Methods("GET")
	router.HandleFunc(measurementsPrefix, r.measurementHandler.AddMeasurement).Methods("POST")
	router.HandleFunc(measurementsPrefix+"/{measurementID}", r.measurementHandler.GetMeasurementByID).Methods("GET")
	router.HandleFunc(measurementsPrefix+"/{measurementID}", r.measurementHandler.UpdateMeasurement).Methods("PUT")
	router.HandleFunc(measurementsPrefix+"/{measurementID}", r.measurementHandler.DeleteMeasurement).Methods("DELETE")

}
