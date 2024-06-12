package port

import (
	"encoding/json"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node"
	nodeapp "github.com/mjannello/smartcompost-webapp/backend/internal/node/app"
	"net/http"
	"time"
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
	var serializedNodes []NodeRestModel
	for _, n := range nodes {
		serializedNode := AppToRestNodeModel(n)
		serializedNodes = append(serializedNodes, serializedNode)
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(serializedNodes); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

type NodeRestModel struct {
	ID           uint64      `json:"id"`
	Description  string      `json:"description"`
	Type         string      `json:"type"`
	LastUpdated  time.Time   `json:"last_updated"`
	Measurements interface{} `json:"measurements,omitempty"`
}

func RestNodeModelToApp(nodeRestModel NodeRestModel) node.Node {
	return node.Node{
		ID:          nodeRestModel.ID,
		Description: nodeRestModel.Description,
		Type:        nodeRestModel.Type,
		LastUpdated: nodeRestModel.LastUpdated,
	}

}

func AppToRestNodeModel(n node.Node) NodeRestModel {
	return NodeRestModel{
		ID:          n.ID,
		Description: n.Description,
		Type:        n.Type,
		LastUpdated: n.LastUpdated,
	}
}
