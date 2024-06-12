package http

import (
	"github.com/gorilla/mux"
	measurementport "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/port"
	nodeport "github.com/mjannello/smartcompost-webapp/backend/internal/node/port"
)

type RouterHandler interface {
	RouteURLs(router *mux.Router)
}

type routerHandler struct {
	nodeHandler        nodeport.Handler
	measurementHandler measurementport.Handler
}

func (r *routerHandler) RouteURLs(router *mux.Router) {
	prefix := "/api"
	nodesPrefix := prefix + "/nodes"
	measurementSuffix := "/measurements"
	router.HandleFunc(nodesPrefix, r.nodeHandler.GetNodes).Methods("GET")
	router.HandleFunc(prefix+"/{nodeID}"+measurementSuffix, r.measurementHandler.GetMeasurementsByNodeID).Methods("GET")
	router.HandleFunc("/api/{node_id}/add_measurement", r.measurementHandler.AddMeasurement).Methods("POST")

}

func NewRouterHandler(nodeHandler nodeport.Handler, measurementHandler measurementport.Handler) RouterHandler {
	return &routerHandler{nodeHandler: nodeHandler, measurementHandler: measurementHandler}
}
