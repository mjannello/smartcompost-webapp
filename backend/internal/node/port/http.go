package port

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	nodemodel "github.com/mjannello/smartcompost-webapp/backend/internal/node"
	nodeapp "github.com/mjannello/smartcompost-webapp/backend/internal/node/app"
)

type Handler interface {
	GetNodes(w http.ResponseWriter, r *http.Request)
	GetNodeByID(w http.ResponseWriter, r *http.Request)
	CreateNode(w http.ResponseWriter, r *http.Request)
	UpdateNode(w http.ResponseWriter, r *http.Request)
	DeleteNode(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	nodeService nodeapp.NodeService
}

func NewNodeHandler(nodeService nodeapp.NodeService) Handler {
	return &handler{
		nodeService: nodeService,
	}
}

type CreateNodeRequest struct {
	SerialNumber string `json:"serial_number"`
	Description  string `json:"description"`
	Model        string `json:"model"`
}

func (h *handler) GetNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Printf("[Handler] GetNodes - Received request %s %s", r.Method, r.URL.Path)

	nodes, err := h.nodeService.GetNodes(ctx)
	if err != nil {
		log.Printf("[Handler] GetNodes - Error getting nodes: %s", err.Error())
		http.Error(w, "Error getting nodes", http.StatusInternalServerError)
		return
	}

	var serializedNodes []NodeRestModel
	for _, n := range nodes {
		serializedNode := AppToRestNodeModel(n)
		serializedNodes = append(serializedNodes, serializedNode)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(serializedNodes); err != nil {
		log.Printf("[Handler] GetNodes - Error encoding response: %s", err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Println("[Handler] GetNodes - Nodes fetched successfully.")
}

func (h *handler) GetNodeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeIDStr := vars["nodeID"]
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] GetNode - Invalid nodeID")
		http.Error(w, "Invalid nodeID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	node, err := h.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		log.Printf("[Handler] GetNode - Error getting node: %s", err.Error())
		http.Error(w, "Error getting node", http.StatusInternalServerError)
		return
	}

	serializedNode := AppToRestNodeModel(node)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(serializedNode); err != nil {
		log.Printf("[Handler] GetNode - Error encoding response: %s", err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("[Handler] GetNode - Node fetched successfully. ID: %d", nodeID)
}

func (h *handler) CreateNode(w http.ResponseWriter, r *http.Request) {
	var createNodeReq CreateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&createNodeReq); err != nil {
		log.Printf("[Handler] CreateNode - Invalid request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	node, err := h.nodeService.CreateNode(ctx, createNodeReq.Description, createNodeReq.Model)
	if err != nil {
		log.Printf("[Handler] CreateNode - Error creating node: %s", err.Error())
		http.Error(w, "Error creating node", http.StatusInternalServerError)
		return
	}

	serializedNode := AppToRestNodeModel(node)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(serializedNode); err != nil {
		log.Printf("[Handler] CreateNode - Error encoding response: %s", err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("[Handler] CreateNode - Node created successfully. ID: %d", node.ID)
}

func (h *handler) UpdateNode(w http.ResponseWriter, r *http.Request) {
	// TODO: Fix last_updated unmarshalling/format
	vars := mux.Vars(r)
	nodeIDStr := vars["nodeID"]
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] UpdateNode - Invalid nodeID")
		http.Error(w, "Invalid nodeID", http.StatusBadRequest)
		return
	}

	var patch PatchNodeModel
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		log.Printf("[Handler] UpdateNode - Invalid request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	existingNode, err := h.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		log.Printf("[Handler] UpdateNode - Node not found: %s", err.Error())
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	if patch.SerialNumber != nil {
		existingNode.SerialNumber = *patch.SerialNumber
	}
	if patch.Description != nil {
		existingNode.Description = *patch.Description
	}
	if patch.Model != nil {
		existingNode.Model = *patch.Model
	}

	updatedNode, err := h.nodeService.UpdateNode(ctx, existingNode)
	if err != nil {
		log.Printf("[Handler] UpdateNode - Error updating node: %s", err.Error())
		http.Error(w, "Error updating node", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("[Handler] UpdateNode - Node updated successfully. ID: %d", updatedNode.ID)
}

func (h *handler) DeleteNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeIDStr := vars["nodeID"]
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 64)
	if err != nil {
		log.Println("[Handler] DeleteNode - Invalid nodeID")
		http.Error(w, "Invalid nodeID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	deletedID, err := h.nodeService.DeleteNode(ctx, nodeID)
	if err != nil {
		log.Printf("[Handler] DeleteNode - Error deleting node: %s", err.Error())
		http.Error(w, "Error deleting node", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("[Handler] DeleteNode - Node deleted successfully. ID: %d", deletedID)
}

type NodeRestModel struct {
	ID           uint64      `json:"id"`
	SerialNumber string      `json:"serial_number"`
	Description  string      `json:"description"`
	Model        string      `json:"model"`
	DateCreated  time.Time   `json:"date_created"`
	LastUpdated  time.Time   `json:"last_updated"`
	Measurements interface{} `json:"measurements,omitempty"`
}

type PatchNodeModel struct {
	SerialNumber *string `json:"serial_number,omitempty"`
	Description  *string `json:"description,omitempty"`
	Model        *string `json:"model,omitempty"`
}

func RestNodeModelToApp(nodeRestModel NodeRestModel) nodemodel.Node {
	return nodemodel.Node{
		ID:           nodeRestModel.ID,
		SerialNumber: nodeRestModel.SerialNumber,
		Description:  nodeRestModel.Description,
		Model:        nodeRestModel.Model,
		DateCreated:  nodeRestModel.DateCreated,
		LastUpdated:  nodeRestModel.LastUpdated,
	}

}

func AppToRestNodeModel(n nodemodel.Node) NodeRestModel {
	return NodeRestModel{
		ID:           n.ID,
		SerialNumber: n.SerialNumber,
		Description:  n.Description,
		Model:        n.Model,
		DateCreated:  n.DateCreated,
		LastUpdated:  n.LastUpdated,
	}
}
