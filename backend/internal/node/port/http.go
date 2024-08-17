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

// GetNodes
// @Summary Get all nodes
// @Description Get a list of all nodes
// @Tags nodes
// @Accept  json
// @Produce  json
// @Success 200 {array} NodeRestModel
// @Failure 500 {string} string "Error getting nodes"
// @Router /api/nodes [get]
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

// GetNodeByID
// @Summary Get a node by ID
// @Description Get a specific node by its ID
// @Tags nodes
// @Accept  json
// @Produce  json
// @Param nodeID path int true "Node ID"
// @Success 200 {object} NodeRestModel
// @Failure 400 {string} string "Invalid nodeID"
// @Failure 404 {string} string "Node not found"
// @Failure 500 {string} string "Error getting node"
// @Router /api/nodes/{nodeID} [get]
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

// CreateNode
// @Summary Create a new node
// @Description Create a new node with the provided information
// @Tags nodes
// @Accept  json
// @Produce  json
// @Param request body CreateNodeRequest true "Node data"
// @Success 201 {object} NodeRestModel
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Error creating node"
// @Router /api/nodes [post]
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

// UpdateNode
// @Summary Update an existing node
// @Description Update a node's information by its ID
// @Tags nodes
// @Accept  json
// @Produce  json
// @Param nodeID path int true "Node ID"
// @Param request body PatchNodeModel true "Updated node data"
// @Success 200 {object} NodeRestModel
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Node not found"
// @Failure 500 {string} string "Error updating node"
// @Router /api/nodes/{nodeID} [patch]
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

// DeleteNode
// @Summary Delete a node by ID
// @Description Delete a specific node by its ID
// @Tags nodes
// @Accept  json
// @Produce  json
// @Param nodeID path int true "Node ID"
// @Success 204
// @Failure 400 {string} string "Invalid nodeID"
// @Failure 404 {string} string "Node not found"
// @Failure 500 {string} string "Error deleting node"
// @Router /api/nodes/{nodeID} [delete]
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
