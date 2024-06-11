package port

import (
	"encoding/json"
	nodeapp "github.com/mjannello/smartcompost-webapp/backend/internal/node/app"
	"net/http"
)

type Handler interface {
	GetNodes(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	nodeService nodeapp.NodeService
}

func NewHTTPHandler(nodeService nodeapp.NodeService) Handler {
	return &handler{nodeService: nodeService}
}

func (h *handler) GetNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nodes, err := h.nodeService.GetNodes(ctx)
	if err != nil {
		http.Error(w, "Error getting nodes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
