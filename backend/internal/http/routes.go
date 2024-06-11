package http

import (
	"github.com/gorilla/mux"
	nodeport "github.com/mjannello/smartcompost-webapp/backend/internal/node/port"
)

type RouterHandler interface {
	RouteURLs(router *mux.Router)
}

type routerHandler struct {
	nodeHandler nodeport.Handler
}

func (r *routerHandler) RouteURLs(router *mux.Router) {
	prefix := "/api/nodes"
	router.HandleFunc(prefix, r.nodeHandler.GetNodes).Methods("GET")
}

func NewRouterHandler(nodeHandler nodeport.Handler) RouterHandler {
	return &routerHandler{nodeHandler: nodeHandler}
}
